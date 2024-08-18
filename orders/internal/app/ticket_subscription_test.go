package app

import (
	"testing"
	"time"

	"github.com/matthxwpavin/ticketing/orders/internal/database"
	"github.com/matthxwpavin/ticketing/streaming"
	"github.com/matthxwpavin/ticketing/streaming/impl/nats"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSubscribeTicketCreated(t *testing.T) {
	t.Parallel()

	mockClient := &nats.MockClient{}
	msg := &streaming.TicketCreatedMessage{
		TicketID:      primitive.NewObjectID().Hex(),
		TicketTitle:   "some_title",
		TicketPrice:   233,
		UserID:        primitive.NewObjectID().Hex(),
		TicketVersion: 1,
	}

	removeErr(mockClient.TicketCreatedPublisher(loggerCtx)).Publish(loggerCtx, msg)

	consumeMessage(t, mockClient)

	ticket, err := db.TicketRepository().FindByID(loggerCtx, msg.TicketID)
	require.NoError(t, err, "could not find the ticket")

	eqMatching := map[any]any{
		ticket.ID:      msg.TicketID,
		ticket.Title:   msg.TicketTitle,
		ticket.Price:   msg.TicketPrice,
		ticket.Version: msg.TicketVersion,
	}
	for got, expected := range eqMatching {
		require.Equal(t, expected, got)
	}
	require.Equal(t, true, mockClient.DidTicketCreatedMessageAck(), "ticket created messages did not ack")
}

func TestSubscribeTicketUpdated(t *testing.T) {
	t.Run("acks the ticket updated message", func(t *testing.T) {
		t.Parallel()
		ticket := buildAndSaveTicket(t)

		// Update the ticket.
		ticket.Version += 1
		ticket.Title = "updated_title"

		msg, didAck := publishAndConsumeTicketUpdatedMessage(t, ticket)
		// Find the updated ticket.
		updatedTicket, err := db.TicketRepository().FindByID(loggerCtx, msg.TicketID)

		// Finding should has no error.
		require.NoError(t, err, "could not find the ticket")
		// The updated ticket should be equal to the original.
		require.Equal(t, *ticket, *updatedTicket, "the updated ticket is not equal to the original")
		// The ticket updated message should be acked.
		require.Equal(t, true, didAck, "the updated ticket message did not ack")
	})

	t.Run("not acks the ticket updated message due to the ticket's version out of sync", func(t *testing.T) {
		t.Parallel()

		ticket := buildAndSaveTicket(t)

		// Update the version to skip one version
		ticket.Version += 2

		_, didAck := publishAndConsumeTicketUpdatedMessage(t, ticket)
		require.Equal(t, false, didAck, "the message acked")
	})

	t.Run("2 updated messages has been ack in out of orders", func(t *testing.T) {
		t.Parallel()

		ticket := buildAndSaveTicket(t)

		// Update the version to skip one version
		ticket.Version += 2
		ticket.Title = "v3"

		// First consumptin with version of 3, expect to the message has not been acked.
		_, didAck := publishAndConsumeTicketUpdatedMessage(t, ticket)
		require.Equal(t, false, didAck, "the v2 message acked")

		// Down the version to fake the ticket of second version.
		ticket.Version -= 1
		ticket.Title = "v2"

		// This time we expect the ticket should be acked.
		_, didAck = publishAndConsumeTicketUpdatedMessage(t, ticket)
		require.Equal(t, true, didAck, "the v1 message did not ack")

		// Increase the version to fake re-delivering of the previous message that has not been acked.
		ticket.Version += 1
		ticket.Title = "v3"
		// This time we expect the ticket should be acked.
		_, didAck = publishAndConsumeTicketUpdatedMessage(t, ticket)
		require.Equal(t, true, didAck, "the v2 message did not ack")

		// Find the updated ticket then expects it to identical to the v3 ticket.
		updatedTicket, err := db.TicketRepository().FindByID(loggerCtx, ticket.ID)
		require.NoError(t, err, "could not find the ticket")
		require.Equal(t, *ticket, *updatedTicket)
	})

}

func buildAndSaveTicket(t *testing.T) *database.Ticket {
	// Build and insert a ticket.
	ticket := &database.Ticket{
		ID:      primitive.NewObjectID().Hex(),
		Title:   "concert",
		Price:   234,
		Version: 1,
	}
	_, err := db.TicketRepository().Insert(loggerCtx, ticket)
	require.NoError(t, err, "could not insert a ticket")
	return ticket
}

func consumeMessage(t *testing.T, nats *nats.MockClient) {
	// Waiting for the message has been published.
	time.Sleep(time.Second * 3)
	_, err := NewService(loggerCtx, db, nats)
	if err != nil {
		t.Fatalf("could not consume the message: %v", err)
	}
	// Waiting for consumer handlers finish their tasks.
	time.Sleep(time.Second * 2)
}

func publishAndConsumeTicketUpdatedMessage(
	t *testing.T,
	ticket *database.Ticket,
) (*streaming.TicketUpdatedMessage, bool) {
	client := &nats.MockClient{}
	msg := &streaming.TicketUpdatedMessage{
		TicketID:      ticket.ID,
		TicketTitle:   ticket.Title,
		TicketPrice:   ticket.Price,
		TicketVersion: ticket.Version,
	}
	removeErr(client.TicketUpdatedPublisher(loggerCtx)).Publish(loggerCtx, msg)

	consumeMessage(t, client)
	return msg, client.DidTicketUpdatedMessageAck()
}

func removeErr[T any](res T, _ error) T {
	return res
}
