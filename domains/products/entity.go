package products

import "gorm.io/gorm"

// Product contains a product's details
// Each product has their associated tax
type Product struct {
	ID          string  `json:"id" gorm:"primaryKey;size:32"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Tax         float64 `json:"tax"`
	TrialExists bool    `json:"trial_exists"`
	gorm.Model
}
