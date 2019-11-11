package fetchers

import (
	"sync"
)

type StatusFetcher struct {
	Coin string
	wg   sync.WaitGroup
}

func NewStatusFetcher(coin string) StatusFetcher {
	return StatusFetcher{Coin: coin}
}

// start the crawl
func (s *StatusFetcher) Run() {
	// for loop
	//for {
	//	// get all address by database.
	//	go func() {
	//		// get blocks / version / ping
	//		// write to channel
	//		// write to database.
	//		// update last_ping_at
	//	}
	//}
}
