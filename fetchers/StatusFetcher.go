package fetchers

import (
	"fmt"
	"github.com/dubuqingfeng/bit-node-crawler/handlers/coins"
	"github.com/dubuqingfeng/bit-node-crawler/models"
	"github.com/dubuqingfeng/bit-node-crawler/utils"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

type StatusFetcher struct {
	Coin  string
	wg    sync.WaitGroup
	count int64
}

func NewStatusFetcher(coin string) StatusFetcher {
	return StatusFetcher{Coin: coin}
}

// start the crawl
func (s *StatusFetcher) Run() {
	// for loop
	s.SyncNodes()
	// for {
	// 	s.SyncNodes()
	// 	time.Sleep(20 * time.Second)
	// }
}

func (s *StatusFetcher) SyncNodes() {
	start := time.Now()
	// get all address by database.
	nodes, err := models.GetAllPeers()
	if err != nil {
		log.Error(err)
		return
	}
	if len(nodes) == 0 {
		return
	}
	for _, node := range nodes {
		s.wg.Add(1)
		go s.SyncNode(node)
	}
	s.wg.Wait()
	log.WithField(fmt.Sprintf("获取 %d 个数据用时", len(nodes)), time.Since(start)).
		WithField("开始时间", start).Info("获取完毕")
}

func (s *StatusFetcher) SyncNode(node models.Result) {
	// get blocks / version / ping
	handler, err := coins.HandlerFactory(s.Coin, node.Address)
	if err != nil {
		log.Error(err)
		return
	}
	// connect the address
	if err := handler.Connect(); err != nil {
		log.Error(err)
		return
	}
	defer handler.DisConnect()
	defer s.wg.Done()
	defer func() {
		s.count = s.count + 1
		fmt.Println(s.count)
	}()
	// send version and version ack
	result, err := handler.Handshake()
	if err != nil {
		log.Error(err)
		return
	}
	// send getaddr
	if err := handler.SendGetAddr(); err != nil {
		log.Error(err)
		return
	}
	addresses, err := handler.GetAddrResponse()
	if err != nil {
		log.Error(err)
		return
	}
	result.CoinType = s.Coin
	result.Peers = strings.Join(addresses, ",")
	if result.Height != node.Height {
		result.HeightChangedAt = time.Now().UTC().Format(utils.UTCDatetime)
	} else {
		result.HeightChangedAt = node.HeightChangedAt
	}
	if err := models.InsertOrUpdateNode(result); err != nil {
		log.Error(err)
		return
	}
}
