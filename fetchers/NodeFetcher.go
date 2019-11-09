package fetchers

import (
	"fmt"
	"github.com/dubuqingfeng/bit-node-crawler/handlers/coins"
	"github.com/dubuqingfeng/bit-node-crawler/helpers"
	"github.com/dubuqingfeng/bit-node-crawler/seeds"
	log "github.com/sirupsen/logrus"
	"reflect"
	"sync"
)

type NodeFetcher struct {
	Coin        string
	workerChan  chan []string
	resultChan  chan string
	wg          sync.WaitGroup
	checkFilter map[string]bool
}

func NewNodeFetcher() NodeFetcher {
	return NodeFetcher{workerChan: make(chan []string), resultChan: make(chan string)}
}

// add addresses to channel
func (n *NodeFetcher) AddPeers(peers []string) {
	for _, peer := range peers {
		// bloom filter
		// TODO need remove duplicate
		// map
		n.wg.Add(1)
		go n.HandleAddress(peer)
	}
}

// run the crawler
func (n *NodeFetcher) Run() {
	// get seeds from default dns
	seedNodes := helpers.GetSeedsFromDNS(seeds.DefaultBTCDnsSeeds)
	// start crawl from seedNodes
	go n.AddPeers(seedNodes)
	// wait group
	go func() { n.wg.Wait() }()
	// receive the worker(peers) / result channel
	for {
		select {
		case peers := <-n.workerChan:
			go n.AddPeers(peers)
		case result := <-n.resultChan:
			n.WriteResult(result)
			// when after the completion of the closed channel
		}
	}
}

// handle address
func (n *NodeFetcher) HandleAddress(address string) {
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
	defer handler.DisConnect()
	defer n.wg.Done()
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
		// add config
		n.workerChan <- addresses
	}
	n.resultChan <- address + "success" + result.Status
}

func (n *NodeFetcher) WriteResult(result string) {
	fmt.Println(result)
}
