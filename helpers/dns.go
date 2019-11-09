package helpers

import (
	log "github.com/sirupsen/logrus"
	"net"
	"sync"
)

// by https://github.com/shazow/btc-crawl/blob/master/seed.go
func GetSeedsFromDNS(dnsSeeds []string) []string {
	wait := sync.WaitGroup{}
	results := make(chan []net.IP)

	for _, seed := range dnsSeeds {
		wait.Add(1)
		go func(address string) {
			defer wait.Done()
			ips, err := net.LookupIP(address)
			if err != nil {
				log.Warningf("Failed to resolve %s: %v", address, err)
				return
			}
			log.Debugf("Resolved %d seeds from %s.", len(ips), address)
			results <- ips
		}(seed)
	}

	go func() {
		wait.Wait()
		close(results)
	}()

	var seeds []string
	for ips := range results {
		for _, ip := range ips {
			seeds = append(seeds, net.JoinHostPort(ip.String(), "8333"))
		}
	}

	log.Infof("Resolved %d seed nodes from %d DNS seeds.", len(seeds), len(dnsSeeds))
	// Note that this will likely include duplicate seeds. The crawler deduplicates them.
	return seeds
}
