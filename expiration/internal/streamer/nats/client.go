package nats

import (
	"context"

	"github.com/matthxwpavin/ticketing/env"
	"github.com/matthxwpavin/ticketing/expiration/internal/streamer"
	"github.com/matthxwpavin/ticketing/streaming/impl/nats"
)

func ConnectFromEnv(ctx context.Context) (streamer.Streamer, error) {
	return nats.Connect(
		ctx,
		env.NatsURL.Value(),
		env.NatsConnName.Value(),
		"expiration-service",
	)
}
