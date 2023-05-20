package vouchers

import "context"

// Manager manages the available operations for vouchers
type Manager interface {
	Create(ctx context.Context, v *Voucher) error
	// Find(ctx context.Context) ([]Voucher, error)
	FindOne(ctx context.Context, id string) (*Voucher, error)
}
