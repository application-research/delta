package core

import (
	"context"
	c "delta/config"
	"delta/utils"
	"fmt"
	model "github.com/application-research/delta-db/db_models"
	"github.com/application-research/delta-db/messaging"
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
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	mdagipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-path/resolver"
	"github.com/labstack/gommon/log"
	trace2 "go.opencensus.io/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"sync"
	"syscall"
	"time"
)

// DeltaNode is a struct that contains a context, a node, an api, a database, a filecoin client, a config, a dispatcher, a
// delta tracer, a meta info, and a websocket broadcast.
// @property Context - The context of the node.
// @property Node - The Whypfs node that this DeltaNode is running on.
// @property Api - The URL of the Delta API
// @property DB - The database connection
// @property FilClient - This is the client that will be used to communicate with the Filecoin network.
// @property Config - This is the configuration for the Delta node.
// @property Dispatcher - This is the Dispatcher that is responsible for dispatching the messages to the appropriate
// handlers.
// @property DeltaTracer - This is a metrics tracer that is used to send metrics to the metrics server.
// @property MetaInfo - This is the metadata of the node. It contains the node's IP address, port, and other information.
// @property {WebsocketBroadcast} WebsocketBroadcast - This is a channel that is used to broadcast messages to the
// `WebsocketBroadcast` is a struct that contains three channels, one for each type of message that can be broadcasted.
// @property {ContentChannel} ContentChannel - This is the channel that will be used to send content to the client.
// @property {PieceCommitmentChannel} PieceCommitmentChannel - This is a channel that will be used to send piece
// commitments to the client.
// @property {ContentDealChannel} ContentDealChannel - This is a channel that will be used to send content deals to the
// client.
// websocket clients.
type DeltaNode struct {
	Context    context.Context
	Node       *whypfs.Node
	Api        url.URL
	DB         *gorm.DB
	FilClient  *fc.FilClient
	LotusApi   v1api.FullNode
	Config     *c.DeltaConfig
	Dispatcher *Dispatcher
	MetaInfo   *model.InstanceMeta

	DeltaEventEmitter  *DeltaEventEmitter
	DeltaMetricsTracer *DeltaMetricsTracer
}

type DeltaEventEmitter struct {
	WebsocketBroadcast WebsocketBroadcast
}

type DeltaMetricsTracer struct {
	Tracer            trace2.Tracer
	DeltaDataReporter *messaging.DeltaMetricsTracer
}

// WebsocketBroadcast `WebsocketBroadcast` is a struct that contains three channels, one for each type of message that can be broadcasted.
// @property {ContentChannel} ContentChannel - This is the channel that will be used to send content to the client.
// @property {PieceCommitmentChannel} PieceCommitmentChannel - This is a channel that will be used to send piece
// commitments to the client.
// @property {ContentDealChannel} ContentDealChannel - This is a channel that will be used to send content deals to the
// client.
type WebsocketBroadcast struct {
	ContentChannel         ContentChannel
	PieceCommitmentChannel PieceCommitmentChannel
	ContentDealChannel     ContentDealChannel
}

// A ContentChannel is a map of websocket connections to booleans and a channel of Content.
// @property Clients - A map of all the clients that are connected to the channel.
// @property Channel - This is the channel that will be used to send messages to all the clients.
// A ContentChannel is a map of websocket connections to booleans and a channel of Content.
// @property Clients - A map of all the clients that are connected to the channel.
// @property Channel - This is the channel that will be used to send messages to all the clients.
type ContentChannel struct {
	Clients map[*ClientChannel]bool
	Channel chan model.Content
}

type ClientChannel struct {
	Conn *websocket.Conn
	Id   string
}

// PieceCommitmentChannel `PieceCommitmentChannel` is a struct that contains a map of websocket connections and a channel of
// `model.PieceCommitment`s.
// @property Clients - A map of all the clients that are connected to the channel.
// @property Channel - This is the channel that will be used to send messages to all the clients.
type PieceCommitmentChannel struct {
	Clients map[*websocket.Conn]bool
	Channel chan model.PieceCommitment
}

// ContentDealChannel `ContentDealChannel` is a struct with a map of `websocket.Conn` pointers and a `chan model.ContentDeal` channel.
// @property Clients - A map of all the clients connected to the channel.
// @property Channel - This is the channel that will be used to send messages to the clients.
type ContentDealChannel struct {
	Clients map[*websocket.Conn]bool
	Channel chan model.ContentDeal
}

// LocalWallet `LocalWallet` is a struct that contains a map of `address.Address` to `key.Key` and a `types.KeyStore` and a
// `sync.Mutex`.
// @property keys - a map of addresses to keys.
// @property keystore - This is the keystore that the wallet uses to store the keys.
// @property lk - A mutex to prevent concurrent access to the wallet.
type LocalWallet struct {
	keys     map[address.Address]*key.Key
	keystore types.KeyStore
	lk       sync.Mutex
}

// GatewayHandler It's a struct that holds a blockstore, a DAGService, a Resolver, and a Node.
// @property bs - The blockstore is the storage layer for the IPFS node. It stores the raw data of the blocks.
// @property dserv - The DAGService is the interface that allows us to interact with the IPFS DAG.
// @property resolver - This is the resolver that will be used to resolve paths to IPFS objects.
// @property node - The Whypfs node that will be used to serve the requests.
type GatewayHandler struct {
	bs blockstore.Blockstore
	// `NewLightNodeParams` is a struct with three fields: `Repo`, `DefaultWalletDir`, and `Config`.
	// @property {string} Repo - The path to the directory where the node will store its data.
	// @property {string} DefaultWalletDir - The default directory where the wallet will be stored.
	// @property Config - This is the configuration object that is used to configure the node.
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

	filclient, err := fc.NewClient(whypfsPeer.Host, api, wallet, walletAddr, whypfsPeer.Blockstore, whypfsPeer.Datastore, whypfsPeer.Config.DatastoreDir.Directory)

	if err != nil {
		panic(err)
	}

	// job dispatcher
	dispatcher := CreateNewDispatcher()

	// delta metrics tracer
	dataTracer := messaging.NewDeltaMetricsTracer()

	openTelemetryTracerProvider := trace.NewTracerProvider(trace.WithSampler(trace.AlwaysSample()))
	defer openTelemetryTracerProvider.Shutdown(context.Background())

	// Register the tracer provider with the global tracer
	otel.SetTracerProvider(openTelemetryTracerProvider)

	// Create a new span
	//tracer := otel.Tracer("example")
	// create the global light node.
	return &DeltaNode{
		Node:       whypfsPeer,
		DB:         db,
		FilClient:  filclient,
		Dispatcher: dispatcher,
		LotusApi:   api,
		Config:     repo.Config,
		DeltaMetricsTracer: &DeltaMetricsTracer{
			DeltaDataReporter: dataTracer,
		},
	}, nil
}

// LotusConnection It takes a string that contains the Lotus full node API address and returns a `v1api.FullNode` interface, a
// `jsonrpc.ClientCloser` interface, and an error
// It takes a string that contains the Lotus full node API address and returns a `v1api.FullNode` interface, a
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

// Get the hostname of the current machine, or return 'unknown' if there's an error.
func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

// ScanHostComputeResources Setting the global node meta.
// > This function sets the global node metadata for the given node
func ScanHostComputeResources(ln *DeltaNode, repo string) *model.InstanceMeta {

	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	totalMemory := memStats.Sys
	totalMemory80 := totalMemory * 90 / 100

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Total memory: %v bytes\n", m.Alloc)
	fmt.Printf("Total system memory: %v bytes\n", m.Sys)
	fmt.Printf("Total heap memory: %v bytes\n", m.HeapSys)
	fmt.Printf("Heap in use: %v bytes\n", m.HeapInuse)
	fmt.Printf("Stack in use: %v bytes\n", m.StackInuse)
	// get the 80% of the total disk usage
	var stat syscall.Statfs_t
	syscall.Statfs(repo, &stat) // blockstore size
	totalStorage := stat.Blocks * uint64(stat.Bsize)
	totalStorage90 := totalStorage * 90 / 100

	// set the number of CPUs
	numCPU := runtime.NumCPU()
	fmt.Printf("Total number of CPUs: %d\n", numCPU)
	fmt.Printf("Number of CPUs that this Delta will use: %d\n", numCPU/(1200/1000))
	fmt.Println(utils.Purple + "Note: Delta instance proactively recalculate resources to use based on the current load." + utils.Reset)
	runtime.GOMAXPROCS(numCPU / (1200 / 1000))

	// delete all data from the instance meta table
	//ln.DB.Model(&model.InstanceMeta{}).Delete(&model.InstanceMeta{}, "id > ?", 0)
	// re-create
	ip, err := GetPublicIP()
	if err != nil {
		fmt.Println("Error getting public IP:", err)
	}

	// if there's already an existing record, update that record
	var instanceMeta model.InstanceMeta
	ln.DB.Model(&model.InstanceMeta{}).Where("id > ?", 0).Order("created_at desc").First(&instanceMeta)

	if instanceMeta.InstanceUuid == "" {
		instanceMeta.InstanceUuid = uuid.New().String()
	}

	instanceMeta = model.InstanceMeta{
		MemoryLimit:      totalMemory80,  // 80%
		StorageLimit:     totalStorage90, // 90%
		InstanceUuid:     instanceMeta.InstanceUuid,
		NumberOfCpus:     uint64(numCPU),
		InstanceNodeName: ln.Config.Node.Name,
		PublicIp:         ip,
		OSDetails:        runtime.GOARCH + " " + runtime.GOOS,
		StorageInBytes:   totalStorage,
		BytesPerCpu:      11000000000,
		SystemMemory:     totalMemory,
		HeapMemory:       m.HeapSys,
		HeapInUse:        m.HeapInuse,
		StackInUse:       m.StackInuse,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		InstanceStart:    time.Now(),
	}

	if instanceMeta.ID > 0 {
		ln.DB.Model(&model.InstanceMeta{}).Save(&instanceMeta)
		ln.Config.Node.InstanceUuid = instanceMeta.InstanceUuid
		ln.MetaInfo = &instanceMeta

	} else {

		ln.DB.Model(&model.InstanceMeta{}).Create(&instanceMeta)
		ln.Config.Node.InstanceUuid = instanceMeta.InstanceUuid
		ln.MetaInfo = &instanceMeta
	}

	return &instanceMeta

}

// CleanUpContentAndPieceComm It updates the status of all the content and piece commitments that were in the process of being transferred or computed
// to failed
func CleanUpContentAndPieceComm(ln *DeltaNode) {

	// if the transfer was started upon restart, then we need to update the status to failed
	ln.DB.Transaction(func(tx *gorm.DB) error {

		rowsAffected := tx.Model(&model.Content{}).Where("status in (?,?)", utils.DEAL_STATUS_TRANSFER_STARTED, utils.CONTENT_PIECE_COMPUTING).Updates(
			model.Content{
				Status:      utils.DEAL_STATUS_TRANSFER_FAILED,
				UpdatedAt:   time.Now(),
				LastMessage: "Transfer failed due to node restart",
			}).RowsAffected

		fmt.Println("Number of rows cleaned up: " + fmt.Sprint(rowsAffected))
		return nil
	})
}
