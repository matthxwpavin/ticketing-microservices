package ordermongo

import (
	"context"
	"fmt"

	"github.com/matthxwpavin/ticketing/database/mongo"
	"github.com/matthxwpavin/ticketing/logging/sugar"
	"github.com/matthxwpavin/ticketing/orders/internal/database"
	"github.com/matthxwpavin/ticketing/streaming"
)

type ticketRepository struct {
	*mongo.Collection[database.Ticket]
}

func (r *ticketRepository) UpdateTicketByTicketUpdatedMessage(
	ctx context.Context,
	tcm *streaming.TicketUpdatedMessage,
) error {
	logger := sugar.FromContext(ctx)
	ticket, err := r.FindByID(ctx, tcm.TicketID)
	if err != nil {
		logger.Errorw("could not find the ticket", "error", err)
		return err
	}
	if ticket == nil {
		logger.Errorln("the ticket not found")
		return err
	}
	if ticket.Version+1 != tcm.TicketVersion {
		err := fmt.Errorf("the ticket's version %v out of sync to update with the message's version: %v", ticket.Version, tcm.TicketVersion)
		logger.Errorln(err)
		return err
	}
	if err := r.UpdateByID(ctx, tcm.TicketID, &database.Ticket{
		ID:      tcm.TicketID,
		Title:   tcm.TicketTitle,
		Price:   tcm.TicketPrice,
		Version: tcm.TicketVersion,
	}); err != nil {
		logger.Errorw("could not update a ticket", "error", err)
		return err
	}
	return nil
}
