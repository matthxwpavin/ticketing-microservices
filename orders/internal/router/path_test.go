package router

import (
	"net/http"
	"testing"

	"github.com/matthxwpavin/ticketing/httptesting"
	"github.com/matthxwpavin/ticketing/orderstatus"
)

func TestPatchOrder(t *testing.T) {
	patchOrderTestCases().Run(t)
}

func patchOrderTestCases() httptesting.TestingList {
	prepared := httptesting.Prepare(h)
	return httptesting.TestingList{
		prepared.Testing(httptesting.TestingSpecifications{
			Name: "returns 204 if cancel an order success",
			TestingRequest: func(t *testing.T) *http.Request {
				userId, orderId := createOrder(t)
				r, err := httptesting.NewRequestPatchJson("/api/orders/"+orderId, map[string]any{
					"orderStatus": orderstatus.Cancelled,
				})
				if err != nil {
					t.Fatalf("could not new a request: %v", err)
				}
				return addJwtWithUserId(t, r, userId)
			},
			StatusCode: http.StatusNoContent,
		}),
	}
}
