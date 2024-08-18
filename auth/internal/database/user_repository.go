package database

import (
	"context"
	"time"
)

type UserRepository interface {
	Insert(context.Context, *User) (string, error)
	FindByEmail(context.Context, string) (*User, error)
	DeleteAll(context.Context) error
}

type User struct {
	ID        string    `bson:"_id"`
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
