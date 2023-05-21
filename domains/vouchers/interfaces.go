package vouchers

import "context"

// Manager manages the available operations for vouchers
type Manager interface {
	Create(ctx context.Context, v *Voucher) error
	FindByCode(ctx context.Context, productID, code string) (*Voucher, error)
	FindOne(ctx context.Context, id string) (*Voucher, error)
}
