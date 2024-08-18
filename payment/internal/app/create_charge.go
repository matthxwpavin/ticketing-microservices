package app

import (
	"context"

	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/orderstatus"
	"github.com/matthxwpavin/ticketing/payment/internal/database"
	"github.com/matthxwpavin/ticketing/payment/internal/stripe"
	"github.com/matthxwpavin/ticketing/ptr"
	"github.com/matthxwpavin/ticketing/serviceutil"
	"github.com/matthxwpavin/ticketing/streaming"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateChargeInput struct {
	OrderId string `json:"orderId"`
}

type CreateOrderOutput struct {
	PaymentIntentId string `json:"paymentIntentId"`
	ClientSecret    string `json:"clientSecret"`
	PaymentId       string `json:"paymentId"`
}

func (s *Service) CreateCharge(ctx context.Context, input *CreateChargeInput) (*CreateOrderOutput, error) {
	logger := sugar.FromContext(ctx).With("input", input)
	logger.Infoln("creating charge...")
	if err := serviceutil.ValidateStruct(input); err != nil {
		logger.Errorw("validate input failed", "error", err)
		return nil, err
	}
	user, err := serviceutil.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	order, err := s.order.FindByID(ctx, input.OrderId)
	if err != nil {
		logger.Errorw("could not find the order", "error", err)
		return nil, err
	}
	if order == nil {
		logger.Errorw("the order not found")
		return nil, serviceutil.NewServiceFailureError("the order not found")
	}
	if *order.UserId != user.UserID {
		logger.Errorw("the user has no permission", "order user id", *order.UserId, "user id", user.UserID)
		return nil, serviceutil.NewServiceFailureError("the user has no permission")
	}
	if *order.Status == orderstatus.Cancelled {
		logger.Errorln("the order is canceled")
		return nil, serviceutil.NewServiceFailureError("the order is canceled")
	}

	paymentIntent, err := s.stripe.Charge(ctx, &stripe.ChargeUsdInput{stripe.BaseChargeInput{
		Amount: ptr.Of(int64(*order.Price)),
	}})
	if err != nil {
		logger.Errorw("could not create a charge", "error", err)
		return nil, serviceutil.NewServiceFailureError("could not create a charge")
	}
	paymentId, err := s.payment.Insert(ctx, &database.Payment{
		Id:                    ptr.Of(primitive.NewObjectID().Hex()),
		OrderId:               order.Id,
		StripePaymentIntentId: &paymentIntent.ID,
	})
	if err != nil {
		logger.Errorw("could not insert a payment", "error", err)
		return nil, serviceutil.NewServiceFailureError("could not create the payment")
	}
	if err := s.paymentPub.Publish(ctx, &streaming.PaymentCreatedMessage{
		PaymentId:             paymentId,
		OrderId:               *order.Id,
		StripePaymentIntentId: paymentIntent.ID,
	}); err != nil {
		logger.Errorw("could not publish a payment created message", "error", err)
		return nil, serviceutil.NewServiceFailureError("could not publish a payment created message")
	}
	return &CreateOrderOutput{
		ClientSecret:    paymentIntent.ClientSecret,
		PaymentIntentId: paymentIntent.ID,
		PaymentId:       paymentId,
	}, nil
}
