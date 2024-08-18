package database

import (
	"context"
)

type OrderRepository interface {
	Insert(context.Context, *TicketIdOrder) (string, error)
	FindByID(context.Context, string) (*TicketIdOrder, error)
	FindAll(context.Context) ([]*TicketIdOrder, error)
	DeleteByID(context.Context, string) error
	DeleteAll(context.Context) error
	UpdateByID(context.Context, string, *TicketIdOrder) error
	FindTicketOrderByOrderID(context.Context, string) (*TicketOrder, error)
	FindTicketOrdersByUserId(ctx context.Context, userId string) ([]*TicketOrder, error)
	FindByTicketIdAndStatuses(
		ctx context.Context,
		ticketId string,
		statuses []string,
	) ([]*TicketIdOrder, error)
	IsTicketReserved(ctx context.Context, ticketId string) (bool, error)
}
