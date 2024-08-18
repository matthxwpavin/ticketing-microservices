package app

import (
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/streaming"
	"github.com/matthxwpavin/ticketing/streaming/impl/nats"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestOrderExpirationHandler(t *testing.T) {
	// Create a mock NATS client.
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	natsMock := &nats.MockClient{}

	// Create a redis client.
	redisCli := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	pub, err := natsMock.OrderCreatedPublisher(ctx)
	require.NoError(t, err, "could not get the ticket created publisher")

	msg := &streaming.OrderCreatedMessage{
		OrderId:        primitive.NewObjectID().Hex(),
		OrderStatus:    "created",
		OrderExpiresAt: time.Now().Add(time.Minute * 15),
	}
	msg.Ticket.Id = primitive.NewObjectID().Hex()
	msg.Ticket.Price = 12121
	err = pub.Publish(ctx, msg)
	require.NoError(t, err, "could not publish a message")

	// Consume the messages.
	_, err = NewService(ctx, natsMock, redisCli)
	require.NoError(t, err, "could not start service")

	// After x minutes since an order created, an expiration completed message will be published.
	// .
	// .
	// .

	// Get an expiration completed consumer
	consumer, err := natsMock.ExpirationCompletedConsumer(ctx, streaming.DefaultConsumeErrorHandler(ctx), "")
	require.NoError(t, err, "could not get the expiration completed consumer")

	// Consume the expiration completed message.
	_, err = consumer.Consume(ctx, func(msg *streaming.ExpirationCompletedMessage, ack streaming.AckFunc) {
		sugar.FromContext(ctx).Infow("received an order expiration completed message", "msg", msg)
		ack()
	})
	require.NoError(t, err, "could not consume the expiration completed message")

	// Wait for terminating signal to finish the test.
	<-ctx.Done()

	// Require the order created message has been acked.
	require.Equal(t, true, natsMock.DidOrderCreatedMessageAck(), "the order created message was not acked")

	// Require the expiration completed message has been acked.
	require.Equal(t, true, natsMock.DidExpirationCompletedMessageAck(), "the expiration completed message was not acked")
}
