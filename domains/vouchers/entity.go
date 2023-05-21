package vouchers

import (
	"fmt"

	"gorm.io/gorm"
)

// VoucherType config
type VoucherType int

const (
	// VoucherTypeUnknown unknown voucher type
	VoucherTypeUnknown VoucherType = iota

	// VoucherTypePercentage voucher type with percentage discount
	VoucherTypePercentage

	// VoucherTypeAmount voucher type with amount discount
	VoucherTypeAmount
)

// Voucher contains the voucher details
type Voucher struct {
	ID          string      `json:"id" gorm:"primaryKey;size:32"`
	VoucherType VoucherType `json:"voucher_type" gorm:""`
	ProductID   string      `json:"product_id"`
	Code        string      `json:"code"`
	Active      bool        `json:"active"`
	Percentage  float64     `json:"percentage"`
	Amount      float64     `json:"amount"`
	Limit       uint        `json:"limit"`
	gorm.Model
}

// String returns the string representation of the voucher type
func (v VoucherType) String() string {
	switch v {
	case VoucherTypePercentage:
		return "percentage"
	case VoucherTypeAmount:
		return "amount"
	default:
		return "unknown"
	}
}

// Validate validates the voucher
func (v Voucher) Validate() error {
	if v.VoucherType == VoucherTypeUnknown {
		return fmt.Errorf("unknown voucher")
	}
	if v.VoucherType == VoucherTypeAmount && v.Amount == 0 {
		return fmt.Errorf("amount voucher should have amount")
	}

	if v.VoucherType == VoucherTypePercentage && v.Percentage == 0 {
		return fmt.Errorf("percent voucher should have the percentage")
	}

	if v.VoucherType == VoucherTypePercentage && v.Percentage > 100 {
		return fmt.Errorf("percentage cannot be more than 100")
	}

	return nil
}

// Calculate calculates the currency amount for a given voucher
func (v Voucher) Calculate(amount float64) float64 {
	if v.VoucherType == VoucherTypeUnknown {
		return 0
	}
	if v.VoucherType == VoucherTypeAmount {
		return v.Amount
	}

	return v.Percentage / 100 * amount
}

// VoucherTypeFromString returns the voucher type from a given string
func VoucherTypeFromString(s string) VoucherType {
	switch s {
	case VoucherTypeAmount.String():
		return VoucherTypeAmount
	case VoucherTypePercentage.String():
		return VoucherTypePercentage
	default:
		return VoucherTypeUnknown
	}
}
