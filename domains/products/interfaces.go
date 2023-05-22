package products

import "context"

// Manager manages the available operations for products
type Manager interface {
	Create(ctx context.Context, p *Product) error
	Find(ctx context.Context) ([]Product, error)
	FindOne(ctx context.Context, id string) (*Product, error)
	FindByIDs(ctx context.Context, ids ...string) (map[string]Product, error)
}
