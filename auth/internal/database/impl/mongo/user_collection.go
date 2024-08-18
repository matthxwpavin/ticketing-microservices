package mongo

import (
	"context"

	"github.com/matthxwpavin/ticketing/auth/internal/database"
	"github.com/matthxwpavin/ticketing/database/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type userCollection struct {
	*mongo.Collection[database.User]
}

func (c *userCollection) FindByEmail(ctx context.Context, email string) (*database.User, error) {
	return c.Collection.FindOne(ctx, bson.D{{"email", email}})
}
