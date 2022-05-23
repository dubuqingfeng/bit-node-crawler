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
	checkFilter    sync.Map
	completeFilter sync.Map
}

func NewNodeFetcher(coin string) NodeFetcher {
	return NodeFetcher{
		workerChan: make(chan []string, utils.Config.Concurrency),
		resultChan: make(chan models.Result),
		Coin:       coin}
}

// add addresses to channel
func (n *NodeFetcher) AddPeers(peers []string) {
	for _, peer := range peers {
		// bloom filter
		// state, ok := n.checkFilter[peer]
		state, ok := n.checkFilter.Load(peer)
		if !ok || state != nil {
			n.wg.Add(1)
			go n.HandleAddress(peer)
		}
	}
}

// run the crawler
func (n *NodeFetcher) Run() {
	// get seeds from default dns
	seedNodes := helpers.GetSeedsFromDNS(seeds.GetDefaultDNSSeeds(n.Coin, ""))
	log.Info(seedNodes)
	// start crawl from seedNodes
	go n.AddPeers(seedNodes)
	// wait group
	go func() { n.wg.Wait() }()
	// receive the worker(peers) / result channel
	for {
		// close the channel when completed
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
	n.checkFilter.Store(address, true)
	// n.checkFilter[address] = true
	// initialize the handler
	handler, err := coins.HandlerFactory(n.Coin, address)
	if err != nil {
		log.WithField("address", address).WithField("step", "init").Error(err)
		return
	}
	// connect the address
	if err := handler.Connect(); err != nil {
		log.WithField("address", address).WithField("step", "connect").Error(err)
		return
	}
	// FILO stack
	defer func() { n.completeFilter.Store(address, true) }()
	defer n.wg.Done()
	defer handler.DisConnect()
	// send version and version ack
	result, err := handler.Handshake()
	if err != nil {
		log.WithField("address", address).WithField("step", "handshake").Error(err)
		return
	}
	// send getaddr
	if err := handler.SendGetAddr(); err != nil {
		log.WithField("address", address).WithField("step", "getaddr").Error(err)
	}
	// get addresses
	addresses, err := handler.GetAddrResponse()
	if err != nil {
		log.WithField("address", address).WithField("step", "getaddrresponse").Error(err)
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
	log.Info(result)
	if err := models.InsertOrUpdatePeer(result, n.Coin); err != nil {
		log.Error(err)
	}
}
