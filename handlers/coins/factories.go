package coins

import (
	"github.com/dubuqingfeng/bit-node-crawler/handlers"
)

func HandlerFactory(coin string, address string) (handlers.Handler, error) {
	//switch strings.ToLower(coin) {
	//return nil, errors.New("not support this chain")
	return NewBTCHandler(address), nil
}
