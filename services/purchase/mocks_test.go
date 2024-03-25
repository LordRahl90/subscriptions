package purchase

import (
	"context"
	"errors"

	"subscriptions/domains/plans"
	"subscriptions/domains/products"
	"subscriptions/domains/subscription"
	"subscriptions/domains/vouchers"
)

var errMockNotInitialized = errors.New("mock not initialized")

var (
	_ ProductManager      = (*ProductMock)(nil)
	_ PlanManager         = (*SubscriptionPlanMock)(nil)
	_ SubscriptionManager = (*SubscriptionMock)(nil)
	_ VoucherManager      = (*VoucherMocks)(nil)
)

// ProductMock mock products
type ProductMock struct {
	FindOneFunc   func(ctx context.Context, id string) (*products.Product, error)
	FindByIDsFunc func(ctx context.Context, ids ...string) (map[string]products.Product, error)
}

// FindOne mocks finding one product
func (pm *ProductMock) FindOne(ctx context.Context, id string) (*products.Product, error) {
	if pm.FindOneFunc == nil {
		return nil, errMockNotInitialized
	}
	return pm.FindOneFunc(ctx, id)
}

// FindByIDs mocks finding products by ids
func (pm *ProductMock) FindByIDs(ctx context.Context, ids ...string) (map[string]products.Product, error) {
	if pm.FindByIDsFunc == nil {
		return nil, errMockNotInitialized
	}
	return pm.FindByIDsFunc(ctx, ids...)
}

// SubscriptionPlanMock mocks subscription plans
type SubscriptionPlanMock struct {
	FindOneFunc func(ctx context.Context, id string) (*plans.SubscriptionPlan, error)
	CreateFunc  func(ctx context.Context, v *plans.SubscriptionPlan) error
}

func (sp *SubscriptionPlanMock) FindOne(ctx context.Context, id string) (*plans.SubscriptionPlan, error) {
	if sp.FindOneFunc == nil {
		return nil, errMockNotInitialized
	}
	return sp.FindOneFunc(ctx, id)
}

// Create mocks creating a new plan
func (sp *SubscriptionPlanMock) Create(ctx context.Context, v *plans.SubscriptionPlan) error {
	if sp.CreateFunc == nil {
		return errMockNotInitialized
	}
	return sp.CreateFunc(ctx, v)
}

// SubscriptionMock mocks subscription service
type SubscriptionMock struct {
	CreateFunc func(ctx context.Context, s *subscription.Subscription) error
}

// Create mocks creating new subscription
func (sm *SubscriptionMock) Create(ctx context.Context, s *subscription.Subscription) error {
	if sm.CreateFunc == nil {
		return errMockNotInitialized
	}
	return sm.CreateFunc(ctx, s)
}

// VoucherMocks mocks the voucher interface
type VoucherMocks struct {
	FindByCodeFunc func(ctx context.Context, productID, code string) (*vouchers.Voucher, error)
}

// FindByCode implements vouchers.Manager
func (vm *VoucherMocks) FindByCode(ctx context.Context, productID, code string) (*vouchers.Voucher, error) {
	if vm.FindByCodeFunc == nil {
		return nil, errMockNotInitialized
	}
	return vm.FindByCodeFunc(ctx, productID, code)
}
