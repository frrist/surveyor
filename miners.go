package main

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

var minerInfos = map[string]string{
	"12D3KooWAT134zy5QdXdnwLS5Jkrp9jpMmY8ao1asRN7kwNAdjoC": "/ip4/222.209.189.157/tcp/8066",
	"12D3KooWJNn4H5XVihwD4oEbouxJ7yyCAyLSsCB57XCHsLmcifRe": "/ip4/211.48.44.75/tcp/24001",
	"12D3KooWGicw4S44fb4hof72EY6HdQxtfZnbDcKcMKTGqMUd6QWe": "/ip4/114.86.142.162/tcp/24001",
	"12D3KooWRbJAznZsQVVi3FvUUzyqhZxBrXZDneEUwkfa1RnhR4S3": "/ip4/154.42.3.34/tcp/24001",
	"12D3KooWN8h1SduL4yEtmVo8mpuuYFD1pYTxhyrMqHhGA8SWvHZk": "/ip4/86.123.188.55/tcp/24001",
	"12D3KooWDKTnsWcGFhKJ9Gy4T8ULTuA73j1Yo1cEy8hrgxZTXGPa": "/ip4/47.108.179.232/tcp/24001",
	"12D3KooWQoUeEjqcxFf59Mfw2xmdMCNFm49Cwqy63mXKs7GtimKh": "/ip4/183.232.117.122/tcp/4001",
	"12D3KooWEZb2AAVTe1jSKnioqwcSzGSsTGD2dXpQT2qeYKfFZ7Tm": "/ip4/172.16.1.48/tcp/58418",
	"12D3KooWFWjjGANCcWq3vjEXydtYi8iwX1447zQA5fcaGiyw8VaK": "/ip4/202.160.106.54/tcp/28019",
	"12D3KooWMBrbCkYszgKZHq4WhdYLHPQc57JPJQNGyRNUx5oHtdvU": "/ip4/202.160.106.54/tcp/28018",
	"12D3KooWB4D2FSgP3wcadixg3RGq9PizQsyHyLuQYJRuajJa3kQN": "/ip4/202.160.106.54/tcp/28019",
	"12D3KooWNMc84yXh9AiQ2oZxrgmwmVs8R7pGwAqdwQDCQC3YC9Lw": "/ip4/202.160.106.54/tcp/28018",
	"12D3KooWG9QcL71Fd6zEYqdGEDxw5ZYGpBqccyXkKeoUDu6yawDN": "/ip4/103.127.248.194/tcp/10001",
	"12D3KooWFGQKFqY4vidD8h7i7zcQtSA4nhoPDyVbWfGk6G5AcP28": "/ip4/83.238.168.109/tcp/24002",
	"12D3KooWSqS2TSub9xKbjgX76bS19Pf1KWZDkL3CZaXFXn5JkdW7": "/ip4/185.32.160.239/tcp/42001",
	"12D3KooWC3VhMDKE44mDRLtE2j5D82tjseeV2Z7kpui44FvAaC1T": "/ip4/66.61.208.206/tcp/41372",
	"12D3KooWA3hqGbpsyF4L6vScuQFdi3nJxzGpuf5gT7jS2BAxkKuT": "/ip4/66.61.208.206/tcp/41372",
	"12D3KooWFMmYND7ArbQJ1PauzCcwt1q6ZFxE5E8jEEhToYF36EHr": "/ip4/66.61.208.206/tcp/41372",
	"12D3KooWAsv9bXDMCfdGEgembSsSZ9RrMcEgsfsdzHQaTPrSJP4A": "/ip4/66.61.208.206/tcp/41372",
	"12D3KooWQnyVm56fziAm2wXsLEddmbYDGH2sVsaBUMg9Q4QWKNCM": "/ip4/61.164.212.195/tcp/24001",
	"12D3KooWSyE7CspqGinUWdDqvZdVRJS3dnnJ4SXFnbu9MnW7DTNd": "/ip4/10.8.1.7/tcp/6666",
	"12D3KooWBvuSa1wDXmrxsrBUm3Tm26sT4RbKUwwENfspgYsmLjx8": "/ip4/208.97.197.18/tcp/11337",
	"12D3KooWFMPX3BkTCHGUcKhRktusAJQJxBZfJZUvE9sFcfK4r3oE": "/ip4/86.123.188.55/tcp/24001",
	"12D3KooWPnsqPeV6aKrncYxbyP3Bz8ZWZT4YNF9gBxUYD2DXncnE": "/ip4/66.6.126.70/tcp/15001",
	"12D3KooWBSyBYP26He59btRU8cbJnm3tEDPvuRQQnw58V8K6NDFM": "/ip4/182.18.86.72/tcp/46888",
	"12D3KooWBxFMAiY9ALptaw8rTg1tfMrWVxp6gbcCM7JrJb5zj4vR": "/ip4/211.181.56.218/tcp/24001",
	"12D3KooWDbkakgtEcUYnDopUDJWKFbSDKM1gQDz1jT6XLcdwHbWX": "/ip4/122.114.37.226/tcp/12000",
	"12D3KooWMYT1V3GAYLPH6aRe8MMjQbvQ5qgFQZWNPRSwB1RQPtRX": "/ip4/117.52.173.162/tcp/24003",
	"12D3KooWAESujTdiW3DpuUQyU1UGakiUmw2geR1QKZRRMgYxias1": "/ip4/218.78.187.146/tcp/21735",
	"12D3KooWK6UYYpVbqqPquP8eCSF6DUFGykgiJe2d9j3tii7S3a6N": "/ip4/10.7.3.2/tcp/1025",
	"12D3KooWM8jLc6ukRaiWkS8MeGR8LvrMK3qCeuEEikVPhgXhxgqJ": "/ip4/61.97.250.14/tcp/1235",
	"12D3KooWFmAudEtLLYRhTgUwDu6h1Wnf9JKfSfkC4YR6qJPBZeBt": "/ip4/67.212.85.196/tcp/10906",
	"12D3KooWJ6o3TjmVkbd9dfLfDewsHbonZY1ufxD25PqBRBXkdDgz": "/ip4/183.93.252.68/tcp/1235",
	"12D3KooWRvuctf6byqUb1dUXs9mgXHzWvJJp6JLEQSjfvsegNHna": "/ip4/121.178.127.135/tcp/24001",
	"12D3KooWG86gyFfD9co3tvuy7qHwcB637pJFwePX5ZSX9XfmFmTT": "/ip4/10.8.1.7/tcp/6666",
	"12D3KooWKQ3mbUyWmNzNeNQnom31myqZc9ZMxyW8My6oJHMVUoEF": "/ip4/10.8.1.9/tcp/6666",
	"12D3KooWSQQTfeUbgjWsAbZPuH6NZvxmwXBtEtJZaAHhGM1etzck": "/ip4/211.107.174.12/tcp/24001",
	"12D3KooWF3Digx6FXjRAwkJeFzVFaqTvtPrWXqqHusTPobb23pPj": "/ip4/54.219.44.132/tcp/8993",
	"12D3KooWNRMpZMVSYUQHFoLxaY3ZLr7qYsHD587gJuKHufcXs7Qd": "/ip4/18.144.19.252/tcp/7743",
	"12D3KooWQMduXUCDfc14u4PPd3wRwz6pxseNRELtwQN9wcpCv1ut": "/ip4/54.193.243.116/tcp/8654",
	"12D3KooWF9mXJA1tonrDLPyUNJMxJNLLpqQPntuocuULEqNCLRFT": "/ip4/18.190.24.1/tcp/9473",
	"12D3KooWSDbfUkotGQiT39sni5fuyPdGb3YzmJBtDRa3zZ9YbFUz": "/ip4/34.219.53.128/tcp/8083",
	"12D3KooWQoJfb2ueBDHN2ipEQqrrsbG9idB3jaCu1cmnvQDk21oi": "/ip4/52.27.153.66/tcp/8791",
	"12D3KooWLYaqBwdVkNh1e19AEpyveRdbtjwdTwtSR9rYEQNCJ3fN": "/ip4/35.90.15.192/tcp/8152",
}

func MinerPeerAddrInfo() []peer.AddrInfo {
	var out []peer.AddrInfo
	for id, addr := range minerInfos {
		pid, err := peer.Decode(id)
		if err != nil {
			panic(err)
		}
		maddr, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			panic(err)
		}
		out = append(out, peer.AddrInfo{
			ID:    pid,
			Addrs: []multiaddr.Multiaddr{maddr},
		})
	}
	return out
}
