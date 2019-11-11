package main

import (
	"github.com/dubuqingfeng/bit-node-crawler/dbs"
	"github.com/dubuqingfeng/bit-node-crawler/fetchers"
	"github.com/dubuqingfeng/bit-node-crawler/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
)

func init() {
	utils.InitConfig("./configs/config.yaml")
	dbs.InitMySQLDB()
}

func main() {
	// new node fetcher
	if utils.Config.Debug {
		go func() {
			if err := http.ListenAndServe("0.0.0.0:6061", nil); err != nil {
				log.Error(err)
			}
		}()
	}
	fetcher := fetchers.NewNodeFetcher("btc")
	fetcher.Run()
	// exit sign
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig // Wait for ^C signal
		log.Warningf("Interrupt signal detected, shutting down.")
	}()
}
