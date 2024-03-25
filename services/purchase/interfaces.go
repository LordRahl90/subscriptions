package purchase

import (
	"context"

	"subscriptions/domains/plans"
	"subscriptions/domains/products"
	"subscriptions/domains/subscription"
	"subscriptions/domains/vouchers"
)

type ProductManager interface {
	FindOne(ctx context.Context, id string) (*products.Product, error)
}

type PlanManager interface {
	FindOne(ctx context.Context, id string) (*plans.SubscriptionPlan, error)
}

type VoucherManager interface {
	FindByCode(ctx context.Context, productID, code string) (*vouchers.Voucher, error)
}

type SubscriptionManager interface {
	Create(ctx context.Context, s *subscription.Subscription) error
}
