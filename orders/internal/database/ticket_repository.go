package database

import (
	"context"

	"github.com/matthxwpavin/ticketing/streaming"
)

type TicketRepository interface {
	Insert(context.Context, *Ticket) (string, error)
	FindByID(context.Context, string) (*Ticket, error)
	FindAll(context.Context) ([]*Ticket, error)
	DeleteByID(context.Context, string) error
	DeleteAll(context.Context) error
	UpdateByID(context.Context, string, *Ticket) error
	UpdateTicketByTicketUpdatedMessage(
		ctx context.Context,
		tcm *streaming.TicketUpdatedMessage,
	) error
}
