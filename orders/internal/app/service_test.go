package app

import (
	"context"
	"testing"

	"github.com/matthxwpavin/ticketing/database/mongo"
	"github.com/matthxwpavin/ticketing/orders/internal/database/impl/ordermongo"
	"github.com/matthxwpavin/ticketing/testsetup"
)

var db *ordermongo.DB

var loggerCtx context.Context

func TestMain(m *testing.M) {
	testsetup.Setup(m, &ordermongo.DbConfig, func(ctx context.Context, connected *mongo.DB) error {
		loggerCtx = ctx
		db = &ordermongo.DB{DB: connected}
		return nil
	})
}
