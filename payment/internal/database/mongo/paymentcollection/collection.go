package paymentcollection

import (
	"github.com/matthxwpavin/ticketing/database/mongo"
	"github.com/matthxwpavin/ticketing/ptr"
)

const ColName = "payments"

var (
	Id = &mongo.NamedProperty{
		Name: "_id",
		Property: &mongo.Property{
			BSONType:    ptr.Of(mongo.BSONTypeString),
			Description: ptr.Of("id of an order"),
		},
		IsRequired: true,
	}
	OrderId = &mongo.NamedProperty{
		Name: "order_id",
		Property: &mongo.Property{
			BSONType:    ptr.Of(mongo.BSONTypeString),
			Description: ptr.Of("must be a string and is required"),
		},
		IsRequired: true,
	}
	StripePaymentIntentId = &mongo.NamedProperty{
		Name: "stripe_payment_intent_id",
		Property: &mongo.Property{
			BSONType:    ptr.Of(mongo.BSONTypeString),
			Description: ptr.Of("must be a string and is required"),
		},
		IsRequired: true,
	}
)

var MigrationOptions = &mongo.MigrationOptions{
	CollectionName: ColName,
	Validator: &mongo.Validator{
		Schema: &mongo.Schema{
			Properties: []*mongo.NamedProperty{
				Id,
				OrderId,
				StripePaymentIntentId,
			},
		},
	},
}
