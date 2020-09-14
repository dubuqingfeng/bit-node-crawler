package main

import (
	"github.com/dubuqingfeng/bit-node-crawler/dbs"
	"github.com/dubuqingfeng/bit-node-crawler/utils"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	utils.InitConfig("./configs/config.yaml")
	// utils.ConfigLocalFileSystemLogger("./logs/", "monitor.log", 7*time.Hour*24, time.Second*20)
	dbs.InitMySQLDB()
}

func main() {
	c := cron.New()
	// err := c.AddFunc("0 * * * * *", func() {
	// 	fetch := fetchers.NewAPIFetcher()
	// 	fetch.Handle()
	// })
	// if err != nil {
	// 	log.Error(err)
	// }
	c.Start()
	if utils.Config.Debug {
		go func() {
			if err := http.ListenAndServe("0.0.0.0:6060", nil); err != nil {
				log.Error(err)
			}
		}()
	}
	// does not stop any jobs already running
	defer c.Stop()
	// blocking forever
	select {}
}
