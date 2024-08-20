package app

import (
	"context"

	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/orders/internal/database"
	"github.com/matthxwpavin/ticketing/streaming"
)

func (s *Service) handleTicketCreated(ctx context.Context) streaming.JsonMessageHandler[streaming.TicketCreatedMessage] {
	return func(tcm *streaming.TicketCreatedMessage, ack streaming.AckFunc) {
		logger := sugar.FromContext(ctx).With("sub.", "ticket:created")
		logger.Infow("a ticket has been created", "id", tcm.TicketID, "title", tcm.TicketTitle, "price", tcm.TicketPrice)
		if _, err := s.tr.Insert(ctx, &database.Ticket{
			ID:      tcm.TicketID,
			Title:   tcm.TicketTitle,
			Price:   tcm.TicketPrice,
			Version: tcm.TicketVersion,
		}); err != nil {
			logger.Errorw("could not insert a ticket", "error", err)
			return
		}
		if err := ack(); err != nil {
			logger.Errorw("could not ack the message", "error", err)
		}
	}
}

func (s *Service) handleTicketUpdated(ctx context.Context) streaming.JsonMessageHandler[streaming.TicketUpdatedMessage] {
	return func(msg *streaming.TicketUpdatedMessage, ack streaming.AckFunc) {
		logger := sugar.FromContext(ctx).With("sub.", "ticket:updated")

		logger.Infow("a ticket has been updated", "id", msg.TicketID, "title", msg.TicketTitle, "price", msg.TicketPrice, "version", msg.TicketVersion)
		if err := s.tr.UpdateTicketByTicketUpdatedMessage(ctx, msg); err != nil {
			return
		}
		if err := ack(); err != nil {
			logger.Errorw("could not ack the message", "error", err)
		}
	}
}
