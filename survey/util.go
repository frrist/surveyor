package survey

import "github.com/multiformats/go-multiaddr"

func MultiAddrsToString(maddrs []multiaddr.Multiaddr) []string {
	var out []string
	for _, addr := range maddrs {
		out = append(out, addr.String())
	}
	return out
}
