package app

import (
	"context"

	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/orders/internal/database"
	"github.com/matthxwpavin/ticketing/streaming"
)

func (s *Service) subscribeTicketCreated(ctx context.Context) error {
	logger := sugar.FromContext(ctx).With("sub.", "ticket:created")
	_, err := s.ticketCreatedSub.Consume(ctx, func(tcm *streaming.TicketCreatedMessage, ack streaming.AckFunc) {

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
	})
	if err != nil {
		logger.Errorw("could not consume", "error", err)
	}
	return err
}

func (s *Service) subscribeTicketUpdated(ctx context.Context) error {
	logger := sugar.FromContext(ctx).With("sub.", "ticket:updated")
	_, err := s.ticketUpdatedSub.Consume(ctx, func(tcm *streaming.TicketUpdatedMessage, ack streaming.AckFunc) {

		logger.Infow("a ticket has been updated", "id", tcm.TicketID, "title", tcm.TicketTitle, "price", tcm.TicketPrice, "version", tcm.TicketVersion)
		if err := s.tr.UpdateTicketByTicketUpdatedMessage(ctx, tcm); err != nil {
			return
		}
		if err := ack(); err != nil {
			logger.Errorw("could not ack the message", "error", err)
		}
	})
	if err != nil {
		logger.Errorw("could not consume", "error", err)
	}
	return err
}
