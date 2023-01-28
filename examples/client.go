package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "gotoleg/rpc/gotoleg"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTransactionClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Add(ctx, &pb.TransactionRequest{LocalID: time.Now().String(), Service: "", Phone: "+99362254853", Amount: "100"})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Println(r)

}
