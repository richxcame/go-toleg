package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "gotoleg/rpc/gotoleg"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("couldn't read .env file: %v", err)
	}
}

func main() {
	addr := os.Getenv("GOTOLEG_PORT")
	fmt.Printf("localhost:%v", addr)
	// Set up a connection to the server.
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%v", addr), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTransactionClient(conn)

	// Contact the server and print out its response.
	ctx := metadata.AppendToOutgoingContext(context.Background(), "api_key", os.Getenv("TEST_GOTOLEG_API_KEY"))
	r, err := c.Add(ctx, &pb.TransactionRequest{LocalID: time.Now().String(), Service: "", Phone: os.Getenv("TEST_GOTOLEG_PHONE"), Amount: os.Getenv("TEST_GOTOLEG_AMOUNT")})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Println(r)

}
