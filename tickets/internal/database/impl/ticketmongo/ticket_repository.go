package ticketmongo

import (
	"context"

	"github.com/matthxwpavin/ticketing/database/mongo"
	"github.com/matthxwpavin/ticketing/tickets/internal/database"
	"go.mongodb.org/mongo-driver/bson"
)

type ticketRepository struct {
	*mongo.Collection[database.Ticket]
}

func (r *ticketRepository) FindAvailableTickets(ctx context.Context) ([]*database.Ticket, error) {
	return r.Find(ctx, bson.D{{"order_id", ""}})
}
