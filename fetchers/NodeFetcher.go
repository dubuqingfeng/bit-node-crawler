package fetchers

import (
	"fmt"
	"github.com/dubuqingfeng/bit-node-crawler/helpers"
	"github.com/dubuqingfeng/bit-node-crawler/seeds"
	"sync"
)

type NodeFetcher struct {
	Coin        string
	workerChan  chan []string
	wg          sync.WaitGroup
	checkFilter map[string]bool
}

func NewNodeFetcher() NodeFetcher {
	workerChan := make(chan []string)
	return NodeFetcher{workerChan: workerChan}
}

func (n NodeFetcher) AddPeers(peers []string) {
	for _, peer := range peers {
		// bloom filter
		// TODO need remove duplicate
		n.wg.Add(1)
		go n.HandleAddress(peer)
	}
}

func (n NodeFetcher) HandleAddress(address string) {
	if len(address) < 40 {
		var str []string
		str = append(str, address+"test")
		n.workerChan <- str
	}
	n.wg.Done()
}

func (n NodeFetcher) Run() {
	// get seeds from default dns
	seedNodes := helpers.GetSeedsFromDNS(seeds.DefaultBCHDnsSeeds)
	// start crawl from seedNodes
	go n.AddPeers(seedNodes)
	// wait group
	go func() { n.wg.Wait() }()
	// Receive the worker(peers) channel
	for peers := range n.workerChan {
		fmt.Println(peers)
		go n.AddPeers(peers)
	}
}
