package app

import (
	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/orderstatus"
	"github.com/matthxwpavin/ticketing/ptr"
	"github.com/matthxwpavin/ticketing/streaming"
)

func (s *Service) handleOrderCanceled(msg *streaming.OrderCancelledMessage, ack streaming.AckFunc) {
	logger := sugar.FromContext(s.subscriberCtx).With("msg", msg)
	logger.Infoln("received an order canceled message")
	order, err := s.order.FindByID(s.subscriberCtx, msg.OrderId)
	if err != nil {
		logger.Errorw("could not find the order", "error", err)
		return
	}
	if order == nil {
		logger.Errorw("the order not found")
		return
	}
	order.Status = ptr.Of(orderstatus.Cancelled)
	if err := s.order.UpdateByID(s.subscriberCtx, msg.OrderId, order); err != nil {
		logger.Errorw("could not update the order", "error", err)
		return
	}
	if err := ack(); err != nil {
		logger.Errorw("could not ack the message", "error", err)
	}
}
