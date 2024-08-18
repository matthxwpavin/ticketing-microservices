package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/orders/internal/app"
	"github.com/matthxwpavin/ticketing/rw"
	"github.com/matthxwpavin/ticketing/serviceutil"
)

type handler struct {
	svc *app.Service
}

func (s *handler) listOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res, err := s.svc.ListAllOrders(ctx)
	if err != nil {
		rw.Error(ctx, w, err)
		return
	}
	rw.JSON(ctx, w, res)
}

func (s *handler) createOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := sugar.FromContext(ctx)

	var input app.OrderCreate
	if err := rw.DecodeJSON(r.Body, &input); err != nil {
		logger.Errorw("unable to parse body", "error", err)
		rw.Error(ctx, w, serviceutil.NewServiceFailureError("unable to decode JSON"))
		return
	}

	created, err := s.svc.CreateOrder(ctx, input)
	if err != nil {
		rw.Error(ctx, w, err)
		return
	}
	rw.JSON201(ctx, w, created)
}

func (s *handler) updateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := sugar.FromContext(ctx)

	input := &app.UpdateOrderInput{OrderId: mux.Vars(r)[id]}
	if err := rw.DecodeJSON(r.Body, input); err != nil {
		logger.Errorw("could not decode JSON", "error", err)
		rw.Error(ctx, w, serviceutil.NewServiceFailureError("could not decode JSON"))
		return
	}
	if err := s.svc.UpdateOrder(ctx, input); err != nil {
		rw.Error(ctx, w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *handler) getOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orderId := mux.Vars(r)[id]
	order, err := s.svc.GetOrder(ctx, &app.OrderGet{OrderId: orderId})
	if err != nil {
		rw.Error(ctx, w, err)
		return
	}
	rw.JSON(ctx, w, order)
}
