package models

import (
	"errors"
	"fmt"
	"github.com/dubuqingfeng/bit-node-crawler/dbs"
	"github.com/dubuqingfeng/bit-node-crawler/utils"
	log "github.com/sirupsen/logrus"
	"time"
)

func InsertOrUpdatePeer(peer Result) error {
	conn := utils.Config.GlobalDatabase.Write.Name
	if !dbs.CheckDBConnExists(conn) {
		return errors.New("not found this conn")
	}
	tableName := "peers"
	now := time.Now().UTC()
	// update
	stmt := fmt.Sprintf("INSERT INTO `" + tableName + "` (address, height, peers, user_agent, coin_type," +
		" timestamp, notified_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)  ON DUPLICATE KEY UPDATE " +
		"height = ?, peers = ?, timestamp = ?, user_agent = ?, updated_at = ?, notified_at = ? ")
	_, err := dbs.DBMaps[conn].Exec(stmt, peer.Address, peer.Height,
		peer.Peers, peer.UserAgent, peer.CoinType, peer.Timestamp,
		now, now, now, peer.Height, peer.Peers, peer.Timestamp, peer.UserAgent, now, now)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
