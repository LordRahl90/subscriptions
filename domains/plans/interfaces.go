package plans

import "context"

// Manager manages the available operations for plans
type Manager interface {
	Create(ctx context.Context, v *SubscriptionPlan) error
	Find(ctx context.Context, productID string) ([]SubscriptionPlan, error)
	FindOne(ctx context.Context, id string) (*SubscriptionPlan, error)
}
