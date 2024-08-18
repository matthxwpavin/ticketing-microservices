package router

import (
	"context"
	"net/http"
	"testing"

	"github.com/matthxwpavin/ticketing/database/mongo"
	"github.com/matthxwpavin/ticketing/env"
	"github.com/matthxwpavin/ticketing/orders/internal/database/impl/ordermongo"
	"github.com/matthxwpavin/ticketing/streaming/impl/nats"
	"github.com/matthxwpavin/ticketing/testsetup"
)

var db *ordermongo.DB

var h http.Handler

var setupCtx context.Context

func TestMain(m *testing.M) {
	testsetup.Setup(m, &ordermongo.DbConfig, func(ctx context.Context, connected *mongo.DB) error {

		if err := env.CheckRequiredEnvs([]env.EnvKey{
			env.MongoURI,
			env.JwtSecret,
		}); err != nil {
			return err
		}

		setupCtx = ctx
		db = &ordermongo.DB{DB: connected}
		h = New(ctx, db, &nats.MockClient{})
		return nil
	})
}
