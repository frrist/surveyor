package survey

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
)

func NewDHTProcessor(i *Inspector) *DHTProcessor {
	return &DHTProcessor{inspector: i}
}

type DHTProcessor struct {
	inspector *Inspector
}

func (d *DHTProcessor) Process(ctx context.Context, pai peer.AddrInfo) map[string]interface{} {
	found, err := d.inspector.dht.FindPeer(ctx, pai.ID)
	if err == nil {
		return map[string]interface{}{
			"ID":         pai.ID.String(),
			"Addr":       MultiAddrsToString(pai.Addrs),
			"Found":      true,
			"FoundAddrs": MultiAddrsToString(found.Addrs),
		}
	}
	return map[string]interface{}{
		"ID":         pai.ID.String(),
		"Addr":       MultiAddrsToString(pai.Addrs),
		"Found":      false,
		"FoundAddrs": nil,
		"Error":      err.Error(),
	}
}
