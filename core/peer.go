package core

import (
	"context"
	"fmt"
	"sync"

	"github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	noise "github.com/libp2p/go-libp2p-noise"
	libp2ptls "github.com/libp2p/go-libp2p-tls"
	"github.com/libp2p/go-libp2p/config"
)

var log = logging.Logger("surveyor/core")

type ConfigOpt func(c *Config)

func WithLibp2pOptions(opts []libp2p.Option) ConfigOpt {
	return func(c *Config) {
		c.libp2pOpts = opts
	}
}

func WithKeyPair(pk crypto.PubKey, sk crypto.PrivKey) ConfigOpt {
	return func(c *Config) {
		c.privateKey = sk
		c.publicKey = pk
	}
}

func WithDatastore(ds datastore.Batching) ConfigOpt {
	return func(c *Config) {
		c.datastore = ds
	}
}

type Config struct {
	libp2pOpts []libp2p.Option
	publicKey  crypto.PubKey
	privateKey crypto.PrivKey
	datastore  datastore.Batching
}

type Peer struct {
	dht  *dht.IpfsDHT
	host host.Host
	cfg  *Config
}

func New(routing config.RoutingC, opts ...ConfigOpt) (*Peer, error) {
	if routing == nil {
		return nil, fmt.Errorf("routing required")
	}
	cfg := new(Config)
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.privateKey == nil || cfg.publicKey == nil {
		var err error
		cfg.privateKey, cfg.publicKey, err = crypto.GenerateKeyPair(crypto.Ed25519, -1)
		if err != nil {
			return nil, err
		}
	}

	if cfg.datastore == nil {
		cfg.datastore = dsync.MutexWrap(datastore.NewMapDatastore())
	}

	var ddht *dht.IpfsDHT
	if cfg.libp2pOpts == nil {
		cfg.libp2pOpts = []libp2p.Option{
			libp2p.Identity(cfg.privateKey),
			libp2p.Security(libp2ptls.ID, libp2ptls.New),
			libp2p.Security(noise.ID, noise.New),
			libp2p.DefaultTransports,
			libp2p.Routing(routing),
		}
	}

	h, err := libp2p.New(cfg.libp2pOpts...)
	if err != nil {
		return nil, err
	}

	return &Peer{
		dht:  ddht,
		host: h,
		cfg:  cfg,
	}, nil
}

func (p *Peer) Bootstrap(ctx context.Context, bootstrap []peer.AddrInfo) error {
	connected := make(chan struct{})

	var wg sync.WaitGroup
	for _, pinfo := range bootstrap {
		wg.Add(1)
		go func(pinfo peer.AddrInfo) {
			defer wg.Done()
			err := p.host.Connect(ctx, pinfo)
			if err != nil {
				log.Warn(err)
				return
			}
			log.Infow("Connected", "peer", pinfo.ID)
			connected <- struct{}{}
		}(pinfo)
	}

	go func() {
		wg.Wait()
		close(connected)
	}()

	i := 0
	for range connected {
		i++
	}
	if nPeers := len(bootstrap); i < nPeers/2 {
		log.Warnf("only connected to %d bootstrap peers out of %d", i, nPeers)
	}

	err := p.dht.Bootstrap(ctx)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (p *Peer) Host() host.Host {
	return p.host
}

func (p *Peer) DHT() *dht.IpfsDHT {
	return p.dht
}
