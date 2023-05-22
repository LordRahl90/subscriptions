package requests

// SubscriptionPlan request DTO for products subscription plans
type SubscriptionPlan struct {
	ProductID     string  `json:"product_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	Duration      uint    `json:"duration" binding:"required"`
	TrialDuration uint    `json:"trial_duration"`
}
