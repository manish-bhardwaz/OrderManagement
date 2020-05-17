package store

import (
	"errors"
	"log"
)

var (
	// MaxProductCount is the maximun inventory of a product
	MaxProductCount int32 = 10
	// ErrSoldOut an error representing the product is out of stock
	ErrSoldOut error = errors.New("we're sorry, the product is out of stock, better luck next time")
	// FailedtoListenF format string of listen failuire error
	FailedtoListenF = "failed to listen, error : %v"
	// FailedToServeF format string of serve failure error
	FailedToServeF = "failed to serve, error : %v"
	// FailedToDialF format string of dial failure error
	FailedToDialF = "failed to dial, error : %v"
	// PortQuantiyServer port of Quantity Server
	PortQuantiyServer string = ":7777"
	// PortCartServer port of Cart Server
	PortCartServer string = ":8888"
	// PortCheckoutServer port of Checkout Server
	PortCheckoutServer string = ":9999"
	// PortPaymentServer port of Payment Gateway Server
	PortPaymentServer string = ":10000"
)

// ErrCheck utility error check method
func ErrCheck(err error, message string) {
	if err != nil {
		log.Fatalf(message, err.Error())
	}
}
