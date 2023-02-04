package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"gotoleg/internal/transaction"

	pb "gotoleg/rpc/gotoleg"

	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
)

func main() {
	gotolegPort := os.Getenv("GOTOLEG_PORT")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", gotolegPort))
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
