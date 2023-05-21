package vouchers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVoucherFromString(t *testing.T) {
	t.Parallel()
	table := []struct {
		name, arg string
		exp       VoucherType
	}{
		{
			name: "percent",
			arg:  "percentage",
			exp:  VoucherTypePercentage,
		},
		{
			name: "unknown percent",
			arg:  "percent",
			exp:  VoucherTypeUnknown,
		},
		{
			name: "amount",
			arg:  "amount",
			exp:  VoucherTypeAmount,
		},
		{
			name: "unknown amount",
			arg:  "amt",
			exp:  VoucherTypeUnknown,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := VoucherTypeFromString(tt.arg)
			require.Equal(t, tt.exp, got)
		})
	}
}

func TestVoucherToString(t *testing.T) {
	t.Parallel()
	table := []struct {
		name, exp string
		arg       VoucherType
	}{
		{
			name: "percent",
			arg:  VoucherTypePercentage,
			exp:  "percentage",
		},
		{
			name: "amount",
			arg:  VoucherTypeAmount,
			exp:  "amount",
		},
		{
			name: "unknown",
			arg:  VoucherTypeUnknown,
			exp:  "unknown",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.arg.String()
			require.Equal(t, tt.exp, got)
		})
	}
}

func TestCalculate(t *testing.T) {
	t.Parallel()
	table := []struct {
		name        string
		voucher     Voucher
		exp, amount float64
	}{
		{
			name: "amount",
			voucher: Voucher{
				VoucherType: VoucherTypeAmount,
				Amount:      200.0,
			},
			amount: 2000,
			exp:    200.0,
		},
		{
			name: "percentage",
			voucher: Voucher{
				VoucherType: VoucherTypePercentage,
				Percentage:  20,
			},
			amount: 200,
			exp:    40.0,
		},
		{
			name: "unknown",
			voucher: Voucher{
				VoucherType: VoucherTypeUnknown,
				Percentage:  20,
				Amount:      200,
			},
			amount: 200,
			exp:    0.0,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.voucher.Calculate(tt.amount)
			require.Equal(t, tt.exp, got)
		})
	}
}
