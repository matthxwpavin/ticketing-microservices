package app

import (
	"testing"

	"github.com/matthxwpavin/ticketing/streaming/impl/nats"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestOrderExpirationHandler(t *testing.T) {

	// Build the Service.
	natsMock := &nats.MockClient{}
	svc, err := NewService(ctx, natsMock, redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	}))
	require.NoError(t, err, "could not start service")

	// Put an expiration message to the service's queue.
	acked := svc.orderExpirationQueueHandler(&ExpirationQueuePayload{
		OrderId: primitive.NewObjectID().Hex(),
	})
	// Require the service to ack with the message.
	require.Equal(t, true, acked, "the message is not acked")

	// Require the expiration complete publisher to publish.
	require.Equal(t, true, natsMock.DidExpirationCompletedMessagePublish())

}
