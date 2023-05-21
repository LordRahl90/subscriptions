package subscription

import "context"

// Manager manages the structure for subscriptions
type Manager interface {
	Create(ctx context.Context, s *Subscription) error
	Find(ctx context.Context, userID string) ([]Subscription, error)
	FindByStatus(ctx context.Context, userID string, status Status) ([]Subscription, error)
	FindOne(ctx context.Context, subID string) (*Subscription, error)
	UpdateStatus(ctx context.Context, subID string, status Status) error
}
