package router

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/middleware"
	"github.com/matthxwpavin/ticketing/payment/internal/app"
	"github.com/matthxwpavin/ticketing/payment/internal/database"
	"github.com/matthxwpavin/ticketing/payment/internal/streamer"
	"github.com/matthxwpavin/ticketing/payment/internal/stripe"
)

const (
	group = "/api/payments"

	id = "id"

	pathID = group + "/{" + id + "}"
)

func New(ctx context.Context, db database.Database, s streamer.Streamer, c *stripe.Client) *mux.Router {
	logger := sugar.FromContext(ctx)

	svc, err := app.NewService(ctx, db, s, c)
	if err != nil {
		logger.Panicw("could not build service", "error", err)
	}

	h := &handler{
		svc: svc,
	}

	r := mux.NewRouter()

	r.Use(middleware.PopulateLogger)
	r.Use(middleware.PopulateJWTClaims)

	r.HandleFunc(group, h.createPayment)

	return r
}
