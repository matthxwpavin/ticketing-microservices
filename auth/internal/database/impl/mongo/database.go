package mongo

import (
	"context"

	"github.com/matthxwpavin/ticketing/auth/internal/database"
	"github.com/matthxwpavin/ticketing/database/mongo"
	"github.com/matthxwpavin/ticketing/env"
	"github.com/matthxwpavin/ticketing/ptr"
)

type DB struct {
	*mongo.DB
}

var DbConfig = &mongo.DbConfig{
	Name: "auth",
	Options: []*mongo.MigrationOptions{{
		CollectionName: "users",
		Validator: &mongo.Validator{
			Schema: &mongo.Schema{
				Properties: []*mongo.NamedProperty{
					{
						Name: "email",
						Property: &mongo.Property{
							BSONType:    ptr.Of(mongo.BSONTypeString),
							Description: ptr.Of("must be a string and is required"),
						},
						IsRequired: true,
					},
					{
						Name: "password",
						Property: &mongo.Property{
							BSONType:    ptr.Of(mongo.BSONTypeString),
							Description: ptr.Of("must be a string and is required"),
						},
						IsRequired: true,
					},
				},
			},
		},
	}},
}

func Connect(ctx context.Context) (*DB, error) {
	db := &mongo.DB{
		URI:    env.MongoURI.Value(),
		Config: *DbConfig,
	}
	if err := db.Connect(ctx); err != nil {
		return nil, err
	}
	return &DB{DB: db}, nil
}

func (s *DB) UserRepository() database.UserRepository {
	return &userCollection{mongo.NewCollection[database.User](s.DB, "users")}
}
