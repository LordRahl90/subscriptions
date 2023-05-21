package products

import "gorm.io/gorm"

// Product contains a product's details
// Each product has their associated taxrate
type Product struct {
	ID          string  `json:"id" gorm:"primaryKey;size:32"`
	Name        string  `json:"name" gorm:"unique;size:100"`
	Description string  `json:"description"`
	TaxRate     float64 `json:"tax" gorm:"type:double(30,2)"`
	gorm.Model
}
