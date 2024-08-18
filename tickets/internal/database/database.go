package database

type Database interface {
	TicketRepository() TicketRepository
}
