package mongo

import (
	"context"

	"github.com/matthxwpavin/ticketing/database/mongo"
	"github.com/matthxwpavin/ticketing/payment/internal/database"
	"github.com/matthxwpavin/ticketing/payment/internal/database/mongo/paymentcollection"
	"go.mongodb.org/mongo-driver/bson"
)

type paymentCollection struct {
	*mongo.Collection[database.Payment]
}

func (c *paymentCollection) FindByOrderId(ctx context.Context, orderId string) (*database.Payment, error) {
	return c.FindOne(ctx, bson.D{{paymentcollection.OrderId.Name, orderId}})
}

func (c *paymentCollection) FindByOrderIdAndStripePaymentIntentId(
	ctx context.Context,
	orderId string,
	stripePaymentIntentId string,
) (*database.Payment, error) {
	return c.FindOne(ctx, bson.D{{"$and", bson.A{
		bson.D{{paymentcollection.OrderId.Name, orderId}},
		bson.D{{paymentcollection.StripePaymentIntentId.Name, stripePaymentIntentId}},
	}}})
}
