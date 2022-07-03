package survey

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
)

type PeerAddrInfoProcessor interface {
	Process(ctx context.Context, pai peer.AddrInfo) map[string]interface{}
}
