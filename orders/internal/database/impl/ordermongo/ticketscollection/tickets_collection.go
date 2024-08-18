package ticketscollection

import (
	"github.com/matthxwpavin/ticketing/database/mongo"
	"github.com/matthxwpavin/ticketing/ptr"
)

const Name = "tickets"

var (
	PropTitle = &mongo.NamedProperty{
		Name: "title",
		Property: &mongo.Property{
			BSONType:    ptr.Of(mongo.BSONTypeString),
			Description: ptr.Of("must be a string and is required"),
		},
		IsRequired: true,
	}
	PropPrice = &mongo.NamedProperty{
		Name: "price",
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
			Properties: []*mongo.NamedProperty{PropTitle, PropPrice},
		},
	},
}
