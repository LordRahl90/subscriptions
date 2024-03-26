package servers

import (
	"context"
	"subscriptions/domains/plans"
	"subscriptions/domains/products"
	"subscriptions/domains/subscription"
	"subscriptions/domains/users"
	"subscriptions/domains/vouchers"
)

type userService interface {
	FindByEmail(ctx context.Context, email string) (*users.User, error)
	Create(ctx context.Context, u *users.User) error
}

type productService interface {
	Find(ctx context.Context) ([]products.Product, error)
	FindOne(ctx context.Context, id string) (*products.Product, error)
	Create(ctx context.Context, p *products.Product) error
	FindByIDs(ctx context.Context, ids ...string) (map[string]products.Product, error)
}

type subscriptionService interface {
	FindOne(ctx context.Context, id string) (*subscription.Subscription, error)
	Find(ctx context.Context, userID string) ([]subscription.Subscription, error)
	UpdateStatus(ctx context.Context, subID string, status subscription.Status) error
	Create(ctx context.Context, s *subscription.Subscription) error
}

type voucherService interface {
	FindByCode(ctx context.Context, productID, code string) (*vouchers.Voucher, error)
	Create(ctx context.Context, v *vouchers.Voucher) error
}

type planService interface {
	Create(ctx context.Context, v *plans.SubscriptionPlan) error
	Find(ctx context.Context, productID string) ([]plans.SubscriptionPlan, error)
	FindOne(ctx context.Context, id string) (*plans.SubscriptionPlan, error)
}
