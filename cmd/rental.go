package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/corka149/rental"
	"github.com/corka149/rental/app"
	"github.com/corka149/rental/datastore"
	"github.com/corka149/rental/jobs"
	"github.com/corka149/rental/locales"
	"github.com/corka149/rental/middleware"
	"github.com/corka149/rental/schema"
	"github.com/invopop/ctxi18n"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/gzip"
	"github.com/joho/godotenv"
)

func NewServer(ctx context.Context, getenv func(string) string) (*Server, error) {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Load .env file
	err := godotenv.Load()

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	// Setup config
	config, err := rental.Setup(ctx, getenv)

	if err != nil {
		return nil, fmt.Errorf("failed to setup config: %w", err)
	}

	queries := datastore.New(config.DbPool)

	// Load locales
	if err := ctxi18n.Load(locales.Content); err != nil {
		log.Fatalf("error loading locales: %v", err)
	}

	// Run migration
	err = schema.RunMigration(ctx, config)

	if err != nil {
		return nil, fmt.Errorf("failed to run migration: %w", err)
	}

	// Session store
	db := stdlib.OpenDBFromPool(config.DbPool)
	store, err := postgres.NewStore(db, []byte(config.Secret))
	if err != nil {
		log.Fatalf("failed to create store: %v", err)
	}

	// Run cleanup
	jobs.CleanUp(ctx, queries)

	// Run server

	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middleware.NewI18n())
	router.Use(sessions.Sessions("rental", store))

	app.RegisterRoutes(router, ctx, queries, config)

	s := &Server{
		Router:  router,
		Config:  config,
		Queries: queries,
	}

	return s, nil
}

func (s *Server) Close() {
	s.Config.DbPool.Close()
}

type Server struct {
	Router  *gin.Engine
	Config  *rental.Config
	Queries *datastore.Queries
}

func (s *Server) Run() error {
	address := fmt.Sprintf(":%s", s.Config.Port)

	return s.Router.Run(address)
}
