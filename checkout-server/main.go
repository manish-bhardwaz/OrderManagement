package main

import (
	"log"
	"net"
	"sync"

	"golang.org/x/net/context"

	"OrderManagement/store"

	"google.golang.org/grpc"
)

// CheckoutServer : represents a Quantity Server
type CheckoutServer struct {
	store.QuantityClient
	store.PaymentGatewayClient
	sync.RWMutex
}

// NewCheckoutServer : instantiates a new quantity server
func NewCheckoutServer(qC store.QuantityClient, pC store.PaymentGatewayClient) *CheckoutServer {
	return &CheckoutServer{
		QuantityClient:       qC,
		PaymentGatewayClient: pC,
	}
}

func processPaymentAndUpdateCount(ctx context.Context, c *CheckoutServer, pQ *store.ProductQuantity) error {
	_, err := c.PaymentGatewayClient.ProcessPayment(ctx, &store.PaymentRequest{})
	if err != nil {
		return err
	}
	_, err = c.QuantityClient.Set(ctx, pQ)
	if err != nil {
		return err
	}
	return nil
}

// MakePayment : mimicks Payment Gateway to complete the purchase
func (c *CheckoutServer) MakePayment(ctx context.Context, p *store.UserProduct) (*store.UserProduct, error) {
	c.RWMutex.RLock()
	pQ, err := c.QuantityClient.Get(ctx, &store.Product{
		ID: p.GetProductID(),
	})
	if err != nil {
		c.RWMutex.RUnlock()
		return nil, err
	}
	c.RWMutex.RUnlock()

	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()
	err = processPaymentAndUpdateCount(ctx, c, pQ)
	if err != nil {
		return nil, err
	}
	log.Printf("payment complete, %s sold to %s\n", p.GetProductID(), p.GetID())
	p.Added = false
	p.Paid = true
	return p, nil
}

func main() {

	lC, err := net.Listen("tcp", store.PortCheckoutServer)
	store.ErrCheck(err, store.FailedtoListenF)
	defer lC.Close()

	qCon, err := grpc.Dial(store.PortQuantiyServer, grpc.WithInsecure())
	defer qCon.Close()
	store.ErrCheck(err, store.FailedToDialF)
	qC := store.NewQuantityClient(qCon)

	pCon, err := grpc.Dial(store.PortPaymentServer, grpc.WithInsecure())
	defer pCon.Close()
	store.ErrCheck(err, store.FailedToDialF)
	pC := store.NewPaymentGatewayClient(pCon)

	cS := NewCheckoutServer(qC, pC)
	gS := grpc.NewServer()
	defer gS.Stop()
	store.RegisterCheckoutServer(gS, cS)

	err = gS.Serve(lC)
	store.ErrCheck(err, store.FailedToServeF)
}
