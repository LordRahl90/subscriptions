package plans

import "gorm.io/gorm"

// SubscriptionPlan contains the details of a subscription plan
type SubscriptionPlan struct {
	ID            string  `json:"id" gorm:"primaryKey;size:32"`
	ProductID     string  `json:"product_id" gorm:"size:32"`
	Amount        float64 `json:"amount"`
	Duration      uint    `json:"duration"`
	TrialDuration uint    `json:"trial_duration"` //0 if no trial
	gorm.Model
}
