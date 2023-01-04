package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/richxcame/gotoleg/gotoleg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:6000", "the address to connect to")
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
	r, err := c.Add(ctx, &pb.AddTransactionRequest{LocalID: "5", Service: "", Phone: "62726535", Amount: "100"})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("success: %s", r.GetStatus())

}
