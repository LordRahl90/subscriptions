package mocks

import (
	"context"

	"subscriptions/domains/subscription"
)

var (
	_ subscription.Manager = (*SubscriptionMock)(nil)
)

// SubscriptionMock mocks subscription service
type SubscriptionMock struct {
	CreateFunc       func(ctx context.Context, s *subscription.Subscription) error
	FindFunc         func(ctx context.Context, userID string) ([]subscription.Subscription, error)
	FindByStatusFunc func(ctx context.Context, userID string, status subscription.Status) ([]subscription.Subscription, error)
	FindOneFunc      func(ctx context.Context, subID string) (*subscription.Subscription, error)
	UpdateStatusFunc func(ctx context.Context, subID string, status subscription.Status) error
}

// Create mocks creating new subscription
func (sm *SubscriptionMock) Create(ctx context.Context, s *subscription.Subscription) error {
	if sm.CreateFunc == nil {
		return errMockNotInitialized
	}
	return sm.CreateFunc(ctx, s)
}

// Find mocks finding all user subscriptions
func (sm *SubscriptionMock) Find(ctx context.Context, userID string) ([]subscription.Subscription, error) {
	if sm.FindFunc == nil {
		return nil, errMockNotInitialized
	}
	return sm.FindFunc(ctx, userID)
}

// FindByStatus mocks finding a subscription by status
func (sm *SubscriptionMock) FindByStatus(ctx context.Context, userID string, status subscription.Status) ([]subscription.Subscription, error) {
	if sm.FindByStatusFunc == nil {
		return nil, errMockNotInitialized
	}

	return sm.FindByStatusFunc(ctx, userID, status)
}

// FindOne mocks finding one subscription
func (sm *SubscriptionMock) FindOne(ctx context.Context, subID string) (*subscription.Subscription, error) {
	if sm.FindOneFunc == nil {
		return nil, errMockNotInitialized
	}
	return sm.FindOneFunc(ctx, subID)
}

// UpdateStatus mocks updating status
func (sm *SubscriptionMock) UpdateStatus(ctx context.Context, subID string, status subscription.Status) error {
	if sm.UpdateStatusFunc == nil {
		return errMockNotInitialized
	}
	return sm.UpdateStatusFunc(ctx, subID, status)
}
