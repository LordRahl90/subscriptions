package servers

import (
	"context"

	"subscriptions/domains/users"
)

type userMock struct {
	CreateFunc      func(ctx context.Context, u *users.User) error
	FindFunc        func(ctx context.Context, id string) (*users.User, error)
	FindByEmailFunc func(ctx context.Context, email string) (*users.User, error)
}

// Create implements users.IUserService
func (um *userMock) Create(ctx context.Context, u *users.User) error {
	return um.CreateFunc(ctx, u)
}

// Find implements users.IUserService
func (um *userMock) Find(ctx context.Context, id string) (*users.User, error) {
	return um.FindFunc(ctx, id)
}

// FindByEmail implements users.IUserService
func (um *userMock) FindByEmail(ctx context.Context, email string) (*users.User, error) {
	return um.FindByEmailFunc(ctx, email)
}
