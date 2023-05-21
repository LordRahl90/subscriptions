package requests

import "time"

// Voucher DTO for voucher requests
type Voucher struct {
	VoucherType string    `json:"voucher_type" binding:"required"`
	ProductID   string    `json:"product_id" binding:"required"`
	Code        string    `json:"code" binding:"required"`
	Percentage  float64   `json:"percentage,omitempty"`
	Amount      float64   `json:"amount,omitempty"`
	ExpiresOn   time.Time `json:"expires_on"`
}
