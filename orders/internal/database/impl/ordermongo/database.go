package ordermongo

import (
	"context"

	"github.com/matthxwpavin/ticketing/database/mongo"
	"github.com/matthxwpavin/ticketing/env"
	"github.com/matthxwpavin/ticketing/orders/internal/database"
	"github.com/matthxwpavin/ticketing/orders/internal/database/impl/ordermongo/orderscollection"
	"github.com/matthxwpavin/ticketing/orders/internal/database/impl/ordermongo/ticketscollection"
)

type DB struct {
	*mongo.DB
}

var DbConfig = mongo.DbConfig{
	Name: "orders",
	Options: []*mongo.MigrationOptions{
		orderscollection.MigrationOptions,
		ticketscollection.MigrationOptions,
	},
}

func Connect(ctx context.Context) (*DB, error) {
	db := &mongo.DB{
		URI:    env.MongoURI.Value(),
		Config: DbConfig,
	}
	if err := db.Connect(ctx); err != nil {
		return nil, err
	}
	return &DB{DB: db}, nil
}

func (s *DB) TicketRepository() database.TicketRepository {
	return &ticketRepository{mongo.NewCollection[database.Ticket](s.DB, ticketscollection.Name)}
}

func (s *DB) OrderRepository() database.OrderRepository {
	return &orderRepository{mongo.NewCollection[database.TicketIdOrder](s.DB, orderscollection.Name)}
}
