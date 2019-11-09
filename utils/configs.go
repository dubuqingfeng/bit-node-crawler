package utils

import "github.com/jinzhu/configor"

// config
var Config = struct {
	Name           string `default:"app_name"`
	IsDebug        bool
	MonitorName    string `default:"api-monitor"`
	CoinType       string
	GlobalDatabase MySQLDB
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
