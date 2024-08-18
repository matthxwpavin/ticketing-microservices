package app

import (
	"context"

	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/payment/internal/database"
	"github.com/matthxwpavin/ticketing/payment/internal/streamer"
	"github.com/matthxwpavin/ticketing/payment/internal/stripe"
	"github.com/matthxwpavin/ticketing/streaming"
)

// Service is a type that represent domain/business logic
// of the application. It is high level language to communicate
// what exactly the application do.
type Service struct {
	subscriberCtx context.Context
	order         database.OrderRepository
	stripe        *stripe.Client
	payment       database.PaymentRepository
	paymentPub    streaming.PaymentCreatedPublisher
}

func NewService(
	ctx context.Context,
	db database.Database,
	s streamer.Streamer,
	stripe *stripe.Client,
) (*Service, error) {
	logger := sugar.FromContext(ctx)

	svc := &Service{
		subscriberCtx: ctx,
		order:         db.OrderRepository(),
		stripe:        stripe,
		payment:       db.PaymentRepository(),
	}

	var err error
	svc.paymentPub, err = s.PaymentCreatedPublisher(ctx)
	if err != nil {
		logger.Errorw("could not get the payment created publisher", "error", err)
		return nil, err
	}

	orderCreated, err := s.OrderCreatedConsumer(
		ctx,
		streaming.DefaultConsumeErrorHandler(ctx),
		streaming.OrderCreatedStreamConfig.Subjects[2],
	)
	if err != nil {
		logger.Errorw("could not get the order created consumer", "error", err)
		return nil, err
	}
	if _, err := orderCreated.Consume(ctx, svc.handleOrderCreated); err != nil {
		logger.Errorw("could not consume the order created", "error", err)
		return nil, err
	}

	orderCanceled, err := s.OrderCancelledConsumer(
		ctx,
		streaming.DefaultConsumeErrorHandler(ctx),
		streaming.OrderCanceledStreamConfig.Subjects[0],
	)
	if err != nil {
		logger.Errorw("could not get consume order canceled consumer", "error", err)
		return nil, err
	}
	if _, err := orderCanceled.Consume(ctx, svc.handleOrderCanceled); err != nil {
		logger.Errorw("could not consume the order canceled")
	}

	return svc, nil
}
