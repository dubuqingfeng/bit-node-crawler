package coins

import (
	"errors"
	"github.com/dubuqingfeng/bit-node-crawler/handlers"
	"strings"
)

func HandlerFactory(coin string, address string) (handlers.Handler, error) {
	switch strings.ToLower(coin) {
	case "btc", "bch", "bsv":
		return NewBTCHandler(address), nil
	default:
		return nil, errors.New("not support this chain")
	}
}
