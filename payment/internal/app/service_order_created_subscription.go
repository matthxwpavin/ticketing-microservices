package app

import (
	"github.com/matthxwpavin/ticketing/iferr"
	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/payment/internal/database"
	"github.com/matthxwpavin/ticketing/streaming"
)

func (s *Service) handleOrderCreated(msg *streaming.OrderCreatedMessage, ack streaming.AckFunc) {
	logger := sugar.FromContext(s.subscriberCtx).With("msg", msg)
	logger.Infoln("an order created event received")

	if _, err := s.order.Insert(s.subscriberCtx, &database.Order{
		Id:      &msg.OrderId,
		Version: &msg.OrderVersion,
		Status:  &msg.OrderStatus,
		UserId:  &msg.OrderUserId,
		Price:   &msg.Ticket.Price,
	}); err != nil {
		logger.Errorw("could not insert an order", "error", err)
		return
	}
	err := ack()
	iferr.Log(s.subscriberCtx, err, "could not ack the message: %v", err)
}
