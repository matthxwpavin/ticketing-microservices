package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/matthxwpavin/ticketing/auth/internal/database/impl/mongo"
	"github.com/matthxwpavin/ticketing/auth/internal/router"
	"github.com/matthxwpavin/ticketing/env"
	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/server"
)

func main() {
	// New logger then associate it with the context.
	logger, err := sugar.New()
	if err != nil {
		log.Fatalf("could not new logger: %v", err)
	}
	ctx := sugar.WithContext(context.Background(), logger)

	// Check required environment variables.
	if err := env.CheckRequiredEnvs([]env.EnvKey{
		env.JwtSecret,
		env.NatsURL,
		env.NatsConnName,
		env.DEV,
		env.MongoURI,
	}); err != nil {
		logger.Fatalw("failed to check required environments", "error", err)
	}

	if err := run(ctx); err != nil {
		os.Exit(1)
	}
}

// run initializes app's dependencies which needed by HTTP server
// then listen for incoming requests.
//
// Also it's reserved for dependencies which need to graceful shutdown
// by deferred shutdown/close functions, when there is an error or a
// terminating signal imcoming to stop the server listening, the function
// returns then those deferred functions execute.
func run(ctx context.Context) error {
	logger := sugar.FromContext(ctx)

	// Receive a context that the Done() channel receive when these signal arrived.
	ctx, done := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer done()

	// Connect to the database.
	db, err := mongo.Connect(ctx)
	if err != nil {
		logger.Errorw("could not connect to database", "error", err)
		return err
	}
	// Disconnect the database when main function returns.
	defer db.Disconnect(ctx)

	h := router.New(ctx, db)
	return server.ListenAndServe(ctx, ":3000", h)
}
