package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	// Run cron job to check declined and empty status transactions periodically
	transactionCronInterval := time.Hour
	go handlers.RunTransactionCronJob(transactionCronInterval)

	// Create database pool
	dbPool := db.CreateDB()
	defer dbPool.Close()

	// Start HTTP server
	routes := routes.SetupRoutes()
	// go routes.Run()
	srv := &http.Server{
		Addr:    ":" + os.Getenv("GOTOLEG_PORT"),
		Handler: routes,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Start GRPC server
	grpcPort := os.Getenv("GOTOLEG_PORT")
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", grpcPort, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTransactionServer(grpcServer, &transaction.Server{})

	log.Printf("gRPC server listening on port %s", grpcPort)
	if err := grpcServer.Serve(grpcListener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// Wait for an interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a context with a timeout for the graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
