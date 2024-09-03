package schema

import (
	"context"
	"embed"
	"github.com/corka149/rental"
	"log"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var migrations embed.FS

func RunMigration(ctx context.Context, config *rental.Config) error {

	db := stdlib.OpenDBFromPool(config.DbPool)

	defer db.Close()

	gooseProvider, err := goose.NewProvider(goose.DialectPostgres, db, migrations)

	if err != nil {
		return err
	}

	log.Println("Running migrations...")

	results, err := gooseProvider.Up(ctx)

	if err != nil {
		return err
	}

	for _, r := range results {
		log.Printf("Applied migration: %s\n", r)
	}

	log.Println("Migrations completed")

	return nil
}
