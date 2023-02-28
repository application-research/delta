package core

import (
	"context"
	c "delta/config"
	"delta/utils"
	"fmt"
	model "github.com/application-research/delta-db/db_models"
	"github.com/application-research/delta-db/messaging"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	fc "github.com/application-research/filclient"
	"github.com/application-research/filclient/keystore"
	"github.com/application-research/whypfs-core"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	lapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/client"
	"github.com/filecoin-project/lotus/api/v1api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/wallet"
	"github.com/filecoin-project/lotus/chain/wallet/key"
	cliutil "github.com/filecoin-project/lotus/cli/util"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	mdagipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-path/resolver"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type DeltaNode struct {
	Context            context.Context
	Node               *whypfs.Node
	Api                url.URL
	DB                 *gorm.DB
	FilClient          *fc.FilClient
	Config             *c.DeltaConfig
	Dispatcher         *Dispatcher
	DeltaTracer        *messaging.DeltaMetricsTracer
	MetaInfo           *model.InstanceMeta
	WebsocketBroadcast WebsocketBroadcast
}

type WebsocketBroadcast struct {
	ContentChannel         ContentChannel
	PieceCommitmentChannel PieceCommitmentChannel
	ContentDealChannel     ContentDealChannel
}
type ContentChannel struct {
	Clients map[*websocket.Conn]bool
	Channel chan model.Content
}
type PieceCommitmentChannel struct {
	Clients map[*websocket.Conn]bool
	Channel chan model.PieceCommitment
}

type ContentDealChannel struct {
	Clients map[*websocket.Conn]bool
	Channel chan model.ContentDeal
}

type LocalWallet struct {
	keys     map[address.Address]*key.Key
	keystore types.KeyStore
	lk       sync.Mutex
}

type GatewayHandler struct {
	bs       blockstore.Blockstore
	dserv    mdagipld.DAGService
	resolver resolver.Resolver
	node     *whypfs.Node
}

type NewLightNodeParams struct {
	Repo             string
	DefaultWalletDir string
	Config           *c.DeltaConfig
}

// NewLightNode Creating a new light node.
// > This function creates a new DeltaNode with a new light client
func NewLightNode(repo NewLightNodeParams) (*DeltaNode, error) {

	//	database
	db, err := model.OpenDatabase(repo.Config.Common.DBDSN)
	publicIp, err := GetPublicIP()
	newConfig := &whypfs.Config{
		ListenAddrs: []string{
			"/ip4/0.0.0.0/tcp/6745",
			"/ip4/" + publicIp + "/tcp/6745",
		},
		AnnounceAddrs: []string{
			"/ip4/0.0.0.0/tcp/6745",
			"/ip4/" + publicIp + "/tcp/6745",
		},
	}
	params := whypfs.NewNodeParams{
		Ctx:       context.Background(),
		Datastore: whypfs.NewInMemoryDatastore(),
		Repo:      repo.Repo,
	}
	// node
	params.Config = params.ConfigurationBuilder(newConfig)
	whypfsPeer, err := whypfs.NewNode(params)

	if err != nil {
		panic(err)
	}

	whypfsPeer.BootstrapPeers(c.BootstrapEstuaryPeers())

	//	FilClient
	api, _, err := LotusConnection(utils.LOTUS_API)
	wallet, err := SetupWallet(repo.DefaultWalletDir)
	walletAddr, err := wallet.GetDefault()
	if err != nil {
		panic(err)
	}

	fc, err := fc.NewClient(whypfsPeer.Host, api, wallet, walletAddr, whypfsPeer.Blockstore, whypfsPeer.Datastore, whypfsPeer.Config.DatastoreDir.Directory)

	if err != nil {
		panic(err)
	}

	// job dispatcher
	dispatcher := CreateNewDispatcher()

	// delta metrics tracer
	tracer := messaging.NewDeltaMetricsTracer()

	// create the global light node.
	return &DeltaNode{
		Node:        whypfsPeer,
		DB:          db,
		FilClient:   fc,
		Dispatcher:  dispatcher,
		Config:      repo.Config,
		DeltaTracer: tracer,
	}, nil
}

// LotusConnection It takes a string that contains the Lotus full node API address and returns a `v1api.FullNode` interface, a
// `jsonrpc.ClientCloser` interface, and an error
func LotusConnection(fullNodeApiInfo string) (v1api.FullNode, jsonrpc.ClientCloser, error) {
	info := cliutil.ParseApiInfo(fullNodeApiInfo)

	var api lapi.FullNode
	var closer jsonrpc.ClientCloser
	addr, err := info.DialArgs("v1")
	if err != nil {
		log.Errorf("Error getting v1 API address %s", err)
		return nil, nil, err
	}

	api, closer, err = client.NewFullNodeRPCV1(context.Background(), addr, info.AuthHeader())
	if err != nil {
		log.Fatalf("Error connecting to Lotus %s", err)
	}

	return api, closer, nil
}

// SetupWallet Creating a new wallet and setting it as the default wallet.
// > SetupWallet creates a new wallet and returns it
func SetupWallet(dir string) (*wallet.LocalWallet, error) {
	kstore, err := keystore.OpenOrInitKeystore(dir)
	if err != nil {
		return nil, err
	}

	wallet, err := wallet.NewWallet(kstore)
	if err != nil {
		return nil, err
	}

	addrs, err := wallet.WalletList(context.TODO())
	if err != nil {
		return nil, err
	}

	if len(addrs) == 0 {
		_, err := wallet.WalletNew(context.TODO(), types.KTSecp256k1)
		if err != nil {
			return nil, err
		}
	}

	defaddr, err := wallet.GetDefault()

	fmt.Println("Wallet address is: ", defaddr)

	return wallet, nil
}

// GetPublicIP Getting the public IP of the node.
// > GetPublicIP() returns the public IP address of the machine it's running on
func GetPublicIP() (string, error) {
	resp, err := http.Get("https://ifconfig.me") // important to get the public ip if possible.
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
