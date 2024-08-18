package router

import (
	"net/http"
	"testing"

	"github.com/matthxwpavin/ticketing/httptesting"
	"github.com/matthxwpavin/ticketing/orders/internal/database"
	"github.com/matthxwpavin/ticketing/serviceutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateOrder(t *testing.T) {
	func(t *testing.T) string {
		ticketId := primitive.NewObjectID().Hex()
		prepared := httptesting.Prepare(h)
		prepared.Testing(httptesting.TestingSpecifications{
			Name: "returns 201 if an order created successful",
			TestingRequest: func(t *testing.T) *http.Request {
				if _, err := db.TicketRepository().Insert(setupCtx, &database.Ticket{
					ID:    ticketId,
					Title: "some_title",
					Price: 230300,
				}); err != nil {
					t.Fatalf("could not insert a ticket: %v", err)
				}
				r, err := httptesting.NewRequestPostJson("/api/orders", map[string]any{
					"ticketID": ticketId,
				})
				if err != nil {
					t.Fatalf("could not new a request: %v", err)
				}
				return addJWT(t, r)
			},
			StatusCode: http.StatusCreated,
		}).Run(t)
		return ticketId
	}(t)
	// createTestCases().Run(t)
}

func createTestCases() httptesting.TestingList {
	prepared := httptesting.Prepare(h)

	successCreate := func(t *testing.T) string {
		ticketId := primitive.NewObjectID().Hex()
		prepared.Testing(httptesting.TestingSpecifications{
			Name: "returns 201 if an order created successful",
			TestingRequest: func(t *testing.T) *http.Request {
				if _, err := db.TicketRepository().Insert(setupCtx, &database.Ticket{
					ID:    ticketId,
					Title: "some_title",
					Price: 2303,
				}); err != nil {
					t.Fatalf("could not insert a ticket: %v", err)
				}
				r, err := httptesting.NewRequestPostJson("/api/orders", map[string]any{
					"ticketID": ticketId,
				})
				if err != nil {
					t.Fatalf("could not new a request: %v", err)
				}
				return addJWT(t, r)
			},
			StatusCode: http.StatusCreated,
		}).Run(t)
		return ticketId
	}

	return httptesting.TestingList{
		prepared.Testing(httptesting.TestingSpecifications{
			Name: "Returns a status code other than 404",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()

				return httptesting.NewRequestPost("/api/orders", nil)
			},
			StatusCodeFunc: func(statusCode int) bool {
				return statusCode != http.StatusNotFound && statusCode != http.StatusMethodNotAllowed
			},
		}),
		prepared.Testing(httptesting.TestingSpecifications{
			Name: "Returns 400 with invalid_parameter error type",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()

				req, err := httptesting.NewRequestPostJson("/api/orders", map[string]string{
					"ticketID": "",
				})
				if err != nil {
					t.Fatalf("unable to new a request: %v", err)
				}
				return addJWT(t, req)
			},
			StatusCode: http.StatusBadRequest,
		}).After(func(t *testing.T, r *http.Response) {
			ce, err := serviceutil.NewCustomErrorFrom(r)
			if err != nil {
				t.Fatalf("unable to parse custom error: %v", err)
			}
			if ce.Type != serviceutil.ErrTypeNameInvalidParameter {
				t.Fatalf("error type is unexpected, expected: %v, got: %v", serviceutil.ErrTypeNameInvalidParameter, ce.Type)
			}
		}),
		prepared.Testing(httptesting.TestingSpecifications{
			Name: "Returns 400 on an invalid body request",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()
				return addJWT(t, httptesting.NewRequestPost("/api/orders", nil))
			},
			StatusCode: http.StatusBadRequest,
		}),
		prepared.Testing(httptesting.TestingSpecifications{
			Name: "returns 400 if a ticket doesn't exist",
			TestingRequest: func(t *testing.T) *http.Request {
				r, err := httptesting.NewRequestPostJson("/api/orders", map[string]any{
					"ticketID": primitive.NewObjectID().Hex(),
				})
				if err != nil {
					t.Fatalf("could not new request: %v", err)
				}
				return addJWT(t, r)
			},
			StatusCode: http.StatusBadRequest,
		}),
		prepared.Testing(httptesting.TestingSpecifications{
			Name: "returns 400 if a ticket has already reserved",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()
				reservedId := successCreate(t)
				r, err := httptesting.NewRequestPostJson("/api/orders", map[string]any{
					"ticketID": reservedId,
				})
				if err != nil {
					t.Fatalf("could not new request: %v", err)
				}
				return addJWT(t, r)
			},
			StatusCode: http.StatusBadRequest,
		}),
	}
}
