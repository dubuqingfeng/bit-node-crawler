package seeds

import "strings"

var DefaultBTCDNSSeeds = []string{
	"seed.bitcoin.sipa.be",
	"dnsseed.bluematt.me",
	"dnsseed.bitcoin.dashjr.org",
	"seed.bitcoinstats.com",
	"seed.bitnodes.io",
	"seed.bitcoin.jonasschnelli.ch",
	"seed.btc.petertodd.org",
	"seed.bitcoin.sprovoost.nl",
	"dnsseed.emzy.de",
}

var DefaultBCHDNSSeeds = []string{
	"seed.bitcoinabc.org",
	// "seed-abc.bitcoinforks.org",
	"btccash-seeder.bitcoinunlimited.info",
	// "seed.bitprim.org",
	"seed.deadalnix.me",
	"seed.bchd.cash",

	"seed.flowee.cash",
	"seed-bch.bitcoinforks.org",
}

var DefaultBSVDNSSeeds = []string{
	"seed.bitcoinsv.io",
	"seed.cascharia.com",
	"seed.satoshisvision.network",
}

func GetDefaultDNSSeeds(coin string) []string {
	// TODO add testnet
	coin = strings.ToLower(coin)
	if coin == "btc" {
		return DefaultBTCDNSSeeds
	}
	if coin == "bch" {
		return DefaultBCHDNSSeeds
	}
	if coin == "bsv" {
		return DefaultBSVDNSSeeds
	}
	return DefaultBTCDNSSeeds
}