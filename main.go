package main

import (
	"github.com/dubuqingfeng/bit-node-crawler/fetchers"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

func main() {
	// new node fetcher
	fetcher := fetchers.NewNodeFetcher()
	fetcher.Run()
	// exit sign
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig // Wait for ^C signal
		log.Warningf("Interrupt signal detected, shutting down.")
		//crawler.Shutdown()
	}()
}
