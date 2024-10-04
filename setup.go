// Package rental contains the setup function for the rental package.
// The configuration can be done by the following environment variables:
// ADMIN_USER: The username for the admin user. Default is "admin".
// ADMIN_PASSWORD: The password for the admin user. Default is "password".
// DB_URL: The URL for the database. Default is "postgres://myadmin:mypassword@localhost:5432/rental_db".
// PORT: The port for the application. Default is "8081".
// SECRET: The secret for the application. Default is "secret".
// MODE: The mode for the application. Default is "DEV" which logs and returns extra details.
package rental

import (
	"cmp"
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	AdminUsername     string
	AdminUserPassword string
	DbPool            *pgxpool.Pool
	Port              string
	Secret            string
}

const (
	defaultDbURL         = "postgres://myadmin:mypassword@localhost:5432/rental_db"
	defaultAdminUser     = "admin"
	defaultAdminPassword = "password"
	defaultPort          = "8081"
	defaultSecret        = "secret"
)

func Setup(ctx context.Context, getenv func(string) string) (*Config, error) {
	mode := cmp.Or(getenv("MODE"), "DEV")

	if mode == "DEV" {
		log.Printf("!!! Running in development mode !!!\n")
	}

	dbUrl := cmp.Or(getenv("DB_URL"), defaultDbURL)

	dbConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Printf("Unable to parse database URL: %v\n", err)
		return nil, err
	}

	if mode == "DEV" {
		logger := log.New(os.Stdout, "rental", log.LstdFlags)
		dbConfig.ConnConfig.Tracer = &rentalQueryTracer{log: logger}
	}

	dbpool, err := pgxpool.NewWithConfig(ctx, dbConfig)

	if err != nil {
		return nil, err
	}

	adminUser := cmp.Or(getenv("ADMIN_USER"), defaultAdminUser)
	adminPassword := cmp.Or(getenv("ADMIN_PASSWORD"), defaultAdminPassword)
	port := cmp.Or(getenv("PORT"), defaultPort)
	secret := cmp.Or(getenv("SECRET"), defaultSecret)

	return &Config{
		AdminUsername:     adminUser,
		AdminUserPassword: adminPassword,
		DbPool:            dbpool,
		Port:              port,
		Secret:            secret,
	}, nil
}

type rentalQueryTracer struct {
	log *log.Logger
}

func (tracer *rentalQueryTracer) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData) context.Context {
	tracer.log.Printf("Executing command '%s' (with '%s')", data.SQL, data.Args)

	return ctx
}

func (tracer *rentalQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
}
