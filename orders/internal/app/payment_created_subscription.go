package app

import (
	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/orderstatus"
	"github.com/matthxwpavin/ticketing/streaming"
)

func (s *Service) handlePaymentCreated(msg *streaming.PaymentCreatedMessage, ack streaming.AckFunc) {
	logger := sugar.FromContext(s.subscriberCtx).With("msg", msg)
	logger.Infoln("received a payment created message...")

	order, err := s.or.FindByID(s.subscriberCtx, msg.OrderId)
	if err != nil {
		logger.Errorw("could not find the order", "error", err)
		return
	}
	if order == nil {
		logger.Errorw("the order not found")
		return
	}

	order.Status = orderstatus.Complete
	if err := s.or.UpdateByID(s.subscriberCtx, order.ID, order); err != nil {
		logger.Errorw("could not update the order", "error", err)
		return
	}

	if err := ack(); err != nil {
		logger.Errorw("could not ack the message", "error", err)
	}
}
