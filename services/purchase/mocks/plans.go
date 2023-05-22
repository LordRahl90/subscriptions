package mocks

import (
	"context"

	"subscriptions/domains/plans"
)

var (
	_ plans.Manager = (*SubscritionPlanMock)(nil)
)

// SubscritionPlanMock mocks subscription plans
type SubscritionPlanMock struct {
	CreateFunc  func(ctx context.Context, v *plans.SubscriptionPlan) error
	FindFunc    func(ctx context.Context, productID string) ([]plans.SubscriptionPlan, error)
	FindOneFunc func(ctx context.Context, id string) (*plans.SubscriptionPlan, error)
}

// Create mocks creating a new plan
func (sp *SubscritionPlanMock) Create(ctx context.Context, v *plans.SubscriptionPlan) error {
	if sp.CreateFunc == nil {
		return errMockNotInitialized
	}
	return sp.CreateFunc(ctx, v)
}

// Find mocks finding all plans for a given product
func (sp *SubscritionPlanMock) Find(ctx context.Context, productID string) ([]plans.SubscriptionPlan, error) {
	if sp.FindFunc == nil {
		return nil, errMockNotInitialized
	}
	return sp.FindFunc(ctx, productID)
}

// FindOne mocks finding one plan
func (sp *SubscritionPlanMock) FindOne(ctx context.Context, id string) (*plans.SubscriptionPlan, error) {
	if sp.FindOneFunc == nil {
		return nil, errMockNotInitialized
	}
	return sp.FindOneFunc(ctx, id)
}
