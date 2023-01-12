package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/richxcame/gotoleg/gotoleg"
	"github.com/richxcame/gotoleg/internal/transaction"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTransactionServer(s, &transaction.Server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// package main

// import (
// 	"context"

// 	"github.com/richxcame/gotoleg/gotoleg"
// 	"github.com/richxcame/gotoleg/internal/transaction"
// )

// func main() {
// 	// time, err := utility.GetEpoch()
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// fmt.Println(time)

// 	// balance, err := utility.CheckBalance()
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// fmt.Println(balance)

// 	// services, err := utility.GetServices()
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// fmt.Println(services)
// 	// fmt.Println(hmacsha1.Generate("l0NWdYaEnlRo199FuxnU+vpQGT/HAUe1yftdTf5yjs1urf9mtClvO3EYGqoXPkqIgDWy9mOjKi230YKREh2fBw==", "1672084585:demiryol"))
// 	s := transaction.Server{}
// 	s.Add(context.Background(), &gotoleg.AddTransactionRequest{LocalID: "3", Service: "", Phone: "62726535", Amount: "100"})

// 	// transaction.Add()
// 	// transaction.CheckStatus("2")

// }
