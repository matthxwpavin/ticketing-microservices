package router

import (
	"context"
	"net/http"
	"testing"

	"github.com/matthxwpavin/ticketing/database/mongo"
	"github.com/matthxwpavin/ticketing/streaming/impl/nats"
	"github.com/matthxwpavin/ticketing/testsetup"
	"github.com/matthxwpavin/ticketing/tickets/internal/database/impl/ticketmongo"
)

var db *ticketmongo.DB

var h http.Handler

var loggerCtx context.Context

func TestMain(m *testing.M) {
	testsetup.Setup(m, &ticketmongo.DbConfig, func(ctx context.Context, connected *mongo.DB) error {

		loggerCtx = ctx
		db = &ticketmongo.DB{DB: connected}
		h = New(ctx, db, &nats.MockClient{})
		return nil
	})
}
