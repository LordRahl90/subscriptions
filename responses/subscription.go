package responses

import "time"

// Subscription response DTO for subscription
type Subscription struct {
	ID            string    `json:"id"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	Duration      uint      `json:"duration"`
	TrialDuration uint      `json:"trial_duration"`
	Price         float64   `json:"price"`
	Tax           float64   `json:"tax"`
	Discount      float64   `json:"discount,omitempty"`
	Total         float64   `json:"total"`
}
