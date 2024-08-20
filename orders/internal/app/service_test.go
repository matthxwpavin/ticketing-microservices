package app

import (
	"context"
	"testing"

	"github.com/matthxwpavin/ticketing/database/mongo"
	"github.com/matthxwpavin/ticketing/orders/internal/database/impl/ordermongo"
	"github.com/matthxwpavin/ticketing/testsetup"
)

var db *ordermongo.DB

var ctx context.Context

func TestMain(m *testing.M) {
	testsetup.Setup(m, &ordermongo.DbConfig, func(loggerCtx context.Context, connected *mongo.DB) error {
		ctx = loggerCtx
		db = &ordermongo.DB{DB: connected}
		return nil
	})
}
