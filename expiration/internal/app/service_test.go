package app

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/matthxwpavin/ticketing/env"
	"github.com/matthxwpavin/ticketing/logging/sugar"
)

var ctx context.Context

func TestMain(m *testing.M) {
	logger, err := sugar.New()
	if err != nil {
		log.Fatalf("could not initialize logger: %v", err)
	}

	ctx = sugar.WithContext(context.Background(), logger)
	os.Setenv("NATS_URL", "nats://nats-srv:4222")
	os.Setenv("NATS_CONN_NAME", "nats_connection_name")
	os.Setenv("REDIS_HOST", "expiration-redis-srv")
	if err := env.CheckRequiredEnvs([]env.EnvKey{
		env.NatsConnName,
		env.NatsURL,
		env.RedisHost,
	}); err != nil {
		logger.Fatalw("could not load env.", "error", err)
	}
	os.Exit(m.Run())
}
