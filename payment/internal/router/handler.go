package router

import (
	"net/http"

	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/payment/internal/app"
	"github.com/matthxwpavin/ticketing/rw"
	"github.com/matthxwpavin/ticketing/serviceutil"
)

type handler struct {
	svc *app.Service
}

func (h *handler) createPayment(w http.ResponseWriter, r *http.Request) {
	logger := sugar.FromContext(r.Context())

	input := new(app.CreateChargeInput)
	if err := rw.DecodeJSON(r.Body, input); err != nil {
		logger.Errorw("could not decode the request body", "error", err)
		rw.Error(r.Context(), w, serviceutil.NewServiceFailureError("could not decode the request's body"))
		return
	}

	out, err := h.svc.CreateCharge(r.Context(), input)
	if err != nil {
		rw.Error(r.Context(), w, err)
		return
	}

	rw.JSON201(r.Context(), w, out)
}
