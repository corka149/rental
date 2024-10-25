// Package rental contains the setup function for the rental package.
// The configuration can be done by the following environment variables:
// ADMIN_USER: The username for the admin user. Default is "admin".
// ADMIN_PASSWORD: The password for the admin user. Default is "password".
// DB_HOST: The host for the database. Default is "localhost".
// DB_PORT: The port for the database. Default is "5432".
// DB_USER: The user for the database. Default is "myadmin".
// DB_PASSWORD: The password for the database. Default is "mypassword".
// DB_NAME: The name for the database. Default is "rental_db".
// PORT: The port for the application. Default is "8081".
// SECRET: The secret for the application. Default is "secret".
package rental

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
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
	defaultDbHost        = "localhost"
	defaultDbPort        = "5432"
	defaultDbUser        = "myadmin"
	defaultDbName        = "rental_db"
	defaultDbPassword    = "mypassword"
	defaultAdminUser     = "admin"
	defaultAdminPassword = "password"
	defaultPort          = "8081"
	defaultSecret        = "secret"
)

func Setup(ctx context.Context, getenv func(string) string) (*Config, error) {
	isDebugMode := gin.Mode() == gin.DebugMode

	dbHost := cmp.Or(getenv("DB_HOST"), defaultDbHost)
	dbPort := cmp.Or(getenv("DB_PORT"), defaultDbPort)
	dbUser := cmp.Or(getenv("DB_USER"), defaultDbUser)
	dbPassword := cmp.Or(getenv("DB_PASSWORD"), defaultDbPassword)
	dbName := cmp.Or(getenv("DB_NAME"), defaultDbName)

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	dbConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Printf("Unable to parse database URL: %v\n", err)
		return nil, err
	}

	if isDebugMode {
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

	args := ""

	for i, arg := range data.Args {
		if i > 0 {
			args += ", "
		}
		args += fmt.Sprintf("$%d=%s", i+1, arg)
	}

	tracer.log.Printf("Executing command '%s' (with '%s')\n", data.SQL, args)

	return ctx
}

func (tracer *rentalQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
}
