package responses

import "time"

// Voucher DTO for voucher responses
type Voucher struct {
	ID          string    `json:"id" gorm:"primaryKey;size:32"`
	VoucherType string    `json:"voucher_type" gorm:""`
	ProductID   string    `json:"product_id"`
	Code        string    `json:"code"`
	Active      bool      `json:"active"`
	Percentage  float64   `json:"discount"`
	Amount      float64   `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
