package users

import (
	"context"
)

type IUserService interface {
	Create(ctx context.Context, u *User) error
	Find(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}
