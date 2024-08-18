package app

import (
	"github.com/matthxwpavin/ticketing/iferr"
	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/orders/internal/database"
	"github.com/matthxwpavin/ticketing/orderstatus"
	"github.com/matthxwpavin/ticketing/streaming"
)

func (s *Service) handleOrderExpiration(msg *streaming.ExpirationCompletedMessage, ack streaming.AckFunc) {
	logger := sugar.FromContext(s.subscriberCtx).With("msg", msg)
	logger.Infoln("received an order expiration message")
	ticketOrder, err := s.or.FindTicketOrderByOrderID(s.subscriberCtx, msg.OrderId)
	if err != nil {
		logger.Errorw("could not find the order", "error", err)
		return
	}
	if ticketOrder == nil {
		logger.Errorw("the order not found")
		return
	}
	if ticketOrder.Status == orderstatus.Complete {
		logger.Infoln("the order has already completed, not cancel")
		if err := ack(); err != nil {
			logger.Errorw("could not ack the message", "error", err)
		}
		return
	}

	ticketOrder.Status = orderstatus.Cancelled
	ticketOrder.Version += 1
	if err := s.or.UpdateByID(s.subscriberCtx, ticketOrder.ID, &database.TicketIdOrder{
		Order:    ticketOrder.Order,
		TicketId: ticketOrder.Ticket.ID,
	}); err != nil {
		logger.Errorw("could not update the order", "error", err)
		return
	}
	pubMsg := &streaming.OrderCancelledMessage{
		OrderId:      ticketOrder.ID,
		OrderVersion: ticketOrder.Version,
	}
	pubMsg.Ticket.Id = ticketOrder.Ticket.ID
	err = s.orderCancelledPub.Publish(s.subscriberCtx, pubMsg)

	iferr.Log(s.subscriberCtx, err, "could not publish an order canceled message: %v", err)

	iferr.Log(s.subscriberCtx, ack(), "could not ack the message: %v", err)
}
