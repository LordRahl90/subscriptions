package mocks

import (
	"context"
	"subscriptions/domains/vouchers"
)

var (
	_ vouchers.Manager = (*VoucherMocks)(nil)
)

// VoucherMocks mocks the voucher interface
type VoucherMocks struct {
	CreateFunc     func(ctx context.Context, v *vouchers.Voucher) error
	FindByCodeFunc func(ctx context.Context, code string) (*vouchers.Voucher, error)
	FindOneFunc    func(ctx context.Context, id string) (*vouchers.Voucher, error)
}

// Create implements vouchers.Manager
func (vm *VoucherMocks) Create(ctx context.Context, v *vouchers.Voucher) error {
	if vm.CreateFunc == nil {
		return errMockNotInitialized
	}
	return vm.CreateFunc(ctx, v)
}

// FindByCode implements vouchers.Manager
func (vm *VoucherMocks) FindByCode(ctx context.Context, code string) (*vouchers.Voucher, error) {
	if vm.FindByCodeFunc == nil {
		return nil, errMockNotInitialized
	}
	return vm.FindByCodeFunc(ctx, code)
}

// FindOne implements vouchers.Manager
func (vm *VoucherMocks) FindOne(ctx context.Context, id string) (*vouchers.Voucher, error) {
	if vm.FindOneFunc == nil {
		return nil, errMockNotInitialized
	}
	return vm.FindOneFunc(ctx, id)
}
