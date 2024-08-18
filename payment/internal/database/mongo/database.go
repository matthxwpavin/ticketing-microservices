package mongo

import (
	"context"

	"github.com/matthxwpavin/ticketing/database/mongo"
	"github.com/matthxwpavin/ticketing/env"
	"github.com/matthxwpavin/ticketing/payment/internal/database"
	"github.com/matthxwpavin/ticketing/payment/internal/database/mongo/ordcollection"
	"github.com/matthxwpavin/ticketing/payment/internal/database/mongo/paymentcollection"
)

type DB struct {
	*mongo.DB
}

var DbConfig = mongo.DbConfig{
	Name:    "payments",
	Options: []*mongo.MigrationOptions{ordcollection.MigrationOptions},
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

func (s *DB) OrderRepository() database.OrderRepository {
	return mongo.NewCollection[database.Order](s.DB, ordcollection.ColName)
}

func (s *DB) PaymentRepository() database.PaymentRepository {
	return &paymentCollection{
		Collection: mongo.NewCollection[database.Payment](s.DB, paymentcollection.ColName),
	}
}
