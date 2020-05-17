package main

import (
	"OrderManagement/store"
	"context"
	"fmt"
	"log"
	"strconv"

	"google.golang.org/grpc"
)

var (
	cO store.CheckoutClient
	cC store.CartClient
)

func purchase(u *store.UserProduct) {
	ctx := context.Background()
	p, err := cC.AddToCart(ctx, u)
	if err != nil {
		//log.Println(fmt.Sprintf("%s, %s", u.GetID(), err.Error()))
		return
	}
	p, err = cO.MakePayment(ctx, p)
	if err != nil {
		//log.Println(fmt.Sprintf("%s, %s", u.GetID(), err.Error()))
		return
	}
	log.Println(fmt.Sprintf("%s sold to %s", u.GetProductID(), u.GetID()))
	return
}

func main() {

	checkOutServerConn, err := grpc.Dial(store.PortCheckoutServer, grpc.WithInsecure())
	store.ErrCheck(err, store.FailedToDialF)
	defer checkOutServerConn.Close()
	cO = store.NewCheckoutClient(checkOutServerConn)

	addToCartServerCon, err := grpc.Dial(store.PortCartServer, grpc.WithInsecure())
	store.ErrCheck(err, store.FailedToDialF)
	defer addToCartServerCon.Close()
	cC = store.NewCartClient(addToCartServerCon)

	var (
		maxNoOfGoRtns int32 = 5000
		noOfChkOuts   int32 = 10000
		semaphore     chan struct{}
		users         []*store.UserProduct
	)

	semaphore = make(chan struct{}, maxNoOfGoRtns)
	users = make([]*store.UserProduct, noOfChkOuts)

	for i := 0; i < int(noOfChkOuts); i++ {
		semaphore <- struct{}{}
		users[i] = &store.UserProduct{ID: "User " + strconv.Itoa(i+1), Added: false, Paid: false, ProductID: "Product X"}
		go func(u *store.UserProduct) {
			purchase(u)
			<-semaphore
		}(users[i])
	}
}
