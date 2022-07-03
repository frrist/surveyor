package survey

import (
	"context"
	"time"

	"github.com/ip2location/ip2location-go"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	madns "github.com/multiformats/go-multiaddr-dns"
)

var ipLocationDB *ip2location.DB

func init() {
	db, err := ip2location.OpenDB("./IP2LOCATION-LITE-DB5.BIN")
	if err != nil {
		panic(err)
	}
	ipLocationDB = db
}

func NewPeerLocation(i *Inspector) *PeerLocation {
	return &PeerLocation{inspector: i}
}

type PeerLocation struct {
	inspector *Inspector
}

func (p *PeerLocation) Process(ctx context.Context, pai peer.AddrInfo) map[string]interface{} {
	for _, addr := range pai.Addrs {
		resolved, err := resolveMultiaddr(ctx, addr)
		if err != nil {
			continue
		}
		for _, r := range resolved {
			ipv4, ipv4Err := r.ValueForProtocol(multiaddr.P_IP4)
			if ipv4Err != nil {
				continue
			}

			all, err := ipLocationDB.Get_all(ipv4)
			if err != nil {
				return map[string]interface{}{
					"ID":    pai.ID.String(),
					"Error": err.Error(),
				}
			}
			return map[string]interface{}{
				"ID":        pai.ID,
				"Addrs":     MultiAddrsToString(pai.Addrs),
				"City":      all.City,
				"Country":   all.Country_long,
				"Latitude":  all.Latitude,
				"Longitude": all.Longitude,
			}
		}
	}
	return nil
}

func resolveMultiaddr(ctx context.Context, ma multiaddr.Multiaddr) ([]multiaddr.Multiaddr, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return madns.Resolve(ctx, ma)
}
