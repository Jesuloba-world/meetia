package db

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"

	"github.com/meetia/backend/internal/config"
)

var DB *bun.DB

func Initialize(cfg *config.Config) *bun.DB {
	dsn := cfg.DBUrl
	if dsn == "" {
		log.Fatal("database url required, set it in environment variables")
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Failed to parse database DSN: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to create connection pool: %v", err)
	}

	sqldb := stdlib.OpenDBFromPool(pool)
	DB = bun.NewDB(sqldb, pgdialect.New())

	// Ping the database to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := DB.PingContext(ctx); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// DB.AddQueryHook(bundebug.NewQueryHook(
	// 	bundebug.WithVerbose(true),
	// ))

	slog.Info("Connected to database successfully")
	return DB
}

func Close() error {
	return DB.Close()
}
