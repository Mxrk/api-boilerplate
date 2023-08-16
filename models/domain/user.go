package domain

import (
	"context"
	"time"
)

type User struct {
	ID        int
	Username  string
	Password  string
	CreatedAt time.Time
}

// UserService represents a service for managing users.
type UserService interface {
	CreateUser(context.Context, string, string) (int, error)
	GetUser(context.Context, int) (User, error)
	GetUserFromUsername(context.Context, string) (User, bool)
	DeleteUser(context.Context, int) error
}
