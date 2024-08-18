package app

import (
	"testing"

	"github.com/matthxwpavin/ticketing/orderstatus"
	"github.com/matthxwpavin/ticketing/payment/internal/database"
	"github.com/matthxwpavin/ticketing/payment/internal/stripe"
	"github.com/matthxwpavin/ticketing/ptr"
	"github.com/matthxwpavin/ticketing/streaming"
	"github.com/matthxwpavin/ticketing/streaming/impl/nats"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestHandleOrderCanceled(t *testing.T) {
	svc, err := NewService(ctx, db, &nats.MockClient{}, stripe.NewClientTest())
	require.NoError(t, err, "could not initialize service")

	// Build an order and insert it.
	order := &database.Order{
		Id:      ptr.Of(primitive.NewObjectID().Hex()),
		Status:  ptr.Of(orderstatus.Created),
		UserId:  ptr.Of(primitive.NewObjectID().Hex()),
		Price:   ptr.Of(int32(1234)),
		Version: ptr.Of(int32(1)),
	}
	_, err = db.OrderRepository().Insert(ctx, order)
	require.NoError(t, err, "could not insert the order")

	// Emit an order canceled message.
	msg := &streaming.OrderCancelledMessage{
		OrderId: *order.Id,
	}

	// Handle the order canceled message.
	var ack bool
	svc.handleOrderCanceled(msg, func() error {
		ack = true
		return nil
	})
	require.Equal(t, true, ack, "the message has not been acked")

	// Find the order and check is status canceled.
	order, err = db.OrderRepository().FindByID(ctx, msg.OrderId)
	require.NoError(t, err, "could not find the order")
	require.NotEmpty(t, order, "the order not found")
	require.Equal(t, orderstatus.Cancelled, *order.Status)
}
