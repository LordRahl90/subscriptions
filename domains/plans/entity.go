package plans

import "gorm.io/gorm"

// SubscriptionPlan contains the details of a subscription plan
type SubscriptionPlan struct {
	ID        string `json:"id" gorm:"primaryKey;size:32"`
	ProductID string `json:"product_id" gorm:"size:32"`
	// Amount cost/price of this plan
	Amount float64 `json:"amount" gorm:"type:double(30,2)"`
	// Duration this is duration of the subscription in months
	Duration uint `json:"duration"`
	// TrialDuration duration of the trial 0 if there's none
	TrialDuration uint `json:"trial_duration"`
	gorm.Model
}
