// Package rental contains the setup function for the rental package.
// The configuration can be done by the following environment variables:
// ADMIN_USER: The username for the admin user. Default is "admin".
// ADMIN_PASSWORD: The password for the admin user. Default is "password".
// DB_URL: The URL for the database. Default is "postgres://myadmin:mypassword@localhost:5432/rental_db".
// PORT: The port for the application. Default is "8081".
package rental

import (
	"cmp"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	AdminUsername     string
	AdminUserPassword string
	DbPool            *pgxpool.Pool
	Port              string
}

const (
	defaultDbURL         = "postgres://myadmin:mypassword@localhost:5432/rental_db"
	defaultAdminUser     = "admin"
	defaultAdminPassword = "password"
	defaultPort          = "8081"
)

func Setup(ctx context.Context, getenv func(string) string) (*Config, error) {
	dbUrl := cmp.Or(getenv("DB_URL"), defaultDbURL)

	dbpool, err := pgxpool.New(ctx, dbUrl)

	if err != nil {
		return nil, err
	}

	adminUser := cmp.Or(getenv("ADMIN_USER"), defaultAdminUser)
	adminPassword := cmp.Or(getenv("ADMIN_PASSWORD"), defaultAdminPassword)
	port := cmp.Or(getenv("PORT"), defaultPort)

	return &Config{
		AdminUsername:     adminUser,
		AdminUserPassword: adminPassword,
		DbPool:            dbpool,
		Port:              port,
	}, nil
}
