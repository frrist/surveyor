package ipfs

import (
	"context"

	ipfslite "github.com/hsanjuan/ipfs-lite"
	"github.com/ipfs/go-bitswap"
	"github.com/ipfs/go-bitswap/network"
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	exchange "github.com/ipfs/go-ipfs-exchange-interface"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-ipns"
	logging "github.com/ipfs/go-log/v2"
	"github.com/ipfs/go-merkledag"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	record "github.com/libp2p/go-libp2p-record"
	"github.com/libp2p/go-libp2p/config"

	"github.com/frrist/surveyor/core"
)

var log = logging.Logger("surveyor/ipfs")

type Config struct {
	Bootstrap []peer.AddrInfo
	Datastore datastore.Batching
	Routing   config.RoutingC
}

type Peer struct {
	core *core.Peer
	cfg  *Config

	exch            exchange.Interface
	bstore          blockstore.Blockstore
	bserv           blockservice.BlockService
	ipld.DAGService // become a DAG service
}

func New(ctx context.Context, cfg *Config, opts ...core.ConfigOpt) (*Peer, error) {
	if cfg.Bootstrap == nil {
		cfg.Bootstrap = ipfslite.DefaultBootstrapPeers()
	}

	if cfg.Datastore == nil {
		cfg.Datastore = dsync.MutexWrap(datastore.NewMapDatastore())
	}

	if cfg.Routing == nil {
		cfg.Routing = func(h host.Host) (routing.PeerRouting, error) {
			ddht, err := dht.New(
				ctx,
				h,
				dht.Datastore(cfg.Datastore),
				dht.NamespacedValidator("pk", record.PublicKeyValidator{}),
				dht.NamespacedValidator("ipns", ipns.Validator{KeyBook: h.Peerstore()}),
				dht.Concurrency(50),
				dht.Mode(dht.ModeClient),
			)
			return ddht, err
		}
	}

	p, err := core.New(cfg.Routing, opts...)
	if err != nil {
		return nil, err
	}

	ipfsPeer := &Peer{
		core: p,
		cfg:  cfg,
	}
	if err := ipfsPeer.init(ctx); err != nil {
		return nil, err
	}
	return ipfsPeer, nil
}

func (p *Peer) init(ctx context.Context) error {
	if err := p.setupBlockstore(ctx); err != nil {
		return err
	}
	if err := p.setupBlockService(ctx); err != nil {
		return err
	}
	if err := p.setupDAGService(); err != nil {
		_ = p.bserv.Close()
		return err
	}
	return nil
}

func (p *Peer) Bootstrap(ctx context.Context) error {
	return p.core.Bootstrap(ctx, p.cfg.Bootstrap)
}

// Session returns a session-based NodeGetter.
func (p *Peer) Session(ctx context.Context) ipld.NodeGetter {
	ng := merkledag.NewSession(ctx, p.DAGService)
	if ng == p.DAGService {
		log.Warn("DAGService does not support sessions")
	}
	return ng
}

func (p *Peer) setupBlockstore(ctx context.Context) error {
	bs := blockstore.NewBlockstore(p.cfg.Datastore)
	bs = blockstore.NewIdStore(bs)
	cachedbs, err := blockstore.CachedBlockstore(ctx, bs, blockstore.DefaultCacheOpts())
	if err != nil {
		return err
	}
	p.bstore = cachedbs
	return nil
}

func (p *Peer) setupBlockService(ctx context.Context) error {
	bswapnet := network.NewFromIpfsHost(p.core.Host(), p.core.DHT())
	bswap := bitswap.New(ctx, bswapnet, p.bstore)
	p.bserv = blockservice.New(p.bstore, bswap)
	p.exch = bswap
	return nil
}

func (p *Peer) setupDAGService() error {
	p.DAGService = merkledag.NewDAGService(p.bserv)
	return nil
}
