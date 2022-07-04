package ipfs

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

var DefaultBootstrapAddresses = []string{
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",
	"/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",      // mars.i.ipfs.io
	"/ip4/104.131.131.82/udp/4001/quic/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ", // mars.i.ipfs.io
}

var IPFSPeers []peer.AddrInfo

func initIPFSPeers() {
	mainMaddrs := make([]multiaddr.Multiaddr, len(DefaultBootstrapAddresses))

	for i, p := range DefaultBootstrapAddresses {
		maddr, err := multiaddr.NewMultiaddr(p)
		if err != nil {
			panic(err)
		}
		mainMaddrs[i] = maddr
	}
	var err error
	IPFSPeers, err = peer.AddrInfosFromP2pAddrs(mainMaddrs...)
	if err != nil {
		panic(err)
	}
}
