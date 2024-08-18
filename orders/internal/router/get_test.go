package router

import (
	"net/http"
	"testing"
	"time"

	"github.com/matthxwpavin/ticketing/httptesting"
	"github.com/matthxwpavin/ticketing/orders/internal/database"
	"github.com/matthxwpavin/ticketing/orderstatus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetOrder(t *testing.T) {
	getOrderTestCases().Run(t)
}

func getOrderTestCases() httptesting.TestingList {
	prepared := httptesting.Prepare(h)
	return httptesting.TestingList{
		prepared.Testing(httptesting.TestingSpecifications{
			Name: "returns 400 if the order id do not exist",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()
				return addJWT(t, httptesting.NewRequestGet("/api/orders/"+primitive.NewObjectID().Hex()))
			},
			StatusCode: http.StatusBadRequest,
		}),
		prepared.Testing(httptesting.TestingSpecifications{
			Name: "returns 401 if the user is not authenticated",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()
				return httptesting.NewRequestGet("/api/orders/" + primitive.NewObjectID().Hex())
			},
			StatusCode: http.StatusUnauthorized,
		}),
		prepared.Testing(httptesting.TestingSpecifications{
			Name: "returns 400 if the order is not owned by the user",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()
				_, orderId := createOrder(t)
				return addJWT(t, httptesting.NewRequestGet("/api/orders/"+orderId))
			},
			StatusCode: http.StatusBadRequest,
		}),
		prepared.Testing(httptesting.TestingSpecifications{
			Name: "returns an user's order",
			TestingRequest: func(t *testing.T) *http.Request {
				t.Parallel()

				userId, orderId := createOrder(t)
				return addJwtWithUserId(t, httptesting.NewRequestGet("/api/orders/"+orderId), userId)
			},
			StatusCode: http.StatusOK,
		}),
	}
}

func createOrder(t *testing.T) (string, string) {
	tr := db.TicketRepository()
	or := db.OrderRepository()

	ticketId := primitive.NewObjectID().Hex()
	if _, err := tr.Insert(setupCtx, &database.Ticket{
		ID:    ticketId,
		Title: "some_title",
		Price: 3232,
	}); err != nil {
		t.Fatalf("could not insert a ticket: %v", err)
	}

	userId := primitive.NewObjectID().Hex()
	orderId := primitive.NewObjectID().Hex()
	if _, err := or.Insert(setupCtx, &database.TicketIdOrder{
		Order: database.Order{
			ID:        orderId,
			Status:    orderstatus.Created,
			ExpiresAt: time.Now().Add(13 * time.Minute),
			UserID:    userId,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		TicketId: ticketId,
	}); err != nil {
		t.Fatalf("could not insert an order: %v", err)
	}
	return userId, orderId
}
