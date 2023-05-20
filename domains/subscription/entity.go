package subscription

import (
	"time"

	"gorm.io/gorm"
)

// Subscription contains the subscription details
type Subscription struct {
	ID                 string    `json:"id" gorm:"primaryKey;size:32"`
	UserID             string    `json:"user_id" gorm:"primaryKey;size:32"`
	ProductID          string    `json:"product_id" gorm:"primaryKey;size:32"`
	SubscriptionPlanID string    `json:"subscription_plan_id" gorm:"primaryKey;size:32"`
	Duration           int       `json:"duration"`
	Amount             float64   `json:"amount"`
	Discount           float64   `json:"discount"`
	Cancelled          bool      `json:"cancelled"`
	Tax                float64   `json:"tax"`
	TrialEnds          time.Time `json:"trial_ends"`
	gorm.Model
}
