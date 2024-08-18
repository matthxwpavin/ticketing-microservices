package router

import (
	"net/http"
	"testing"

	"github.com/matthxwpavin/ticketing/httptesting"
	"github.com/matthxwpavin/ticketing/jwtclaims"
	"github.com/matthxwpavin/ticketing/jwtcookie"
	"github.com/matthxwpavin/ticketing/rw"
	"github.com/matthxwpavin/ticketing/serviceutil"
	"github.com/matthxwpavin/ticketing/tickets/internal/app"
	"github.com/matthxwpavin/ticketing/tickets/internal/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestHandler(t *testing.T) {
	handlerTestCases().Run(t)
}

func goodTicketReq(t *testing.T) *http.Request {
	r, err := httptesting.NewRequestPostJson(group, &app.Ticket{
		Title: "valid title",
		Price: 2000,
	})
	if err != nil {
		t.Fatalf("failed to new request: %v", err)
	}
	return r
}

var (
	userID1 = primitive.NewObjectID().Hex()
	userID2 = primitive.NewObjectID().Hex()
)

func defaultAddSignInCookie1(t *testing.T, r *http.Request) *http.Request {
	return addSignInCookie(t, r, "zxcv@gmail.com", userID1)
}

func defaultAddSignInCookie2(t *testing.T, r *http.Request) *http.Request {
	return addSignInCookie(t, r, "abcd@gmail.com", userID2)
}

func addSignInCookie(t *testing.T, r *http.Request, email, userID string) *http.Request {
	signed, err := jwtclaims.IssueToken(jwtclaims.Metadata{
		Email:  email,
		UserID: userID,
	})
	if err != nil {
		t.Fatalf("failed to sign JWT, erroro: %v", err)
	}
	r.AddCookie(jwtcookie.New(signed))
	return r
}

func createGoodTicket(t *testing.T) string {
	var ticketID string
	httptesting.Run(t, httptesting.Testing{
		Handler: h,
		Specs: httptesting.TestingSpecifications{
			Name: "Create a good ticket",
			TestingRequest: func(t *testing.T) *http.Request {

				r, err := httptesting.NewRequestPostJson("/api/tickets", &app.TicketCreate{Title: "a title", Price: 3233})
				if err != nil {
					t.Fatalf("Could not create a ticket: %v", err)
				}
				return defaultAddSignInCookie1(t, r)
			},
			StatusCode: http.StatusCreated,
		},
		AfterRun: func(t *testing.T, r *http.Response) {
			ticket := new(app.Ticket)
			if err := rw.DecodeJSON(r.Body, ticket); err != nil {
				t.Fatalf("Could not decode the ticket response: %v", err)
			}
			ticketID = ticket.ID
		},
	})
	return ticketID
}

func handlerTestCases() httptesting.TestingList {
	ht := httptesting.Prepare(h)

	tc := db.TicketRepository()

	const group = "/api/tickets"
	pathID := func(id string) string {
		return group + "/" + id
	}

	goodTicket := &app.Ticket{
		Title: "a title",
		Price: 3233,
	}

	var updatedTicketID string
	return httptesting.TestingList{
		ht.Testing(httptesting.TestingSpecifications{
			Name: "Handler listens to /api/tickets for a POST request",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()
				return httptesting.NewRequestPost(group, nil)
			},
			StatusCodeFunc: func(statusCode int) bool {
				return statusCode != http.StatusNotFound
			},
		}),
		ht.Testing(httptesting.TestingSpecifications{
			Name: "Can only be accessed if the user is signed in",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()
				return goodTicketReq(t)
			},
			StatusCode: http.StatusUnauthorized,
		}),
		ht.Testing(httptesting.TestingSpecifications{
			Name: "Error if an invalid title is provided",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()
				r, err := httptesting.NewRequestPostJson(group, &app.Ticket{
					Title: "",
					Price: 1,
				})
				if err != nil {
					t.Fatalf("failed to new request: %v", err)
				}
				return defaultAddSignInCookie1(t, r)
			},
			StatusCode: http.StatusBadRequest,
		}),
		ht.Testing(httptesting.TestingSpecifications{
			Name: "Error if an invalid price is provided",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()

				r, err := httptesting.NewRequestPostJson(group, &app.Ticket{
					Title: "valid title",
					Price: -10,
				})
				if err != nil {
					t.Fatalf("failed to new request: %v", err)
				}

				return defaultAddSignInCookie1(t, r)
			},
			StatusCode: http.StatusBadRequest,
		}).After(func(t *testing.T, r *http.Response) {
			body, err := serviceutil.NewCustomErrorFrom(r)
			if err != nil {
				t.Fatalf("failed to decode to CustomError: %v", err)
			}
			if body.Type != serviceutil.ErrTypeNameInvalidParameter {
				t.Fatalf(
					"error type is not expected, expected: %s, received: %s",
					serviceutil.ErrTypeNameInvalidParameter,
					body.Type,
				)
			}

			fieldErrs, ok := body.Value.([]any)
			if !ok {
				t.Fatalf("Val is not an expected type, Val type: %T, value: %v", body.Value, body.Value)
			}
			priceErr := false
			for _, err := range fieldErrs {
				m := err.(map[string]any)
				if m["field"] == "Price" {
					priceErr = true
				}
			}
			if !priceErr {
				t.Fatalf("field 'Price' is not in errors")
			}
		}),
		ht.Testing(httptesting.TestingSpecifications{
			Name: "Create a ticket with valid inputs",
			TestingRequest: func(t *testing.T) *http.Request {
				// t.Parallel()

				tickets, err := tc.FindAll(loggerCtx)
				if err != nil {
					t.Fatalf("Could not find tickets: %v", err)
				}
				if len(tickets) > 0 {
					t.Fatalf("There are tickets left in the collection")
				}

				return defaultAddSignInCookie1(t, goodTicketReq(t))
			},
			StatusCode: http.StatusCreated,
		}).After(func(t *testing.T, r *http.Response) {
			tickets, err := tc.FindAll(loggerCtx)
			if err != nil {
				t.Fatalf("Could not find tickets after created: %v", err)
			}
			const expectedQty = 1
			if len(tickets) != expectedQty {
				t.Fatalf("tickets qty is not expected, expected: %v, received: %v", expectedQty, len(tickets))
			}
			tkt := tickets[0]
			if tkt.Price != 2000 {
				t.Fatalf("the ticket's price is not expected, expected: %v, received: %v", 2000, tkt.Price)
			}
			if tkt.Title != "valid title" {
				t.Fatalf("the ticket's title is not expected, expected: %v, received: %v", "valid title", tkt.Title)
			}

		}),
		ht.Testing(httptesting.TestingSpecifications{
			Name: "Returns other status code than 401 if an user is signed in",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()

				return defaultAddSignInCookie1(t, goodTicketReq(t))
			},
			StatusCodeFunc: func(statusCode int) bool {
				return statusCode != http.StatusUnauthorized
			},
		}),
		ht.Testing(httptesting.TestingSpecifications{
			Name: "Returns 400 if the ticket is not found",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()

				r := httptesting.NewRequestGet(pathID(primitive.NewObjectID().Hex()))
				return defaultAddSignInCookie1(t, r)
			},
			StatusCode: http.StatusBadRequest,
		}),
		ht.Testing(httptesting.TestingSpecifications{
			Name: "Returns the ticket if it is found",
			TestingRequest: func(t *testing.T) *http.Request {
				// t.Parallel()

				id := createGoodTicket(t)

				r := httptesting.NewRequestGet(pathID(id))
				return defaultAddSignInCookie1(t, r)
			},
			StatusCode: http.StatusOK,
		}).After(func(t *testing.T, r *http.Response) {
			tk := new(app.Ticket)
			if err := rw.DecodeJSON(r.Body, tk); err != nil {
				t.Fatalf("Could not decode the ticket: %v", err)
			}

			if tk.Title != goodTicket.Title {
				t.Fatalf("The title is not expected, expected: %v, received: %v", goodTicket.Title, tk.Title)
			}
			if tk.Price != goodTicket.Price {
				t.Fatalf("The price is not expected, expected: %v, received: %v", goodTicket.Price, tk.Price)
			}
		}),
		ht.Testing(httptesting.TestingSpecifications{
			Name: "Create some tickets and get them back of /api/tickets",
			TestingRequest: func(t *testing.T) *http.Request {
				tc.DeleteAll(loggerCtx)

				createGoodTicket(t)
				createGoodTicket(t)
				createGoodTicket(t)

				return defaultAddSignInCookie1(t, httptesting.NewRequestGet(group))
			},
			StatusCode: http.StatusOK,
		}).After(func(t *testing.T, r *http.Response) {
			var tickets []app.Ticket
			if err := rw.DecodeJSON(r.Body, &tickets); err != nil {
				t.Fatalf("Could not decode tickets: %v", err)
			}

			const qty = 3
			if len(tickets) != qty {
				t.Fatalf("Tickets qty is not expected, expected: %v, got: %v", qty, len(tickets))
			}
		}),
		ht.Testing(httptesting.TestingSpecifications{
			Name: "Update a ticket",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()

				updatedTicketID = createGoodTicket(t)

				r, err := httptesting.NewRequestPutJson(pathID(updatedTicketID), &app.Ticket{
					Title: "updated",
					Price: 1000000000,
				})
				if err != nil {
					t.Fatalf("Could not new update ticket request: %v", err)
				}
				return defaultAddSignInCookie1(t, r)
			},
			StatusCode: http.StatusOK,
		}).After(func(t *testing.T, r *http.Response) {
			_, err := tc.FindByID(loggerCtx, updatedTicketID)
			if err != nil {
				t.Fatalf("Could not decode the updated ticket: %v", err)
			}

		}),
		ht.Testing(httptesting.TestingSpecifications{
			Name: "Delete a ticket",
			TestingRequest: func(t *testing.T) *http.Request {

				if err := tc.DeleteAll(loggerCtx); err != nil {
					t.Fatalf("Could not drop tickets, error: %v", err)
				}
				id := createGoodTicket(t)

				r := httptesting.NewRequestDelete(pathID(id))
				return defaultAddSignInCookie1(t, r)
			},
			StatusCode: http.StatusOK,
		}).After(func(t *testing.T, r *http.Response) {
			tickets, err := tc.FindAll(loggerCtx)
			if err != nil {
				t.Fatalf("Could not find tickets: %v", err)
			}

			const ticketsQty = 0
			if len(tickets) != ticketsQty {
				t.Fatalf("Tickets Qty is unexpected, expected: %v, got: %v", ticketsQty, len(tickets))
			}
		}),
		ht.Testing(httptesting.TestingSpecifications{
			Name: "Returns 404 if a path not match",
			TestingRequest: func(t *testing.T) *http.Request {
				return httptesting.NewRequestGet("/abcd")
			},
			StatusCode: http.StatusNotFound,
		}),
		ht.Testing(httptesting.TestingSpecifications{
			Name: "Returns 400 when try to update a ticket that has no permission on it",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()

				if err := tc.DeleteAll(loggerCtx); err != nil {
					t.Fatalf("Could not drop tickets: %v", err)
				}

				id := createGoodTicket(t)

				r, err := httptesting.NewRequestPutJson(pathID(id), &app.Ticket{
					Title: "updated",
					Price: 1000000000,
				})
				if err != nil {
					t.Fatalf("Could not new update ticket request: %v", err)
				}
				return defaultAddSignInCookie2(t, r)
			},
			StatusCode: http.StatusBadRequest,
		}),
		ht.Testing(httptesting.TestingSpecifications{
			Name: "Returns 400 when try to delete a ticket that has no permission on it",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()

				if err := tc.DeleteAll(loggerCtx); err != nil {
					t.Fatalf("Could not drop tickets: %v", err)
				}

				id := createGoodTicket(t)

				r := httptesting.NewRequestDelete(pathID(id))
				return defaultAddSignInCookie2(t, r)
			},
			StatusCode: http.StatusBadRequest,
		}),
		ht.Testing(httptesting.TestingSpecifications{
			Name: "Update a reserved ticket",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()

				ticket := &database.Ticket{
					ID:      primitive.NewObjectID().Hex(),
					Title:   "some_title",
					Price:   123,
					UserID:  primitive.NewObjectID().Hex(),
					Version: 2,
					OrderId: primitive.NewObjectID().Hex(),
				}
				if _, err := db.TicketRepository().Insert(loggerCtx, ticket); err != nil {
					t.Fatalf("could not insert an ordered ticket: %v", err)
				}
				r, err := httptesting.NewRequestPutJson(pathID(ticket.ID), &app.Ticket{
					Title: "updated",
					Price: 1000000000,
				})
				if err != nil {
					t.Fatalf("Could not new update ticket request: %v", err)
				}
				return addSignInCookie(t, r, "some@email.com", ticket.UserID)
			},
			StatusCode: http.StatusBadRequest,
		}),
	}
}
