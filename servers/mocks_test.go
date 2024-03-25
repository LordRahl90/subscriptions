package servers

import (
	"context"

	"subscriptions/domains/plans"
	"subscriptions/domains/products"
	"subscriptions/domains/subscription"
	"subscriptions/domains/vouchers"
)

type productServiceMock struct {
	FindFunc      func(ctx context.Context) ([]products.Product, error)
	FindOneFunc   func(ctx context.Context, id string) (*products.Product, error)
	CreateFunc    func(ctx context.Context, p *products.Product) error
	FindByIDsFunc func(ctx context.Context, ids ...string) (map[string]products.Product, error)
}

func (ps *productServiceMock) Find(ctx context.Context) ([]products.Product, error) {
	return ps.FindFunc(ctx)
}

func (ps *productServiceMock) FindOne(ctx context.Context, id string) (*products.Product, error) {
	return ps.FindOneFunc(ctx, id)
}

func (ps *productServiceMock) Create(ctx context.Context, p *products.Product) error {
	return ps.CreateFunc(ctx, p)
}

func (ps *productServiceMock) FindByIDs(ctx context.Context, ids ...string) (map[string]products.Product, error) {
	return ps.FindByIDsFunc(ctx, ids...)
}

type subscriptionServiceMock struct {
	CreateFunc       func(ctx context.Context, s *subscription.Subscription) error
	FindOneFunc      func(ctx context.Context, id string) (*subscription.Subscription, error)
	FindFunc         func(ctx context.Context, userID string) ([]subscription.Subscription, error)
	UpdateStatusFunc func(ctx context.Context, subID string, status subscription.Status) error
}

func (ss *subscriptionServiceMock) Create(ctx context.Context, s *subscription.Subscription) error {
	return ss.CreateFunc(ctx, s)
}

func (ss *subscriptionServiceMock) FindOne(ctx context.Context, id string) (*subscription.Subscription, error) {
	return ss.FindOneFunc(ctx, id)
}

func (ss *subscriptionServiceMock) Find(ctx context.Context, userID string) ([]subscription.Subscription, error) {
	return ss.FindFunc(ctx, userID)
}

func (ss *subscriptionServiceMock) UpdateStatus(ctx context.Context, subID string, status subscription.Status) error {
	return ss.UpdateStatusFunc(ctx, subID, status)
}

type voucherServiceMock struct {
	CreateFunc     func(ctx context.Context, v *vouchers.Voucher) error
	FindByCodeFunc func(ctx context.Context, productID, code string) (*vouchers.Voucher, error)
}

func (vs *voucherServiceMock) Create(ctx context.Context, v *vouchers.Voucher) error {
	return vs.CreateFunc(ctx, v)
}

func (vs *voucherServiceMock) FindByCode(ctx context.Context, productID, code string) (*vouchers.Voucher, error) {
	return vs.FindByCodeFunc(ctx, productID, code)
}

type planServiceMock struct {
	CreateFunc  func(ctx context.Context, v *plans.SubscriptionPlan) error
	FindFunc    func(ctx context.Context, productID string) ([]plans.SubscriptionPlan, error)
	FindOneFunc func(ctx context.Context, id string) (*plans.SubscriptionPlan, error)
}

func (ps *planServiceMock) Find(ctx context.Context, productID string) ([]plans.SubscriptionPlan, error) {
	return ps.FindFunc(ctx, productID)
}

func (ps *planServiceMock) FindOne(ctx context.Context, id string) (*plans.SubscriptionPlan, error) {
	return ps.FindOneFunc(ctx, id)
}

func (ps *planServiceMock) Create(ctx context.Context, v *plans.SubscriptionPlan) error {
	return ps.CreateFunc(ctx, v)
}
