package app

import (
	"context"
	"testing"

	mongodb "github.com/matthxwpavin/ticketing/database/mongo"
	"github.com/matthxwpavin/ticketing/payment/internal/database/mongo"
	"github.com/matthxwpavin/ticketing/testsetup"
)

var db *mongo.DB

var ctx context.Context

func TestMain(m *testing.M) {
	testsetup.Setup(m, &mongo.DbConfig, func(loggerCtx context.Context, connected *mongodb.DB) error {
		ctx = loggerCtx
		db = &mongo.DB{DB: connected}
		return nil
	})
}
