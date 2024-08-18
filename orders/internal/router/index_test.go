package router

import (
	"net/http"
	"testing"
	"time"

	"github.com/matthxwpavin/ticketing/httptesting"
	"github.com/matthxwpavin/ticketing/jwtclaims"
	"github.com/matthxwpavin/ticketing/jwtcookie"
	"github.com/matthxwpavin/ticketing/orders/internal/database"
	"github.com/matthxwpavin/ticketing/orderstatus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestIndex(t *testing.T) {
	indexTestCases().Run(t)
}

func indexTestCases() httptesting.TestingList {
	prepared := httptesting.Prepare(h)

	return httptesting.TestingList{
		prepared.Testing(httptesting.TestingSpecifications{
			Name: "Returns a status code other than 404",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()
				return httptesting.NewRequestGet("/api/orders")
			},
			StatusCodeFunc: func(statusCode int) bool { return statusCode != http.StatusNotFound },
		}),
		prepared.Testing(httptesting.TestingSpecifications{
			Name: "Returns 401 on an unauthenticated request",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()
				return httptesting.NewRequestGet("/api/orders")
			},
			StatusCode: http.StatusUnauthorized,
		}),
		prepared.Testing(httptesting.TestingSpecifications{
			Name: "Returns 200 on an authenicated request",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()

				return addJWT(t, httptesting.NewRequestGet("/api/orders"))
			},
			StatusCode: http.StatusOK,
		}),
		prepared.Testing(httptesting.TestingSpecifications{
			Name: "returns 3 user's orders",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()

				tr := db.TicketRepository()
				ticket1 := &database.Ticket{
					ID:    primitive.NewObjectID().Hex(),
					Title: "ticket_1",
					Price: 1000,
				}
				ticket2 := &database.Ticket{
					ID:    primitive.NewObjectID().Hex(),
					Title: "ticket_2",
					Price: 2000,
				}
				ticket3 := &database.Ticket{
					ID:    primitive.NewObjectID().Hex(),
					Title: "ticket_3",
					Price: 3000,
				}
				for _, ticket := range []*database.Ticket{ticket1, ticket2, ticket3} {
					if _, err := tr.Insert(setupCtx, ticket); err != nil {
						t.Fatalf("could not insert ticket: %v", err)
					}
				}

				or := db.OrderRepository()
				userId := primitive.NewObjectID().Hex()
				order1 := &database.TicketIdOrder{
					Order: database.Order{
						ID:        primitive.NewObjectID().Hex(),
						Status:    orderstatus.Created,
						ExpiresAt: time.Now().Add(15 * time.Minute),
						UserID:    userId,
					},
					TicketId: ticket1.ID,
				}
				order2 := &database.TicketIdOrder{
					Order: database.Order{
						ID:        primitive.NewObjectID().Hex(),
						Status:    orderstatus.Created,
						ExpiresAt: time.Now().Add(15 * time.Minute),
						UserID:    userId,
					},
					TicketId: ticket2.ID,
				}
				order3 := &database.TicketIdOrder{
					Order: database.Order{
						ID:        primitive.NewObjectID().Hex(),
						Status:    orderstatus.Created,
						ExpiresAt: time.Now().Add(15 * time.Minute),
						UserID:    userId,
					},
					TicketId: ticket3.ID,
				}
				for _, order := range []*database.TicketIdOrder{order1, order2, order3} {
					if _, err := or.Insert(setupCtx, order); err != nil {
						t.Fatalf("could not insert order: %v", err)
					}
				}
				return addJwtWithUserId(t, httptesting.NewRequestGet("/api/orders"), userId)
			},
			StatusCode: http.StatusOK,
		}),
	}
}

func addJwtWithUserId(t *testing.T, r *http.Request, userId string) *http.Request {
	jwt, err := jwtclaims.IssueToken(jwtclaims.Metadata{
		Email:  "some@example.com",
		UserID: userId,
	})
	if err != nil {
		t.Fatalf("could not issue jwt: %v", err)
	}
	r.AddCookie(jwtcookie.New(jwt))
	return r

}

func addJWT(t *testing.T, r *http.Request) *http.Request {
	return addJwtWithUserId(t, r, primitive.NewObjectID().Hex())
}
