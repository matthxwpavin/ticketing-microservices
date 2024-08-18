package database

import (
	"context"
	"time"
)

type TicketRepository interface {
	Insert(context.Context, *Ticket) (string, error)
	FindByID(context.Context, string) (*Ticket, error)
	FindAll(context.Context) ([]*Ticket, error)
	DeleteByID(context.Context, string) error
	DeleteAll(context.Context) error
	UpdateByID(context.Context, string, *Ticket) error
	FindAvailableTickets(context.Context) ([]*Ticket, error)
}

type Ticket struct {
	ID        string    `bson:"_id"`
	Title     string    `bson:"title,omitempty"`
	Price     int32     `bson:"price,omitempty"`
	UserID    string    `bson:"user_id,omitempty"`
	CreatedAt time.Time `bson:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty"`
	Version   int32     `bson:"version,omitempty"`
	OrderId   string    `bson:"order_id"`
}
