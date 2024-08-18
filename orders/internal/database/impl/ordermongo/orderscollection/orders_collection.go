package orderscollection

import (
	"github.com/matthxwpavin/ticketing/database/mongo"
	"github.com/matthxwpavin/ticketing/orderstatus"
	"github.com/matthxwpavin/ticketing/ptr"
)

const Name = "orders"

var (
	PropUserID = &mongo.NamedProperty{
		Name: "user_id",
		Property: &mongo.Property{
			BSONType:    ptr.Of(mongo.BSONTypeString),
			Description: ptr.Of("must be a string and is required"),
		},
		IsRequired: true,
	}
	PropStatus = &mongo.NamedProperty{
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
	PropExpiredAt = &mongo.NamedProperty{
		Name: "expires_at",
		Property: &mongo.Property{
			BSONType:    ptr.Of(mongo.BSONTypeDate),
			Description: ptr.Of("must be a date and is required"),
		},
		IsRequired: true,
	}
	PropTicket = &mongo.NamedProperty{
		Name: "ticket_id",
		Property: &mongo.Property{
			BSONType:    ptr.Of(mongo.BSONTypeString),
			Description: ptr.Of("must be a string and is required"),
		},
		IsRequired: true,
	}
	PropVersion = &mongo.NamedProperty{
		Name: "version",
		Property: &mongo.Property{
			BSONType:    ptr.Of(mongo.BSONTypeInt),
			Description: ptr.Of("must be an int and is required"),
		},
		IsRequired: true,
	}
)

var MigrationOptions = &mongo.MigrationOptions{
	CollectionName: Name,
	Validator: &mongo.Validator{
		Schema: &mongo.Schema{
			Properties: []*mongo.NamedProperty{PropUserID, PropStatus, PropExpiredAt, PropVersion},
		},
	},
}
