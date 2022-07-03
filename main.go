package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	logging "github.com/ipfs/go-log/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"

	"github.com/frrist/surveyor/survey"
)

var MainnetBoostrappers = []string{
	"/dns4/bootstrap-0.mainnet.filops.net/tcp/1347/p2p/12D3KooWCVe8MmsEMes2FzgTpt9fXtmCY7wrq91GRiaC8PHSCCBj",
	"/dns4/bootstrap-1.mainnet.filops.net/tcp/1347/p2p/12D3KooWCwevHg1yLCvktf2nvLu7L9894mcrJR4MsBCcm4syShVc",
	"/dns4/bootstrap-2.mainnet.filops.net/tcp/1347/p2p/12D3KooWEWVwHGn2yR36gKLozmb4YjDJGerotAPGxmdWZx2nxMC4",
	"/dns4/bootstrap-3.mainnet.filops.net/tcp/1347/p2p/12D3KooWKhgq8c7NQ9iGjbyK7v7phXvG6492HQfiDaGHLHLQjk7R",
	"/dns4/bootstrap-4.mainnet.filops.net/tcp/1347/p2p/12D3KooWL6PsFNPhYftrJzGgF5U18hFoaVhfGk7xwzD8yVrHJ3Uc",
	"/dns4/bootstrap-5.mainnet.filops.net/tcp/1347/p2p/12D3KooWLFynvDQiUpXoHroV1YxKHhPJgysQGH2k3ZGwtWzR4dFH",
	"/dns4/bootstrap-6.mainnet.filops.net/tcp/1347/p2p/12D3KooWP5MwCiqdMETF9ub1P3MbCvQCcfconnYHbWg6sUJcDRQQ",
	"/dns4/bootstrap-7.mainnet.filops.net/tcp/1347/p2p/12D3KooWRs3aY1p3juFjPy8gPN95PEQChm2QKGUCAdcDCC4EBMKf",
	"/dns4/bootstrap-8.mainnet.filops.net/tcp/1347/p2p/12D3KooWScFR7385LTyR4zU1bYdzSiiAb5rnNABfVahPvVSzyTkR",

	"/dns4/lotus-bootstrap.ipfsforce.com/tcp/41778/p2p/12D3KooWGhufNmZHF3sv48aQeS13ng5XVJZ9E6qy2Ms4VzqeUsHk",

	"/dns4/bootstrap-0.starpool.in/tcp/12757/p2p/12D3KooWGHpBMeZbestVEWkfdnC9u7p6uFHXL1n7m1ZBqsEmiUzz",
	"/dns4/bootstrap-1.starpool.in/tcp/12757/p2p/12D3KooWQZrGH1PxSNZPum99M1zNvjNFM33d1AAu5DcvdHptuU7u",

	"/dns4/node.glif.io/tcp/1235/p2p/12D3KooWBF8cpp65hp2u9LK5mh19x67ftAam84z9LsfaquTDSBpt",

	"/dns4/bootstrap-0.ipfsmain.cn/tcp/34721/p2p/12D3KooWQnwEGNqcM2nAcPtRR9rAX8Hrg4k9kJLCHoTR5chJfz6d",
	"/dns4/bootstrap-1.ipfsmain.cn/tcp/34723/p2p/12D3KooWMKxMkD5DMpSWsW7dBddKxKT7L2GgbNuckz9otxvkvByP",

	"/dns4/bootstarp-0.1475.io/tcp/61256/p2p/12D3KooWQjaNmbz9b1XmheQB3RWsRjKSzuRLfjeiDZHyX7Y5RcBr",
}

var MainnetMiners = []string{
	"/ip4/222.209.189.157/tcp/8066",
	"/ip4/211.48.44.75/tcp/24001",
	"/ip4/114.86.142.162/tcp/24001",
	"/ip4/154.42.3.34/tcp/24001",
	"/ip4/86.123.188.55/tcp/24001",
	"/ip4/47.108.179.232/tcp/24001",
	"/ip4/183.232.117.122/tcp/4001",
	"/ip4/172.16.1.48/tcp/58418",
	"/ip4/172.16.1.48/tcp/58418",
	"/ip4/202.160.106.54/tcp/28019",
	"/ip4/202.160.106.54/tcp/28018",
	"/ip4/202.160.106.54/tcp/28019",
	"/ip4/202.160.106.54/tcp/28018",
	"/ip4/103.127.248.194/tcp/10001",
	"/ip4/103.127.248.196/tcp/10001",
	"/ip4/83.238.168.109/tcp/24002",
	"/ip4/185.32.160.239/tcp/42001",
	"/ip4/66.61.208.206/tcp/41372",
	"/ip4/66.61.208.206/tcp/41372",
	"/ip4/103.127.248.194/tcp/10001",
	"/ip4/103.127.248.196/tcp/10001",
	"/ip4/66.61.208.206/tcp/41372",
	"/ip4/66.61.208.206/tcp/41372",
	"/ip4/61.164.212.195/tcp/24001",
	"/ip4/10.8.1.7/tcp/6666",
	"/ip4/208.97.197.18/tcp/11337",
	"/ip4/86.123.188.55/tcp/24001",
	"/ip4/66.6.126.70/tcp/15001",
	"/ip4/182.18.86.72/tcp/46888",
	"/ip4/211.181.56.218/tcp/24001",
	"/ip4/122.114.37.226/tcp/12000",
	"/ip4/117.52.173.162/tcp/24003",
	"/ip4/218.78.187.146/tcp/21735",
	"/ip4/10.7.3.2/tcp/1025",
	"/ip4/10.7.3.2/tcp/1025",
	"/ip4/61.97.250.14/tcp/1235",
	"/ip4/182.18.86.72/tcp/46888",
	"/ip4/67.212.85.196/tcp/10906",
	"/ip4/183.93.252.68/tcp/1235",
	"/ip4/121.178.127.135/tcp/24001",
	"/ip4/10.8.1.7/tcp/6666",
	"/ip4/10.7.3.2/tcp/1025",
	"/ip4/10.8.1.9/tcp/6666",
	"/ip4/211.107.174.12/tcp/24001",
	"/ip4/54.219.44.132/tcp/8993",
	"/ip4/18.144.19.252/tcp/7743",
	"/ip4/54.193.243.116/tcp/8654",
	"/ip4/18.190.24.1/tcp/9473",
	"/ip4/35.90.15.192/tcp/8152",
	"/ip4/34.219.53.128/tcp/8083",
}

var log = logging.Logger("main")

func main() {
	logging.SetAllLoggers(logging.LevelInfo)
	logging.SetLogLevel("bitswap", "error")
	logging.SetLogLevel("connmgr", "error")
	logging.SetLogLevel("dht/RtRefreshManager", "error")
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app := &cli.App{
		Name: "surveyor",
		Commands: []*cli.Command{
			runCmd,
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "repo",
				Usage: "specify default node repo location",
				Value: filepath.Join(home, ".surveyor"),
			},
		},
	}

	app.Setup()
	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stdout, err.Error())
		os.Exit(1)
	}
}

var runCmd = &cli.Command{
	Name: "run",
	Before: func(cctx *cli.Context) error {
		return ensureRepoExists(cctx.String("repo"))
	},
	Action: func(cctx *cli.Context) error {
		repoDir := cctx.String("repo")

		inspect, err := survey.New(cctx.Context, &survey.Config{
			Libp2pKeyFile: filepath.Join(repoDir, "libp2p.key"),
		})
		if err != nil {
			return err
		}
		ctx := context.Background()
		results := make(chan *survey.Result)
		go inspect.Run(ctx, MinerPeerAddrInfo(), results, survey.NewAgentProcessor(inspect), survey.NewDHTProcessor(inspect), survey.NewPeerLocation(inspect))

		thing := make(map[string][]interface{})

		for res := range results {
			thing[res.ID.String()] = append(thing[res.ID.String()], res)
		}

		PrintJson(thing)

		return nil
	},
}

func ensureRepoExists(dir string) error {
	st, err := os.Stat(dir)
	if err == nil {
		if st.IsDir() {
			return nil
		}
		return fmt.Errorf("repo dir was not a directory")
	}

	if !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return nil
}

func PrintJson(obj interface{}) error {
	resJson, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling json: %w", err)
	}

	fmt.Println(string(resJson) + ",")
	return nil
}
