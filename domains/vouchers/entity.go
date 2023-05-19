package vouchers

import "gorm.io/gorm"

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
	ID        string  `json:"id" gorm:"primaryKey;size:32"`
	ProductID string  `json:"product_id"`
	Code      string  `json:"code"`
	Status    bool    `json:"status"`
	Discount  float64 `json:"discount"`
	Amount    float64 `json:"amount"`
	gorm.Model
}

// VoucherUsage tracks the voucher usage
type VoucherUsage struct {
	ID             string `json:"id" gorm:"primaryKey;size:32"`
	VoucherID      string `json:"voucher_id"`
	SubscriptionID string `json:"subscription_id"`
	UserID         string `json:"user_id"`
	gorm.Model
}
