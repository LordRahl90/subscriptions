package responses

import "time"

// SubscriptionPlan response DTO for subscription plans
type SubscriptionPlan struct {
	ID            string    `json:"id"`
	ProductID     string    `json:"product_id"`
	Amount        float64   `json:"amount"`
	Duration      uint      `json:"duration"`
	TrialDuration uint      `json:"trial_duration"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
