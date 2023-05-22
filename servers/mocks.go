package servers

import (
	"context"

	"subscriptions/domains/users"
)

var _ users.IUserService = (*userMock)(nil)

type userMock struct {
	CreateFunc      func(ctx context.Context, u *users.User) error
	FindFunc        func(ctx context.Context, id string) (*users.User, error)
	FindByEmailFunc func(ctx context.Context, email string) (*users.User, error)
}

// Create implements users.IUserService
func (*userMock) Create(ctx context.Context, u *users.User) error {
	panic("unimplemented")
}

// Find implements users.IUserService
func (*userMock) Find(ctx context.Context, id string) (*users.User, error) {
	panic("unimplemented")
}

// FindByEmail implements users.IUserService
func (*userMock) FindByEmail(ctx context.Context, email string) (*users.User, error) {
	panic("unimplemented")
}
