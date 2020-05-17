package main

import (
	"net"

	"golang.org/x/net/context"

	"OrderManagement/store"

	"google.golang.org/grpc"
)

// QuantityServer : represents a Quantity Server
type QuantityServer struct {
}

// NewQuantityServer : instantiates a new quantity server
func NewQuantityServer() *QuantityServer {
	return &QuantityServer{}
}

// Get : returns the available quantity of given product
func (s *QuantityServer) Get(ctx context.Context, p *store.Product) (*store.ProductQuantity, error) {
	return &store.ProductQuantity{
		ProductID: p.GetID(),
		Quantity:  store.MaxProductCount,
	}, nil
}

// Set : decreases the product count by 1
func (s *QuantityServer) Set(ctx context.Context, p *store.ProductQuantity) (*store.ProductQuantity, error) {
	if store.MaxProductCount > 0 {
		store.MaxProductCount = store.MaxProductCount - 1
		return &store.ProductQuantity{
			ProductID: p.GetProductID(),
			Quantity:  store.MaxProductCount,
		}, nil
	}
	return nil, store.ErrSoldOut
}

func main() {

	l, err := net.Listen("tcp", store.PortQuantiyServer)
	store.ErrCheck(err, store.FailedtoListenF)
	defer l.Close()

	qS := NewQuantityServer()
	gS := grpc.NewServer()
	defer gS.Stop()
	store.RegisterQuantityServer(gS, qS)

	err = gS.Serve(l)
	store.ErrCheck(err, store.FailedtoListenF)
}
