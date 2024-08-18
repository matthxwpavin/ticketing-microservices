package database

type Database interface {
	UserRepository() UserRepository
}
