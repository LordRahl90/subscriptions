package subscription

import "context"

type Manager interface {
	Create(ctx context.Context, s *Subscription) error
	Find(ctx context.Context, userID string) ([]Subscription, error)
	FindOne(ctx context.Context, subID string) (*Subscription, error)
	Update(ctx context.Context, subID string)
}
