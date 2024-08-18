package database

import (
	"time"
)

type Database interface {
	OrderRepository() OrderRepository
	TicketRepository() TicketRepository
}

type Order struct {
	ID        string    `bson:"_id"`
	Status    string    `bson:"status"`
	ExpiresAt time.Time `bson:"expires_at"`
	UserID    string    `bson:"user_id,omitempty"`
	CreatedAt time.Time `bson:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty"`
	Version   int32     `bson:"version"`
}

type Ticket struct {
	ID      string `bson:"_id"`
	Title   string `bson:"title"`
	Price   int32  `bson:"price"`
	Version int32  `bson:"version"`
}

type TicketOrder struct {
	Order  `bson:"inline"`
	Ticket Ticket `bson:"ticket"`
}

type TicketIdOrder struct {
	Order    `bson:"inline"`
	TicketId string `bson:"ticket_id"`
}
