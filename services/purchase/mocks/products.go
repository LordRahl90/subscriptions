package mocks

import (
	"context"

	"subscriptions/domains/products"
)

var (
	_ products.Manager = (*ProductMock)(nil)
)

// ProductMock mock products
type ProductMock struct {
	CreateFunc  func(ctx context.Context, p *products.Product) error
	FindFunc    func(ctx context.Context) ([]products.Product, error)
	FindOneFunc func(ctx context.Context, id string) (*products.Product, error)
}

// Create mocks creating a new product
func (pm *ProductMock) Create(ctx context.Context, p *products.Product) error {
	if pm.CreateFunc == nil {
		return errMockNotInitialized
	}
	return pm.CreateFunc(ctx, p)
}

// Find mocks finding all products
func (pm *ProductMock) Find(ctx context.Context) ([]products.Product, error) {
	if pm.FindFunc == nil {
		return nil, errMockNotInitialized
	}
	return pm.FindFunc(ctx)
}

// FindOne mocks finding one product
func (pm *ProductMock) FindOne(ctx context.Context, id string) (*products.Product, error) {
	if pm.FindOneFunc == nil {
		return nil, errMockNotInitialized
	}
	return pm.FindOneFunc(ctx, id)
}
