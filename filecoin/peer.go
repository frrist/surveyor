package filecoin

import (
	"context"
	"fmt"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"

	"github.com/frrist/surveyor/core"
)

var log = logging.Logger("surveyor/filecoin")

type Config struct {
	Bootstrap         []peer.AddrInfo
	Datastore         datastore.Batching
	DHTProtocolPrefix protocol.ID
}

type Peer struct {
	core *core.Peer
	cfg  Config
}

func New(ctx context.Context, cfg Config, opts ...core.ConfigOpt) (*Peer, error) {
	if cfg.Bootstrap == nil {
		cfg.Bootstrap = MainnetPeers
	}

	if cfg.Datastore == nil {
		cfg.Datastore = dsync.MutexWrap(datastore.NewMapDatastore())
	}

	if cfg.DHTProtocolPrefix == "" {
		cfg.DHTProtocolPrefix = MainnetDHTPrefix
	}

	p, err := core.New(ctx, cfg.DHTProtocolPrefix, opts...)
	if err != nil {
		return nil, err
	}

	return &Peer{
		core: p,
		cfg:  cfg,
	}, nil
}

func (p *Peer) Bootstrap(ctx context.Context) error {
	return p.core.Bootstrap(ctx, p.cfg.Bootstrap)
}

type PeerAgentProtocols struct {
	Peer      peer.AddrInfo
	Agent     string
	Protocols []string
}

func (p *Peer) GetAllPeerAgentProtocols(ctx context.Context, whos []peer.AddrInfo, by time.Duration, workers int) chan *PeerAgentProtocols {
	out := make(chan *PeerAgentProtocols)
	wp := workerpool.New(workers)
	for _, who := range whos {
		wp.Submit(func() {
			found, err := p.GetPeerAgentProtocolsWithTimeout(ctx, who, by)
			if err != nil {
				return
			}
			out <- found
		})
	}
	go func() {
		wp.StopWait()
		close(out)
	}()
	return out
}

func (p *Peer) GetPeerAgentProtocolsWithTimeout(ctx context.Context, who peer.AddrInfo, by time.Duration) (*PeerAgentProtocols, error) {
	ctx, cancel := context.WithTimeout(ctx, by)
	defer cancel()
	found, err := p.GetPeerAgentProtocols(ctx, who)
	if err != nil {
		if err == context.DeadlineExceeded {
			log.Infow("deadline exceeded for peer", "error", err)
			return nil, err
		}
		log.Infow("connecting to peer failed", "error", err, "peer", who.String())
		return nil, err
	}
	return found, nil
}

func (p *Peer) GetPeerAgentProtocols(ctx context.Context, who peer.AddrInfo) (*PeerAgentProtocols, error) {
	if err := p.core.Host().Connect(ctx, who); err != nil {
		return nil, fmt.Errorf("connecting to peer %s: %w", who.Addrs, err)
	}

	protos, err := p.core.Host().Peerstore().GetProtocols(who.ID)
	if err != nil {
		return nil, err
	}

	agent, err := p.core.Host().Peerstore().Get(who.ID, "AgentVersion")
	if err != nil {
		return nil, err
	}
	return &PeerAgentProtocols{
		Peer:      who,
		Agent:     agent.(string),
		Protocols: protos,
	}, nil
}

// FindAllPeers queries the DHT for all peers in `whos` using `workers` goroutines, each peer has a find timeout of `by`.
// FindAllPeers will close the returned channel when all find operations have completed. Errors for finding peers are ignored.
func (p *Peer) FindAllPeers(ctx context.Context, whos []peer.ID, by time.Duration, workers int) chan peer.AddrInfo {
	out := make(chan peer.AddrInfo)
	wp := workerpool.New(workers)
	for _, who := range whos {
		wp.Submit(func() {
			found, err := p.FindPeerWithTimeout(ctx, who, by)
			if err != nil {
				return
			}
			out <- found
		})
	}
	go func() {
		wp.StopWait()
		close(out)
	}()
	return out
}

func (p *Peer) FindPeerWithTimeout(ctx context.Context, who peer.ID, by time.Duration) (peer.AddrInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, by)
	defer cancel()
	found, err := p.FindPeer(ctx, who)
	if err != nil {
		if err == context.DeadlineExceeded {
			log.Infow("deadline exceeded for peer", "error", err)
			return peer.AddrInfo{}, err
		}
		log.Infow("finding peer failed", "error", err, "peer", who.String())
		return peer.AddrInfo{}, err
	}
	return found, nil
}

func (p *Peer) FindPeer(ctx context.Context, who peer.ID) (peer.AddrInfo, error) {
	found, err := p.core.DHT().FindPeer(ctx, who)
	if err != nil {
		return peer.AddrInfo{}, err
	}
	return found, nil
}
