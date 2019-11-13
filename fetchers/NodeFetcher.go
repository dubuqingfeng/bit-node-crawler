package fetchers

import (
	"github.com/dubuqingfeng/bit-node-crawler/handlers/coins"
	"github.com/dubuqingfeng/bit-node-crawler/helpers"
	"github.com/dubuqingfeng/bit-node-crawler/models"
	"github.com/dubuqingfeng/bit-node-crawler/seeds"
	"github.com/dubuqingfeng/bit-node-crawler/utils"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
	"sync"
)

type NodeFetcher struct {
	Coin           string
	workerChan     chan []string
	resultChan     chan models.Result
	wg             sync.WaitGroup
	checkFilter    map[string]bool
	completeFilter map[string]bool
}

func NewNodeFetcher(coin string) NodeFetcher {
	return NodeFetcher{workerChan: make(chan []string, utils.Config.Concurrency),
		resultChan: make(chan models.Result), Coin: coin}
}

// add addresses to channel
func (n *NodeFetcher) AddPeers(peers []string) {
	for _, peer := range peers {
		// bloom filter
		state, ok := n.checkFilter[peer]
		if !ok || !state {
			n.wg.Add(1)
			go n.HandleAddress(peer)
		}
	}
}

// run the crawler
func (n *NodeFetcher) Run() {
	// get seeds from default dns
	seedNodes := helpers.GetSeedsFromDNS(seeds.GetDefaultDNSSeeds(n.Coin))
	// start crawl from seedNodes
	go n.AddPeers(seedNodes)
	// wait group
	go func() { n.wg.Wait() }()
	// receive the worker(peers) / result channel
	for {
		// close the channel when completed
		if len(n.checkFilter) <= len(n.completeFilter) {
			close(n.workerChan)
			close(n.resultChan)
			break
		}
		select {
		case peers := <-n.workerChan:
			go n.AddPeers(peers)
		case result := <-n.resultChan:
			go n.WriteResult(result)
		}
	}
}

// handle address
func (n *NodeFetcher) HandleAddress(address string) {
	n.checkFilter[address] = true
	// initialize the handler
	handler, err := coins.HandlerFactory(n.Coin, address)
	if err != nil {
		log.Error(err)
		return
	}
	// connect the address
	if err := handler.Connect(); err != nil {
		log.Error(err)
		return
	}
	// FILO stack
	defer func() { n.completeFilter[address] = true }()
	defer n.wg.Done()
	defer handler.DisConnect()
	// send version and version ack
	result, err := handler.Handshake()
	if err != nil {
		log.Error(err)
		return
	}
	// send getaddr
	if err := handler.SendGetAddr(); err != nil {
		log.Error(err)
	}
	// get addresses
	addresses, err := handler.GetAddrResponse()
	if err != nil {
		log.Error(err)
	}
	if len(addresses) != 0 && !reflect.DeepEqual(addresses, []string{address}) {
		log.Info(addresses)
		log.Info(address)
		if !utils.Config.OnlyFetchDefaultSeeds {
			n.workerChan <- addresses
		}
	}
	result.CoinType = n.Coin
	result.Peers = strings.Join(addresses, ",")
	n.resultChan <- result
}

func (n *NodeFetcher) WriteResult(result models.Result) {
	// write to database.
	if err := models.InsertOrUpdatePeer(result); err != nil {
		log.Error(err)
	}
}
