package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"gotoleg/internal/db"
	"gotoleg/internal/transaction"

	pb "gotoleg/rpc/gotoleg"
	"gotoleg/web/handlers"
	"gotoleg/web/routes"

	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
)

func main() {
	// run cron job, for declined and empty status transactions
	interval := time.Hour
	go handlers.CronJob(interval)

	// Close db pool
	pool := db.CreateDB()
	defer pool.Close()

	// HTTP server
	r := routes.SetupRoutes()
	go r.Run()

	// GRPC server
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
