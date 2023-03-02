package db

import (
	"context"
	"gotoleg/pkg/logger"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func CreateDB() *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatalf("Failed to create new pool %v", err)
	}
	if err := dbpool.Ping(context.Background()); err != nil {
		logger.Fatalf("Couldn't connect to database %v", err)
	}

	DB = dbpool

	return dbpool
}
