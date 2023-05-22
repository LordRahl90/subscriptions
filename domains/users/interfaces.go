package users

import (
	"context"
)

// IUserService manages the user service
type IUserService interface {
	Create(ctx context.Context, u *User) error
	Find(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}
