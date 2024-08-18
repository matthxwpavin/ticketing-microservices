package app

import (
	"context"
	"testing"

	"github.com/matthxwpavin/ticketing/database/mongo"
	"github.com/matthxwpavin/ticketing/testsetup"
	"github.com/matthxwpavin/ticketing/tickets/internal/database/impl/ticketmongo"
)

var db *ticketmongo.DB

var ctx context.Context

func TestMain(m *testing.M) {
	testsetup.Setup(m, &ticketmongo.DbConfig, func(loggerCtx context.Context, connected *mongo.DB) error {
		ctx = loggerCtx
		db = &ticketmongo.DB{DB: connected}
		return nil
	})
}
