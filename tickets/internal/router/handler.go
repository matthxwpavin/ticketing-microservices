package router

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/rw"
	"github.com/matthxwpavin/ticketing/tickets/internal/app"
	"go.uber.org/zap"
)

type handler struct {
	svc *app.Service
}

func (s *handler) listTickets(w http.ResponseWriter, r *http.Request) {
	ctx, _ := s.ctxAndLoggerFrom(r)
	res, err := s.svc.ListAllTickets(ctx)
	if err != nil {
		rw.Error(ctx, w, err)
		return
	}
	rw.JSON(ctx, w, res)
}

func (s *handler) createTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := sugar.FromContext(ctx)
	input := new(app.TicketCreate)
	if err := rw.DecodeJSON(r.Body, input); err != nil {
		logger.Errorw("Failed to decode a ticket", "error", err)
		rw.Error(ctx, w, err)
		return
	}
	res, err := s.svc.CreateTicket(ctx, input)
	if err != nil {
		rw.Error(ctx, w, err)
		return
	}
	rw.JSON201(ctx, w, res)
}

func (s *handler) updateTicket(w http.ResponseWriter, r *http.Request) {
	ctx, logger := s.ctxAndLoggerFrom(r)
	input := new(app.TicketUpdate)
	if err := rw.DecodeJSON(r.Body, input); err != nil {
		logger.Errorw("Failed to decode a ticket", "error", err)
		rw.Error(ctx, w, err)
		return
	}
	ticketID := mux.Vars(r)[id]
	if err := s.svc.UpdateTicket(ctx, ticketID, input); err != nil {
		rw.Error(ctx, w, err)
		return
	}
}

func (s *handler) deleteTicket(w http.ResponseWriter, r *http.Request) {
	ticketID := mux.Vars(r)[id]
	if err := s.svc.DeleteTicket(r.Context(), ticketID); err != nil {
		rw.Error(r.Context(), w, err)
	}
}

func (s *handler) getTicket(w http.ResponseWriter, r *http.Request) {
	ticketID := mux.Vars(r)[id]
	tk, err := s.svc.GetTicket(r.Context(), ticketID)
	if err != nil {
		rw.Error(r.Context(), w, err)
		return
	}
	rw.JSON(r.Context(), w, tk)
}

func (s *handler) ctxAndLoggerFrom(r *http.Request) (context.Context, *zap.SugaredLogger) {
	return r.Context(), sugar.FromContext(r.Context())
}
