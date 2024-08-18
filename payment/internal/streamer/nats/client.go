package nats

import (
	"github.com/matthxwpavin/ticketing/env"
	"github.com/matthxwpavin/ticketing/streaming/impl/nats"
	"golang.org/x/net/context"
)

func ConnectFromEnv(ctx context.Context) (*nats.Client, error) {
	return nats.Connect(ctx, env.NatsURL.Value(), env.NatsConnName.Value(), "payment-service")
}
