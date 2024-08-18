package app

import (
	"testing"
	"time"

	"github.com/matthxwpavin/ticketing/streaming"
	"github.com/matthxwpavin/ticketing/streaming/impl/nats"
	"github.com/matthxwpavin/ticketing/tickets/internal/database"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSubscribeOrderCreated(t *testing.T) {
	t.Parallel()
	// Build a new ticket then save it.
	ticket := &database.Ticket{
		ID:        primitive.NewObjectID().Hex(),
		Title:     "some_ticket",
		Price:     1234,
		UserID:    primitive.NewObjectID().Hex(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   1,
	}
	_, err := db.TicketRepository().Insert(ctx, ticket)
	require.NoError(t, err, "could not insert a ticket")

	// Fake an order created message.
	msg := &streaming.OrderCreatedMessage{
		OrderId:        primitive.NewObjectID().Hex(),
		OrderStatus:    "created",
		OrderExpiresAt: time.Now().Add(15 * time.Minute),
		OrderVersion:   1,
	}
	msg.Ticket.Id = ticket.ID
	msg.Ticket.Price = ticket.Price

	natsMock := &nats.MockClient{}
	pub, _ := natsMock.OrderCreatedPublisher(ctx)
	pub.Publish(ctx, msg)
	// Consume the order created message
	_, err = NewService(ctx, db, natsMock)
	// Expect consumption success.
	require.NoError(t, err, "could not consume the message")

	// Expect the message has been acked.
	require.Equal(t, true, natsMock.DidOrderCreatedMessageAck(), "the order created message didn't ack")

	// Find the updated ticket.
	updatedTicket, err := db.TicketRepository().FindByID(ctx, msg.Ticket.Id)
	require.NoError(t, err, "could not find the updated ticket")

	// Expect the updated ticket has an order id updated to it and must be equal to the message's order id
	require.Equal(t, msg.OrderId, updatedTicket.OrderId)
}

func TestSubscribeOrderCanceled(t *testing.T) {
	t.Parallel()

	// Build a new ticket then save it.
	ticket := &database.Ticket{
		ID:        primitive.NewObjectID().Hex(),
		Title:     "some_ticket",
		Price:     1234,
		UserID:    primitive.NewObjectID().Hex(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   2, // ticket's version is 2 here becuase the order id is not empty (an order has been created)
		OrderId:   primitive.NewObjectID().Hex(),
	}
	_, err := db.TicketRepository().Insert(ctx, ticket)
	require.NoError(t, err, "could not insert a ticket")

	// Fake an order canceled message.
	msg := &streaming.OrderCancelledMessage{
		OrderId:      ticket.OrderId,
		OrderVersion: 2, // This order's version is 2 becuase the 1st version is when the order created.
	}
	msg.Ticket.Id = ticket.ID

	// Publish the order cancelled message.
	natsMock := &nats.MockClient{}
	pub, _ := natsMock.OrderCancelledPublisher(ctx)
	pub.Publish(ctx, msg)

	// Consume the message.
	_, err = NewService(ctx, db, natsMock)
	require.NoError(t, err, "could not consume the order cancelled message")

	// Find the updated ticket by the message.
	ticket, err = db.TicketRepository().FindByID(ctx, ticket.ID)
	require.NoError(t, err, "could not find the updated ticket")

	// Expect no order id attatched to the ticket.
	require.Equal(t, "", ticket.OrderId, "the ticket's order id is not empty")
	require.Equal(t, int32(3), ticket.Version, "the ticket's version is unexpected")
	require.Equal(t, true, natsMock.DidOrderCancelledMessageAck(), "the order cancelled messsage has not been acked")
}
