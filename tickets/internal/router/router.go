package router

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/middleware"
	"github.com/matthxwpavin/ticketing/streaming"
	"github.com/matthxwpavin/ticketing/tickets/internal/app"
	"github.com/matthxwpavin/ticketing/tickets/internal/database"
)

const (
	group = "/api/tickets"

	id = "id"

	pathID = group + "/{" + id + "}"
)

func New(ctx context.Context, db database.Database, s streaming.TicketStreamer) *mux.Router {
	logger := sugar.FromContext(ctx)

	r := mux.NewRouter()

	svc, err := app.NewService(ctx, db, s)
	if err != nil {
		logger.Panicw("could not initialize services", "error", err)
	}
	h := &handler{
		svc: svc,
	}

	r.Use(middleware.PopulateLogger)
	r.Use(middleware.PopulateJWTClaims)

	r.HandleFunc(group, h.listTickets).Methods(http.MethodGet)
	r.HandleFunc(group, h.createTicket).Methods(http.MethodPost)
	r.HandleFunc(pathID, h.updateTicket).Methods(http.MethodPut)
	r.HandleFunc(pathID, h.getTicket).Methods(http.MethodGet)
	r.HandleFunc(pathID, h.deleteTicket).Methods(http.MethodDelete)

	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	return r
}
