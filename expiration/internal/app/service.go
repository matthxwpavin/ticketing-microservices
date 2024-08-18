package app

import (
	"context"
	"time"

	"github.com/matthxwpavin/ticketing/expiration/internal/streamer"
	"github.com/matthxwpavin/ticketing/jsondq"
	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/streaming"
	"github.com/redis/go-redis/v9"
)

// Service is a type that represent domain/business logic
// of the application. It is high level language to communicate
// what exactly the application do.
type Service struct {
	listenerCtx            context.Context
	expirationQueue        *jsondq.Queue[ExpirationQueuePayload]
	expirationCompletedPub streaming.ExpirationCompletedPublisher
}

func NewService(ctx context.Context, s streamer.Streamer, c *redis.Client) (*Service, error) {
	logger := sugar.FromContext(ctx)

	svc := &Service{listenerCtx: ctx}
	svc.expirationQueue = jsondq.New(ctx, "order:expiration", c, svc.orderExpirationQueueHandler)
	svc.expirationQueue.StartConsume(ctx)

	orderCreatedConsumer, err := s.OrderCreatedConsumer(
		ctx,
		streaming.DefaultConsumeErrorHandler(ctx),
		streaming.OrderCreatedStreamConfig.Subjects[0],
	)
	if err != nil {
		logger.Errorw("could not get the order created consumer", "error", err)
		return nil, err
	}
	if _, err = orderCreatedConsumer.Consume(ctx, svc.orderCreatedHandler); err != nil {
		logger.Errorw("could not consume the order created event", "error", err)
		return nil, err
	}

	svc.expirationCompletedPub, err = s.ExpirationCompletedPublisher(ctx)
	if err != nil {
		logger.Errorw("could not get expiration completed publisher", "error", err)
		return nil, err
	}

	return svc, nil
}

type ExpirationQueuePayload struct {
	OrderId string `json:"orderId"`
}

func (s *Service) orderCreatedHandler(msg *streaming.OrderCreatedMessage, ack streaming.AckFunc) {
	logger := sugar.FromContext(s.listenerCtx)
	logger.Infow("received an order created message", "msg", msg)
	if err := s.expirationQueue.SendDelayMsg(&ExpirationQueuePayload{
		OrderId: msg.OrderId,
	},
		time.Until(msg.OrderExpiresAt),
	); err != nil {
		logger.Errorw("could not send an order expiration message", "error", err)
		return
	}
	if err := ack(); err != nil {
		logger.Errorw("could not ack the message", "error", err)
	}
}

func (s *Service) orderExpirationQueueHandler(msg *ExpirationQueuePayload) bool {
	logger := sugar.FromContext(s.listenerCtx)
	logger.Infow("an order expired message received", "message", msg)
	if err := s.expirationCompletedPub.Publish(s.listenerCtx, &streaming.ExpirationCompletedMessage{
		OrderId: msg.OrderId,
	}); err != nil {
		logger.Errorw("could not publish an expiration message", "error", err)
		return false
	}
	return true
}
