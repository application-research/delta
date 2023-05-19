package config

import (
	"github.com/caarlos0/env/v6"
	logging "github.com/ipfs/go-log/v2"
	"github.com/joho/godotenv"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

var (
	log                       = logging.Logger("config")
	defaultTestBootstrapPeers []multiaddr.Multiaddr
)

type DeltaConfig struct {
	Node struct {
		Name         string `env:"NODE_NAME" envDefault:"delta-deal-maker"`
		Description  string `env:"NODE_DESCRIPTION" envDefault:"delta-deal-maker"`
		Type         string `env:"NODE_TYPE" envDefault:"delta-deal-maker"`
		InstanceUuid string `env:"INSTANCE_UUID"`
		KeepCopies   bool   `env:"KEEP_COPIES" envDefault:"false"`
	}

	Dispatcher struct {
		MaxCleanupWorkers int `env:"MAX_CLEANUP_WORKERS" envDefault:"1500"`
	}

	Common struct {
		Mode                  string `env:"MODE" envDefault:"standalone"`
		DBDSN                 string `env:"DB_DSN" envDefault:"delta.db"`
		EnableWebsocket       bool   `env:"ENABLE_WEBSOCKET" envDefault:"false"`
		CommpMode             string `env:"COMMP_MODE" envDefault:"fast"` // option "filboost"
		StatsCollection       bool   `env:"STATS_COLLECTION" envDefault:"true"`
		Commit                string `env:"COMMIT"`
		Version               string `env:"VERSION"`
		MaxReplicationFactor  int    `env:"MAX_REPLICATION_FACTOR" envDefault:"6"`
		MaxAutoRetry          int    `env:"MAX_AUTO_RETRY" envDefault:"3"`
		EnableInclusionProofs bool   `env:"ENABLE_INCLUSION_PROOFS" envDefault:"false"`
	}

	// configurable via env vars
	ExternalApis struct {
		LotusApi       string `env:"LOTUS_API" envDefault:"http://api.chain.love"`
		AuthSvcApi     string `env:"AUTH_SVC_API" envDefault:"https://auth.estuary.tech"`
		SpThrottlerApi string `env:"SP_THROTTLER_SVC_API" envDefault:"https://sp-throttler.delta.store"`
		SpSelectionApi string `env:"SP_SELECTION_SVC_API" envDefault:"https://sp-select.delta.store/api/providers"`
		DealStatusApi  string `env:"DEAL_STATUS_API" envDefault:"https://deal-status.estuary.tech"`
	}

	Standalone struct {
		APIKey string `env:"DELTA_AUTH" envDefault:""`
	}
}

func InitConfig() DeltaConfig {
	godotenv.Load() // load from environment OR .env file if it exists
	var cfg DeltaConfig

	if err := env.Parse(&cfg); err != nil {
		log.Fatal("error parsing config: %+v\n", err)
	}

	log.Debug("config parsed successfully")

	return cfg
}

// BootstrapEstuaryPeers Creating a list of multiaddresses that are used to bootstrap the network.
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
