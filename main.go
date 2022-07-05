package main

import (
	"context"
	"time"

	logging "github.com/ipfs/go-log/v2"

	"github.com/frrist/surveyor/collectors/peeragent"
	filecoin2 "github.com/frrist/surveyor/networks/filecoin"
	"github.com/frrist/surveyor/storage"
)

var log = logging.Logger("main")

func main() {
	logging.SetAllLoggers(logging.LevelInfo)
	db, err := storage.NewDatabase("host=192.168.1.125 user=postgres password=password dbname=postgres port=5432 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	res, err := peeragent.QueryAllMinerInfo(db)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.TODO()
	fp, err := filecoin2.New(ctx, filecoin2.Config{
		Bootstrap:         filecoin2.CalibnetPeers,
		Datastore:         nil,
		DHTProtocolPrefix: filecoin2.CalibnetDHTPrefix,
	})
	if err := fp.Bootstrap(ctx); err != nil {
		log.Fatal(err)
	}
	// let it bootstrap
	time.Sleep(time.Second * 5)
	err = peeragent.FindMinerAgentProtocols(db, fp, res)
	if err != nil {
		log.Fatal(err)
	}
}
