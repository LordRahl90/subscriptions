package products

import "gorm.io/gorm"

// Product contains a product's details
// Each product has their associated taxrate
type Product struct {
	ID          string  `json:"id" gorm:"primaryKey;size:32"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	TaxRate     float64 `json:"tax"`
	gorm.Model
}
