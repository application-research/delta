package core

import (
	"context"
	"fmt"
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
	"github.com/ipfs/go-blockservice"
	bsfetcher "github.com/ipfs/go-fetcher/impl/blockservice"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	mdagipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-path/resolver"
	"github.com/ipfs/go-unixfsnode"
	dagpb "github.com/ipld/go-codec-dagpb"
	"github.com/ipld/go-ipld-prime"
	ipldbasicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/ipld/go-ipld-prime/schema"
	"github.com/labstack/gommon/log"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/urfave/cli"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

type LightNode struct {
	Node      *whypfs.Node
	Api       url.URL
	Gw        *GatewayHandler
	DB        *gorm.DB
	Wallet    LocalWallet
	Filclient *fc.FilClient
	Config    *Configuration
}

type LocalWallet struct {
	keys     map[address.Address]*key.Key
	keystore types.KeyStore

	lk sync.Mutex
}

type Configuration struct {
	APINodeAddress string
}

type GatewayHandler struct {
	bs       blockstore.Blockstore
	dserv    mdagipld.DAGService
	resolver resolver.Resolver
	node     *whypfs.Node
}

var defaultTestBootstrapPeers []multiaddr.Multiaddr

// Creating a list of multiaddresses that are used to bootstrap the network.
func BootstrapEstuaryPeers() []peer.AddrInfo {

	for _, s := range []string{
		"/ip4/145.40.90.135/tcp/6746/p2p/12D3KooWNTiHg8eQsTRx8XV7TiJbq3379EgwG6Mo3V3MdwAfThsx",
		"/ip4/139.178.68.217/tcp/6744/p2p/12D3KooWCVXs8P7iq6ao4XhfAmKWrEeuKFWCJgqe9jGDMTqHYBjw",
		"/ip4/147.75.49.71/tcp/6745/p2p/12D3KooWGBWx9gyUFTVQcKMTenQMSyE2ad9m7c9fpjS4NMjoDien",
		"/ip4/147.75.86.255/tcp/6745/p2p/12D3KooWFrnuj5o3tx4fGD2ZVJRyDqTdzGnU3XYXmBbWbc8Hs8Nd",
		"/ip4/3.134.223.177/tcp/6745/p2p/12D3KooWN8vAoGd6eurUSidcpLYguQiGZwt4eVgDvbgaS7kiGTup",
		"/ip4/35.74.45.12/udp/6746/quic/p2p/12D3KooWLV128pddyvoG6NBvoZw7sSrgpMTPtjnpu3mSmENqhtL7",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
		"/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",
	} {
		ma, err := multiaddr.NewMultiaddr(s)
		if err != nil {
			panic(err)
		}
		defaultTestBootstrapPeers = append(defaultTestBootstrapPeers, ma)
	}

	peers, _ := peer.AddrInfosFromP2pAddrs(defaultTestBootstrapPeers...)
	return peers
}

// Add a config to enable gateway or not.
// Add a config to enable content, bucket, commp, replication verifier processor
func NewLightNode(ctx context.Context) (*LightNode, error) {
	//	database
	db, err := OpenDatabase()
	publicIp, err := GetPublicIP()
	newConfig := &whypfs.Config{
		ListenAddrs: []string{
			"/ip4/" + publicIp + "/tcp/6745",
			"/ip4/0.0.0.0/tcp/6745",
		},
		AnnounceAddrs: []string{
			"/ip4/" + publicIp + "/tcp/6745",
			"/ip4/0.0.0.0/tcp/6745",
		},
	}
	params := whypfs.NewNodeParams{
		Ctx:       context.Background(),
		Datastore: whypfs.NewInMemoryDatastore(),
	}
	// node
	params.Config = params.ConfigurationBuilder(newConfig)
	whypfsPeer, err := whypfs.NewNode(params)

	if err != nil {
		panic(err)
	}

	whypfsPeer.BootstrapPeers(BootstrapEstuaryPeers())

	//	Filclient
	api, _, err := LotusConnection("https://api.node.glif.io")
	addr, err := api.WalletDefaultAddress(ctx)

	fmt.Println(addr)
	//wallet := &wallet.LocalWallet{}
	wallet, err := SetupWallet("./wallet")
	walletAddr, err := wallet.GetDefault()
	if err != nil {
		panic(err)
	}

	fc, err := fc.NewClient(whypfsPeer.Host, api, wallet, walletAddr, whypfsPeer.Blockstore, whypfsPeer.Datastore, whypfsPeer.Config.DatastoreDir.Directory)

	if err != nil {
		panic(err)
	}

	// create the global light node.
	return &LightNode{
		Node:      whypfsPeer,
		DB:        db,
		Filclient: fc,
	}, nil
}

func NewFullLightNode(ctx *cli.Context) (*LightNode, error) {

	// database connection
	db, err := OpenDatabase()

	// node
	whypfsPeer, err := whypfs.NewNode(whypfs.NewNodeParams{
		Ctx:       context.Background(),
		Datastore: whypfs.NewInMemoryDatastore(),
	})
	if err != nil {
		panic(err)
	}

	whypfsPeer.BootstrapPeers(BootstrapEstuaryPeers())

	// Filclient
	api, _, err := LotusConnection("http://api.chain.love")
	addr, err := api.WalletDefaultAddress(context.Background())
	wallet := &wallet.LocalWallet{}

	fmt.Println(whypfsPeer.Host.ID().String())
	fc, err := fc.NewClient(whypfsPeer.Host, api, wallet, addr, whypfsPeer.Blockstore, whypfsPeer.Datastore, whypfsPeer.Config.DatastoreDir.Directory)

	amt, err := types.ParseFIL("0.5")

	fc.LockMarketFunds(context.Background(), amt)
	// gateway
	gw, err := NewGatewayHandler(whypfsPeer)

	if err != nil {
		return nil, err
	}

	// create the global light node.
	return &LightNode{
		Node:      whypfsPeer,
		Gw:        gw,
		DB:        db,
		Filclient: fc,
	}, nil

}

func NewGatewayHandler(node *whypfs.Node) (*GatewayHandler, error) {

	bsvc := blockservice.New(node.Blockstore, nil)
	ipldFetcher := bsfetcher.NewFetcherConfig(bsvc)

	ipldFetcher.PrototypeChooser = dagpb.AddSupportToChooser(func(lnk ipld.Link, lnkCtx ipld.LinkContext) (ipld.NodePrototype, error) {
		if tlnkNd, ok := lnkCtx.LinkNode.(schema.TypedLinkNode); ok {
			return tlnkNd.LinkTargetNodePrototype(), nil
		}
		return ipldbasicnode.Prototype.Any, nil
	})

	resolver := resolver.NewBasicResolver(ipldFetcher.WithReifier(unixfsnode.Reify))
	return &GatewayHandler{
		bs:       node.Blockstore,
		dserv:    merkledag.NewDAGService(bsvc),
		resolver: resolver,
		node:     node,
	}, nil
}

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
