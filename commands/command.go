package main

import (
	"github.com/dubuqingfeng/bit-node-crawler/dbs"
	"github.com/dubuqingfeng/bit-node-crawler/fetchers"
	"github.com/dubuqingfeng/bit-node-crawler/utils"
)

// init
func init() {
	utils.InitConfig("./configs/config.yaml")
	dbs.InitMySQLDB()
}

// main
func main() {
	fetch := fetchers.NewStatusFetcher("btc")
	fetch.Run()
}
