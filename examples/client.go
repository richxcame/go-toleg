package main

import (
	"context"
	"log"
	"os"

	pb "gotoleg/rpc/gotoleg"

	"github.com/google/uuid"
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
	gotolegIP := os.Getenv("GOTOLEG_IP")
	gotolegPort := os.Getenv("GOTOLEG_PORT")

	// Set up a connection to the server.
	conn, err := grpc.Dial(gotolegIP+":"+gotolegPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTransactionClient(conn)

	// Contact the server and print out its response.
	ctx := metadata.AppendToOutgoingContext(context.Background(), "api_key", os.Getenv("TEST_GOTOLEG_API_KEY"))
	r, err := c.Add(ctx, &pb.TransactionRequest{LocalID: uuid.New().String(), Service: "", Phone: os.Getenv("TEST_GOTOLEG_PHONE"), Amount: os.Getenv("TEST_GOTOLEG_AMOUNT")})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Println(r)

}
