package app

import (
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/matthxwpavin/ticketing/orders/internal/database"
	"github.com/matthxwpavin/ticketing/orderstatus"
	"github.com/matthxwpavin/ticketing/streaming"
	"github.com/matthxwpavin/ticketing/streaming/impl/nats"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestExpirationListener(t *testing.T) {
	// Run test in parallel.
	t.Parallel()

	// Get context that can be cancel by signal.
	ctx, stop := signal.NotifyContext(loggerCtx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Insert a ticket to database.
	ticket := &database.Ticket{
		ID:      primitive.NewObjectID().Hex(),
		Title:   "ticket_tile",
		Price:   1234,
		Version: 1,
	}
	_, err := db.TicketRepository().Insert(ctx, ticket)
	require.NoError(t, err, "could not insert a ticket")

	// Insert an order to database.
	order := &database.TicketIdOrder{
		Order: database.Order{
			ID:        primitive.NewObjectID().Hex(),
			Status:    orderstatus.Created,
			ExpiresAt: time.Now().Add(expirationWindow),
			UserID:    primitive.NewObjectID().Hex(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   1,
		},
		TicketId: ticket.ID,
	}
	_, err = db.OrderRepository().Insert(ctx, order)
	require.NoError(t, err, "could not insert an order")

	// Build an expiration message.
	msg := &streaming.ExpirationCompletedMessage{OrderId: order.ID}

	// Get a publisher.
	mock := &nats.MockClient{}
	pub, err := mock.ExpirationCompletedPublisher(ctx)
	require.NoError(t, err, "could not get a publisher")

	// Publish the message.
	require.NoError(t, pub.Publish(ctx, msg), "could not publish a message")

	// Consume the message.
	_, err = NewService(ctx, db, mock)
	require.NoError(t, err, "could not cosume the message")

	// Wait for consuming then send a signal(Ctrl+C) to continue the test.
	<-ctx.Done()

	// Require the expiration message to be acked.
	require.Equal(t, true, mock.DidExpirationCompletedMessageAck(), "the expiration message has not been acked")

	// Find the updated order.
	ticketOrder, err := db.OrderRepository().FindTicketOrderByOrderID(loggerCtx, msg.OrderId)
	require.NoError(t, err, "could not find the updated order", "error", err)

	// Require to found.
	require.NotEmpty(t, ticketOrder, "the updated order not found")

	// Require the updated order's status to be canceled.
	require.Equal(t, orderstatus.Cancelled, ticketOrder.Status)

	// Require the updated order's version to be updated.
	require.Equal(t, int32(2), ticketOrder.Version)
}
