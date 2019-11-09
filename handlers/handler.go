package handlers

import "github.com/dubuqingfeng/bit-node-crawler/models"

type Handler interface {
	Connect() error
	DisConnect()
	Handshake() (models.Result, error)
	SendGetAddr() error
	GetAddrResponse() ([]string, error)
}
