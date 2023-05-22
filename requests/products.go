package requests

// Product request DTO for products
type Product struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Tax         float64 `json:"tax" binding:"required"`
	TrialExists bool    `json:"trial_exists"`
}
