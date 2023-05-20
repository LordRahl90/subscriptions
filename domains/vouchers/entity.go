package vouchers

import (
	"fmt"
	"subscriptions/domains/products"

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
	ID          string             `json:"id" gorm:"primaryKey;size:32"`
	VoucherType VoucherType        `json:"voucher_type" gorm:""`
	ProductID   string             `json:"product_id"`
	Products    []products.Product `json:"products" gorm:"foreignKey:ID;references:ProductID"` //one to many relationship with products
	Code        string             `json:"code"`
	Status      bool               `json:"status"`
	Percentage  float64            `json:"discount"`
	Amount      float64            `json:"amount"`
	gorm.Model
}

func (v Voucher) validate() error {
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

// VoucherUsage tracks the voucher usage
type VoucherUsage struct {
	ID             string `json:"id" gorm:"primaryKey;size:32"`
	VoucherID      string `json:"voucher_id"`
	SubscriptionID string `json:"subscription_id"`
	UserID         string `json:"user_id"`
	gorm.Model
}
