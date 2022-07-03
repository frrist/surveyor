package survey

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
)

func NewAgentProcessor(i *Inspector) *AgentProcessor {
	return &AgentProcessor{inspector: i}
}

type AgentProcessor struct {
	inspector *Inspector
}

func (a *AgentProcessor) Process(ctx context.Context, pai peer.AddrInfo) map[string]interface{} {
	if err := a.inspector.host.Connect(ctx, pai); err != nil {
		return map[string]interface{}{
			"ID":    pai.ID.String(),
			"Addr":  MultiAddrsToString(pai.Addrs),
			"Error": err.Error(),
		}
	}
	protos, err := a.inspector.host.Peerstore().GetProtocols(pai.ID)
	if err != nil {
		return map[string]interface{}{
			"ID":    pai.ID.String(),
			"Addr":  MultiAddrsToString(pai.Addrs),
			"Error": err.Error(),
		}
	}
	agent, err := a.inspector.host.Peerstore().Get(pai.ID, "AgentVersion")
	if err != nil {
		return map[string]interface{}{
			"ID":     pai.ID.String(),
			"Addr":   MultiAddrsToString(pai.Addrs),
			"Protos": protos,
			"Error":  err.Error(),
		}
	}

	return map[string]interface{}{
		"ID":        pai.ID.String(),
		"Addr":      MultiAddrsToString(pai.Addrs),
		"Agent":     agent.(string),
		"Protocols": protos,
	}
}
