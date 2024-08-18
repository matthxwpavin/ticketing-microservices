package ordcollection

import (
	"github.com/matthxwpavin/ticketing/database/mongo"
	"github.com/matthxwpavin/ticketing/orderstatus"
	"github.com/matthxwpavin/ticketing/ptr"
)

const ColName = "orders"

var (
	Id = &mongo.NamedProperty{
		Name: "_id",
		Property: &mongo.Property{
			BSONType:    ptr.Of(mongo.BSONTypeString),
			Description: ptr.Of("id of an order"),
		},
		IsRequired: true,
	}
	Status = &mongo.NamedProperty{
		Name: "status",
		Property: &mongo.Property{
			BSONType:    ptr.Of(mongo.BSONTypeString),
			Description: ptr.Of("must be a string and is required"),
			Enum: []string{
				orderstatus.Created,
				orderstatus.Cancelled,
				orderstatus.WaitingPayment,
				orderstatus.Complete,
			},
		},
		IsRequired: true,
	}
	Version = &mongo.NamedProperty{
		Name: "version",
		Property: &mongo.Property{
			BSONType:    ptr.Of(mongo.BSONTypeInt),
			Description: ptr.Of("must be an int and is required"),
		},
		IsRequired: true,
	}
	UserId = &mongo.NamedProperty{
		Name: "user_id",
		Property: &mongo.Property{
			BSONType:    ptr.Of(mongo.BSONTypeString),
			Description: ptr.Of("must be a string and is required"),
		},
		IsRequired: true,
	}
	Price = &mongo.NamedProperty{
		Name: "price",
		Property: &mongo.Property{
			BSONType:    ptr.Of(mongo.BSONTypeInt),
			Description: ptr.Of("must be an int and is required"),
		},
	}
)

var MigrationOptions = &mongo.MigrationOptions{
	CollectionName: ColName,
	Validator: &mongo.Validator{
		Schema: &mongo.Schema{
			Properties: []*mongo.NamedProperty{
				Id, Status, Version, UserId, Price,
			},
		},
	},
}
