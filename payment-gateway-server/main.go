package main

import (
	"net"

	"golang.org/x/net/context"

	"OrderManagement/store"

	"google.golang.org/grpc"
)

// PaymentGatewayServer : represents a Quantity Server
type PaymentGatewayServer struct {
}

// NewPaymentGatewayServer : instantiates a new quantity server
func NewPaymentGatewayServer() *PaymentGatewayServer {
	return &PaymentGatewayServer{}
}

// ProcessPayment : processes the pauyment
func (s *PaymentGatewayServer) ProcessPayment(context.Context, *store.PaymentRequest) (*store.PaymentResponse, error) {
	return &store.PaymentResponse{}, nil
}

func main() {

	l, err := net.Listen("tcp", store.PortPaymentServer)
	store.ErrCheck(err, store.FailedtoListenF)
	defer l.Close()

	pS := NewPaymentGatewayServer()
	gS := grpc.NewServer()
	defer gS.Stop()
	store.RegisterPaymentGatewayServer(gS, pS)

	err = gS.Serve(l)
	store.ErrCheck(err, store.FailedtoListenF)
}
