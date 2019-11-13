package utils

import (
	"github.com/jinzhu/configor"
	"time"
)

// config
var Config = struct {
	Name           string `default:"app_name"`
	Debug          bool
	MonitorName    string `default:"bit-node-crawler"`
	CoinType       string `default:"btc"`
	GlobalDatabase MySQLDB
	// crawl config
	Concurrency           int           `default:"0"`
	ConnectTimeout        time.Duration `default:"10s"`
	UserAgent             string        `default:"/btc-crawl:0.1.2/"`
	OnlyFetchDefaultSeeds bool
}{}

// dsn
type MySQLDSN struct {
	Name string
	DSN  string
}

// mysql db
type MySQLDB struct {
	Read     MySQLDSN
	Write    MySQLDSN
	Timezone string
}

// init config, example: config.example
func InitConfig(files string) {
	err := configor.Load(&Config, files)
	if err != nil {
		panic(err)
	}
}

// get all database configs
func GetAllDatabaseConfigs() map[string]string {
	configs := make(map[string]string)
	AddDatabaseConfig(Config.GlobalDatabase, configs)
	return configs
}

// add database config
func AddDatabaseConfig(value MySQLDB, configs map[string]string) {
	if value.Read.DSN != "" && value.Read.Name != "" {
		configs[value.Read.Name] = value.Read.DSN
	}
	if value.Write.DSN != "" && value.Write.Name != "" {
		configs[value.Write.Name] = value.Write.DSN
	}
}
