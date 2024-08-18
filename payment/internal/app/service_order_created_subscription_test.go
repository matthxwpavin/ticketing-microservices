package app

import (
	"testing"
	"time"

	"github.com/matthxwpavin/ticketing/orderstatus"
	"github.com/matthxwpavin/ticketing/payment/internal/stripe"
	"github.com/matthxwpavin/ticketing/streaming"
	"github.com/matthxwpavin/ticketing/streaming/impl/nats"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestHandleOrderCreated(t *testing.T) {
	svc, err := NewService(ctx, db, &nats.MockClient{}, stripe.NewClientTest())
	require.NoError(t, err, "could not initialize service")

	msg := &streaming.OrderCreatedMessage{
		OrderId:        primitive.NewObjectID().Hex(),
		OrderStatus:    orderstatus.Created,
		OrderVersion:   1,
		OrderExpiresAt: time.Now().Add(15 * time.Minute),
		OrderUserId:    primitive.NewObjectID().Hex(),
	}
	msg.Ticket.Id = primitive.NewObjectID().Hex()
	msg.Ticket.Price = 123

	var ack bool
	svc.handleOrderCreated(msg, func() error {
		ack = true
		return nil
	})
	require.Equal(t, true, ack, "the message has not been acked")

	order, err := db.OrderRepository().FindByID(ctx, msg.OrderId)
	require.NoError(t, err, "could not find the order")
	require.NotEmpty(t, order, "the order not found")
}
