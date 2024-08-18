package router

import (
	"net/http"
	"testing"

	"github.com/matthxwpavin/ticketing/httptesting"
	"github.com/matthxwpavin/ticketing/orderstatus"
	"github.com/matthxwpavin/ticketing/payment/internal/database"
	"github.com/matthxwpavin/ticketing/payment/internal/stripe"
	"github.com/matthxwpavin/ticketing/ptr"
	"github.com/matthxwpavin/ticketing/rw"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateCharge(t *testing.T) {
	t.Run("valid charge creation", func(t *testing.T) {
		t.Parallel()
		// Build an order then insert it
		order := buildAndSaveOrder(t)
		// Do a HTTP request to /api/payments endpoint.
		httptesting.Run(t, httptesting.Testing{
			Handler: h,
			Specs: httptesting.TestingSpecifications{
				Name: "valid charge created",
				TestingRequest: func(t *testing.T) *http.Request {
					r, err := httptesting.NewRequestPostJson("/api/payments", map[string]any{
						"orderId": order.Id,
					})
					require.NoError(t, err, "could not build a request")
					return addJwtCookieWithUserId(t, r, *order.UserId)
				},
				StatusCode: http.StatusCreated,
			},
			AfterRun: func(t *testing.T, r *http.Response) {
				// Parse the response into map.
				var createChargeOutput map[string]any
				require.NoError(t, rw.DecodeJSON(r.Body, &createChargeOutput), "could not decode the response")
				// Get the created payment intent id.
				createdChargeId := createChargeOutput["paymentIntentId"].(string)

				// Use that id to retreive from Stripe.
				paymentIntent, err := stripe.NewClient().GetPaymentIntentById(ctx, createdChargeId)
				require.NoError(t, err, "could not fetch the payment intent")

				// Check the payment intent is not empty.
				require.NotEmpty(t, paymentIntent, "the payment intent is empty")
				// Check the retreived payment intent id is the same to the created one.
				require.Equal(t, createdChargeId, paymentIntent.ID)
				// Check the price amount is the same the created order.
				require.Equal(t, int64(*order.Price), paymentIntent.Amount, "the payment intent amount is not equal the order amount")

				// Find the created payment by the order id and Strip payment intent id.
				payment, err := db.PaymentRepository().FindByOrderIdAndStripePaymentIntentId(
					ctx,
					*order.Id,
					paymentIntent.ID,
				)
				require.NoError(t, err, "could not find the payment")

				// Check the payment is not nil.
				require.NotEmpty(t, payment)
				// Check the payment's order id is equal to the order id.
				require.Equal(t, order.Id, payment.OrderId)
				// Check the payment's Stripe payment intent id is equal to Stripe payment intent id.
				require.Equal(t, paymentIntent.ID, *payment.StripePaymentIntentId)
			},
		})
	})

	t.Run("no permission to create", func(t *testing.T) {
		t.Parallel()

		// Build an order then insert it
		order := buildAndSaveOrder(t)

		httptesting.Run(t, httptesting.Testing{
			Handler: h,
			Specs: httptesting.TestingSpecifications{
				Name: "no permistion to create",
				TestingRequest: func(t *testing.T) *http.Request {
					r, err := httptesting.NewRequestPostJson("/api/payments", map[string]any{
						"orderId": order.Id,
					})
					require.NoError(t, err, "could not build a request")
					return addJwtCookie(t, r)
				},
				StatusCode: http.StatusBadRequest,
			},
		})
	})

	t.Run("request on order's status canceled", func(t *testing.T) {
		order := buildAndSaveOrderWithStatus(t, orderstatus.Cancelled)

		httptesting.Run(t, httptesting.Testing{
			Handler: h,
			Specs: httptesting.TestingSpecifications{
				Name: "request on order's status canceled",
				TestingRequest: func(t *testing.T) *http.Request {
					r, err := httptesting.NewRequestPostJson("/api/payments", map[string]any{
						"orderId": order.Id,
					})
					require.NoError(t, err, "could not build a request")
					return addJwtCookieWithUserId(t, r, *order.UserId)
				},
				StatusCode: http.StatusBadRequest,
			},
		})
	})
}

func buildAndSaveOrder(t *testing.T) *database.Order {
	return buildAndSaveOrderWithStatus(t, orderstatus.Created)
}

func buildAndSaveOrderWithStatus(t *testing.T, orderstatus string) *database.Order {
	order := &database.Order{
		Id:      ptr.Of(primitive.NewObjectID().Hex()),
		Version: ptr.Of(int32(1)),
		Status:  &orderstatus,
		UserId:  ptr.Of(primitive.NewObjectID().Hex()),
		Price:   ptr.Of(int32(122)),
	}
	_, err := db.OrderRepository().Insert(ctx, order)
	require.NoError(t, err, "could not insert the order")
	return order
}
