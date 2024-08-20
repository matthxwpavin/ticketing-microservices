package app

import (
	"context"

	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/streaming"
)

func (s *Service) handleOrderCreated(ctx context.Context) streaming.JsonMessageHandler[streaming.OrderCreatedMessage] {
	return func(msg *streaming.OrderCreatedMessage, ack streaming.AckFunc) {
		logger := sugar.FromContext(ctx)
		logger = logger.With("order_id", msg.OrderId, "ticket_id", msg.Ticket.Id)
		logger.Infoln("an order created event received")
		ticket, err := s.tr.FindByID(ctx, msg.Ticket.Id)
		if err != nil {
			logger.Errorw("could not find the ticket")
			return
		}
		if ticket == nil {
			logger.Errorw("the ticket not found")
			return
		}
		ticket.OrderId = msg.OrderId
		ticket.Version += 1
		if err := s.tr.UpdateByID(ctx, ticket.ID, ticket); err != nil {
			logger.Errorw("could not update the ticket", "error", err)
			return
		}
		if err := s.ticketUpdatedPub.Publish(ctx, &streaming.TicketUpdatedMessage{
			TicketID:      ticket.ID,
			TicketTitle:   ticket.Title,
			TicketPrice:   ticket.Price,
			TicketVersion: ticket.Version,
			OrderId:       ticket.OrderId,
		}); err != nil {
			logger.Errorw("could not publish a ticket updated message", "error", err)
			return
		}
		if err := ack(); err != nil {
			logger.Errorw("could not ack the message", "error", err)
		}
	}
}

func (s *Service) handleOrderCanceled(ctx context.Context) streaming.JsonMessageHandler[streaming.OrderCancelledMessage] {
	return func(msg *streaming.OrderCancelledMessage, ack streaming.AckFunc) {
		logger := sugar.FromContext(ctx).With("order_id", msg.OrderId, "order_version", msg.OrderVersion, "ticket_id", msg.Ticket.Id)

		logger.Infoln("an order cancelled message received")
		ticket, err := s.tr.FindByID(ctx, msg.Ticket.Id)
		if err != nil {
			logger.Errorw("could not find ticket", "error", err)
			return
		}
		if ticket == nil {
			logger.Errorw("the ticket not found")
			return
		}

		ticket.OrderId = ""
		ticket.Version += 1
		if err := s.tr.UpdateByID(ctx, ticket.ID, ticket); err != nil {
			logger.Errorw("could not update the ticket", "error", err)
			return
		}

		if err := ack(); err != nil {
			logger.Errorw("could not ack the message", "error", err)
		}
	}
}
