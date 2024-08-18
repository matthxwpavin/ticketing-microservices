package router

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/middleware"
	"github.com/matthxwpavin/ticketing/orders/internal/app"
	"github.com/matthxwpavin/ticketing/orders/internal/database"
	"github.com/matthxwpavin/ticketing/streaming"
)

const (
	group = "/api/orders"

	id = "id"

	pathID = group + "/{" + id + "}"
)

func New(ctx context.Context, db database.Database, s streaming.OrderStreamer) *mux.Router {
	logger := sugar.FromContext(ctx)

	svc, err := app.NewService(ctx, db, s)
	if err != nil {
		logger.Panicw("could not build service", "error", err)
	}

	h := &handler{
		svc: svc,
	}

	r := mux.NewRouter()

	r.Use(middleware.PopulateLogger)
	r.Use(middleware.PopulateJWTClaims)

	r.HandleFunc(group, h.listOrders).Methods(http.MethodGet)
	r.HandleFunc(group, h.createOrder).Methods(http.MethodPost)
	r.HandleFunc(pathID, h.getOrder).Methods(http.MethodGet)
	r.HandleFunc(pathID, h.updateOrder).Methods(http.MethodPatch)

	return r
}
