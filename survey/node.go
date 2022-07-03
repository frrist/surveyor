package survey

import (
	"context"
	crand "crypto/rand"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	ipfslite "github.com/hsanjuan/ipfs-lite"
	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	dht2 "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-kad-dht/dual"
	"github.com/multiformats/go-multiaddr"
)

var log = logging.Logger("survey/node")

var Boostrappers = []string{
	"/dns4/bootstrap-0.mainnet.filops.net/tcp/1347/p2p/12D3KooWCVe8MmsEMes2FzgTpt9fXtmCY7wrq91GRiaC8PHSCCBj",
	"/dns4/bootstrap-1.mainnet.filops.net/tcp/1347/p2p/12D3KooWCwevHg1yLCvktf2nvLu7L9894mcrJR4MsBCcm4syShVc",
	"/dns4/bootstrap-2.mainnet.filops.net/tcp/1347/p2p/12D3KooWEWVwHGn2yR36gKLozmb4YjDJGerotAPGxmdWZx2nxMC4",
	"/dns4/bootstrap-3.mainnet.filops.net/tcp/1347/p2p/12D3KooWKhgq8c7NQ9iGjbyK7v7phXvG6492HQfiDaGHLHLQjk7R",
	"/dns4/bootstrap-4.mainnet.filops.net/tcp/1347/p2p/12D3KooWL6PsFNPhYftrJzGgF5U18hFoaVhfGk7xwzD8yVrHJ3Uc",
	"/dns4/bootstrap-5.mainnet.filops.net/tcp/1347/p2p/12D3KooWLFynvDQiUpXoHroV1YxKHhPJgysQGH2k3ZGwtWzR4dFH",
	"/dns4/bootstrap-6.mainnet.filops.net/tcp/1347/p2p/12D3KooWP5MwCiqdMETF9ub1P3MbCvQCcfconnYHbWg6sUJcDRQQ",
	"/dns4/bootstrap-7.mainnet.filops.net/tcp/1347/p2p/12D3KooWRs3aY1p3juFjPy8gPN95PEQChm2QKGUCAdcDCC4EBMKf",
	"/dns4/bootstrap-8.mainnet.filops.net/tcp/1347/p2p/12D3KooWScFR7385LTyR4zU1bYdzSiiAb5rnNABfVahPvVSzyTkR",

	"/dns4/lotus-bootstrap.ipfsforce.com/tcp/41778/p2p/12D3KooWGhufNmZHF3sv48aQeS13ng5XVJZ9E6qy2Ms4VzqeUsHk",

	"/dns4/bootstrap-0.starpool.in/tcp/12757/p2p/12D3KooWGHpBMeZbestVEWkfdnC9u7p6uFHXL1n7m1ZBqsEmiUzz",
	"/dns4/bootstrap-1.starpool.in/tcp/12757/p2p/12D3KooWQZrGH1PxSNZPum99M1zNvjNFM33d1AAu5DcvdHptuU7u",

	"/dns4/node.glif.io/tcp/1235/p2p/12D3KooWBF8cpp65hp2u9LK5mh19x67ftAam84z9LsfaquTDSBpt",

	"/dns4/bootstrap-0.ipfsmain.cn/tcp/34721/p2p/12D3KooWQnwEGNqcM2nAcPtRR9rAX8Hrg4k9kJLCHoTR5chJfz6d",
	"/dns4/bootstrap-1.ipfsmain.cn/tcp/34723/p2p/12D3KooWMKxMkD5DMpSWsW7dBddKxKT7L2GgbNuckz9otxvkvByP",

	"/dns4/bootstarp-0.1475.io/tcp/61256/p2p/12D3KooWQjaNmbz9b1XmheQB3RWsRjKSzuRLfjeiDZHyX7Y5RcBr",
}

var MainnetBootstrapPeers []peer.AddrInfo

func init() {
	for _, bs := range Boostrappers {
		pai, err := ParsePeerString(bs)
		if err != nil {
			panic(err)
		}
		MainnetBootstrapPeers = append(MainnetBootstrapPeers, *pai)
	}
}

type Config struct {
	Libp2pKeyFile string
}

type Inspector struct {
	host host.Host
	dht  *dual.DHT
	lite *ipfslite.Peer
}

func New(ctx context.Context, cfg *Config) (*Inspector, error) {
	pk, err := loadOrInitPeerKey(cfg.Libp2pKeyFile)
	if err != nil {
		return nil, err
	}
	ds := ipfslite.NewInMemoryDatastore()

	h, dht, err := ipfslite.SetupLibp2p(ctx, pk, nil, nil, ds,
		[]dht2.Option{
			dht2.ProtocolPrefix(protocol.ID("/fil/kad/testnetnet")),
			dht2.QueryFilter(dht2.PublicQueryFilter),
			dht2.RoutingTableFilter(dht2.PublicRoutingTableFilter),
		},
		append(ipfslite.Libp2pOptionsExtra)...,
	)
	if err != nil {
		return nil, err
	}

	lite, err := ipfslite.New(ctx, ds, h, dht, nil)
	if err != nil {
		return nil, err
	}

	lite.Bootstrap(MainnetBootstrapPeers)

	return &Inspector{
		lite: lite,
		host: h,
		dht:  dht,
	}, nil
}

type Result struct {
	ID   peer.ID
	Data map[string]interface{}
}

func (i *Inspector) Run(ctx context.Context, todo []peer.AddrInfo, result chan *Result, processors ...PeerAddrInfoProcessor) {
	var wg sync.WaitGroup
	for _, t := range todo {
		t := t
		for _, p := range processors {
			p := p
			wg.Add(1)
			go func(p peer.AddrInfo, proc PeerAddrInfoProcessor) {
				defer wg.Done()
				res := proc.Process(ctx, t)
				result <- &Result{
					ID:   t.ID,
					Data: res,
				}
			}(t, p)
		}
	}
	wg.Wait()
	close(result)
}

func loadOrInitPeerKey(kf string) (crypto.PrivKey, error) {
	data, err := ioutil.ReadFile(kf)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		k, _, err := crypto.GenerateEd25519Key(crand.Reader)
		if err != nil {
			return nil, err
		}

		data, err := crypto.MarshalPrivateKey(k)
		if err != nil {
			return nil, err
		}

		if err := ioutil.WriteFile(kf, data, 0600); err != nil {
			return nil, err
		}

		return k, nil
	}
	return crypto.UnmarshalPrivateKey(data)
}
func ParsePeerString(text string) (*peer.AddrInfo, error) {
	// Multiaddr
	if strings.HasPrefix(text, "/") {
		maddr, err := multiaddr.NewMultiaddr(text)
		if err != nil {
			return nil, err
		}
		ainfo, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			return nil, err
		}
		return ainfo, nil
	}
	return nil, peer.ErrInvalidAddr
}
