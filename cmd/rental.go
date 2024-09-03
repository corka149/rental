package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/corka149/rental"
	"github.com/corka149/rental/app"
	"github.com/corka149/rental/datastore"
	"github.com/corka149/rental/schema"
	"log"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/gzip"
	"github.com/joho/godotenv"
)

func main() {
	err := Run(context.Background(), os.Getenv)

	log.Fatalln(err)
}

func Run(ctx context.Context, getenv func(string) string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Load .env file
	err := godotenv.Load()

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	config, err := rental.Setup(ctx, getenv)

	if err != nil {
		return fmt.Errorf("failed to setup config: %w", err)
	}

	defer config.DbPool.Close()

	err = schema.RunMigration(ctx, config)

	if err != nil {
		return fmt.Errorf("failed to run migration: %w", err)
	}

	queries := datastore.New(config.DbPool)

	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	app.RegisterRoutes(router, ctx, queries, config)

	address := fmt.Sprintf(":%s", config.Port)

	err = router.Run(address)

	if err != nil {
		return fmt.Errorf("failed to run server: %w", err)
	}

	return nil
}
