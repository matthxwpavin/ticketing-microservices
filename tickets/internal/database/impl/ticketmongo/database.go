package ticketmongo

import (
	"context"

	"github.com/matthxwpavin/ticketing/database/mongo"
	"github.com/matthxwpavin/ticketing/env"
	"github.com/matthxwpavin/ticketing/ptr"
	"github.com/matthxwpavin/ticketing/tickets/internal/database"
)

type DB struct {
	*mongo.DB
}

var DbConfig = mongo.DbConfig{
	Name: "tickets",
	Options: []*mongo.MigrationOptions{{
		CollectionName: "tickets",
		Validator: &mongo.Validator{
			Schema: &mongo.Schema{
				Properties: []*mongo.NamedProperty{
					{
						Name: "title",
						Property: &mongo.Property{
							BSONType:    ptr.Of(mongo.BSONTypeString),
							Description: ptr.Of("must be a string and is required"),
						},
						IsRequired: true,
					},
					{
						Name: "price",
						Property: &mongo.Property{
							BSONType:    ptr.Of(mongo.BSONTypeInt),
							Description: ptr.Of("must be an int and is required"),
						},
						IsRequired: true,
					},
					{
						Name: "user_id",
						Property: &mongo.Property{
							BSONType:    ptr.Of(mongo.BSONTypeString),
							Description: ptr.Of("must be a string and is required"),
						},
						IsRequired: true,
					},
					{
						Name: "version",
						Property: &mongo.Property{
							BSONType:    ptr.Of(mongo.BSONTypeInt),
							Description: ptr.Of("must be an integer and is required"),
						},
						IsRequired: true,
					},
					{
						Name: "order_id",
						Property: &mongo.Property{
							BSONType:    ptr.Of(mongo.BSONTypeString),
							Description: ptr.Of("an order id indicated that the ticket is being reserved if it is not empty"),
						},
						// IsRequired: true,
					},
				},
			},
		},
	}},
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
	return &ticketRepository{Collection: mongo.NewCollection[database.Ticket](s.DB, "tickets")}
}
