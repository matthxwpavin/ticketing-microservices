package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/matthxwpavin/ticketing/env"
	"github.com/matthxwpavin/ticketing/expiration/internal/app"
	"github.com/matthxwpavin/ticketing/expiration/internal/streamer/nats"
	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/redis/go-redis/v9"
)

func main() {
	// New logger then associate it with the context.
	logger, err := sugar.New()
	if err != nil {
		log.Fatalf("could not new logger: %v", err)
	}
	ctx := sugar.WithContext(context.Background(), logger)

	// Check environment variables.
	if err := env.CheckRequiredEnvs([]env.EnvKey{
		env.NatsConnName,
		env.NatsURL,
		env.RedisHost,
	}); err != nil {
		logger.Fatalw("check env failed", "error", err)
	}

	if err := run(ctx); err != nil {
		os.Exit(1)
	}
}

// run initializes app's dependencies which needed by HTTP server
// then listen for incoming requests.
//
// Also it's reserved for dependencies which need to shutdown gracefully
// by deferred shutdown/close functions, when there is an error or a
// terminating signal imcoming to stop the server listening, the function
// returns then those deferred functions execute.
func run(ctx context.Context) error {
	logger := sugar.FromContext(ctx)

	// Receive a context that the Done() channel receive when these signal arrived.
	ctx, done := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer done()

	// Connect NATS streaming.
	nc, err := nats.ConnectFromEnv(ctx)
	if err != nil {
		logger.Errorw("could not connect to NATS", "error", err)
		return err
	}
	defer nc.Disconenct(ctx)

	redisCli := redis.NewClient(&redis.Options{
		Addr: env.RedisHost.Value() + ":6379",
	})

	_, err = app.NewService(ctx, nc, redisCli)
	if err != nil {
		logger.Errorw("could not start service", "error", err)
		return err
	}

	logger.Infoln("service listening...")
	<-ctx.Done()

	return nil
}
