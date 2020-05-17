package main

import (
	"net"

	"golang.org/x/net/context"

	"OrderManagement/store"

	"google.golang.org/grpc"
)

// CartServer : represents a Quantity Server
type CartServer struct {
}

// NewCartServer : instantiates a new quantity server
func NewCartServer() *CartServer {
	return &CartServer{}
}

// AddToCart : adds the given product to user's cart
func (c *CartServer) AddToCart(ctx context.Context, p *store.UserProduct) (*store.UserProduct, error) {
	p.Added = true
	p.Paid = false
	return p, nil
}

func main() {

	l, err := net.Listen("tcp", store.PortCartServer)
	store.ErrCheck(err, store.FailedtoListenF)
	defer l.Close()

	cS := NewCartServer()
	gS := grpc.NewServer()
	defer gS.Stop()
	store.RegisterCartServer(gS, cS)

	err = gS.Serve(l)
	store.ErrCheck(err, store.FailedToServeF)
}
