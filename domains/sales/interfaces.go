package sales

import "context"

type Manager interface {
	Calculate(ctx context.Context, req *PurchaseRequest) (float64, error)
	Purchase(ctx context.Context, req *PurchaseRequest) error
}
