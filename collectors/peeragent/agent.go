package peeragent

import (
	"context"
	_ "embed"
	"encoding/json"
	"time"

	logging "github.com/ipfs/go-log/v2"
	"github.com/lib/pq"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"gorm.io/datatypes"
	"gorm.io/gorm/clause"

	"github.com/frrist/surveyor/networks/filecoin"
	"github.com/frrist/surveyor/storage"
)

var log = logging.Logger("collectors/peeragent")

type MinerInfoResult struct {
	Height         int64          `gorm:"column:height"`
	MinerID        string         `gorm:"column:miner_id"`
	PeerID         string         `gorm:"column:peer_id"`
	MultiAddresses datatypes.JSON `gorm:"column:multi_addresses"`
}

func QueryAllMinerInfo(db *storage.Database) ([]MinerInfoResult, error) {
	var results []MinerInfoResult
	tx := db.Db.Raw(`
select
    distinct on (miner_id)
    height,
    miner_id,
    peer_id,
    multi_addresses
from
    miner_infos
where
  multi_addresses != 'null' and peer_id != 'null'
order by
    miner_id,
    height desc
`).Scan(&results)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return results, nil
}

func QueryActiveMinerInfo(db *storage.Database) ([]MinerInfoResult, error) {
	var results []MinerInfoResult
	tx := db.Db.Raw(
		`
	select
	  distinct on (pac.miner_id) pac.height,
	  pac.miner_id,
	  m.peer_id,
	  m.multi_addresses
	from
	  power_actor_claims pac
	join
	  miner_infos m on m.miner_id = pac.miner_id
	where
	  pac.height between 0  and current_height()
	  and m.multi_addresses != 'null' and m.peer_id != 'null'
	order by
	  pac.miner_id,
	  pac.height desc
`).Scan(&results)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return results, nil
}

type MinerAgentProtocols struct {
	CreatedAt time.Time
	Height    int64          `gorm:"primaryKey"`
	MinerID   string         `gorm:"primaryKey"`
	Agent     string         `gorm:"index"`
	Protocols pq.StringArray `gorm:"type:text[]"`
}

func FindMinerAgentProtocols(db *storage.Database, fp *filecoin.Peer, todo []MinerInfoResult) error {
	if err := db.Db.AutoMigrate(&MinerAgentProtocols{}); err != nil {
		return err
	}
	tracker := make(map[peer.ID]MinerInfoResult)
	toConnect := make([]peer.AddrInfo, 0, len(todo))
	for _, t := range todo {
		maybe, err := t.MultiAddresses.MarshalJSON()
		if err != nil {
			return err
		}
		var out []string
		if err := json.Unmarshal(maybe, &out); err != nil {
			return err
		}
		var maddrs []multiaddr.Multiaddr
		for _, o := range out {
			maddr, err := multiaddr.NewMultiaddr(o)
			if err != nil {
				return err
			}
			maddrs = append(maddrs, maddr)
		}
		pid, err := peer.Decode(t.PeerID)
		if err != nil {
			return err
		}
		toConnect = append(toConnect, peer.AddrInfo{
			ID:    pid,
			Addrs: maddrs,
		})
		tracker[pid] = t
	}
	results := fp.GetAllPeerAgentProtocols(context.TODO(), toConnect, time.Minute, 100)
	out := make([]MinerAgentProtocols, 0, len(todo))
	for res := range results {
		miner := tracker[res.Peer.ID]
		out = append(out, MinerAgentProtocols{
			Height:    miner.Height,
			MinerID:   miner.MinerID,
			Agent:     res.Agent,
			Protocols: res.Protocols,
		})

	}
	return db.Db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).CreateInBatches(out, 10).Error
}
