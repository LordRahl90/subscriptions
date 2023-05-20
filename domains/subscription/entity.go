package subscription

import (
	"gorm.io/gorm"
)

// UpdateAction type of action that can happen on subscription updates
type UpdateAction int

const (
	// UpdateActionUnpaused action to pause subscription
	UpdateActionPaused UpdateAction = iota + 1
	// UpdateActionUnpaused action to unpause subscription
	UpdateActionUnpaused
)

// Subscription contains the subscription details
type Subscription struct {
	ID                 string  `json:"id" gorm:"primaryKey;size:32"`
	UserID             string  `json:"user_id" gorm:"primaryKey;size:32"`
	ProductID          string  `json:"product_id" gorm:"primaryKey;size:32"`
	SubscriptionPlanID string  `json:"subscription_plan_id" gorm:"primaryKey;size:32"`
	Duration           int     `json:"duration"`
	Amount             float64 `json:"amount"`
	Discount           float64 `json:"discount"`
	Tax                float64 `json:"tax"`
	// we dont charge for trial, given there's a month max anyways
	Trial  bool         `json:"trial"`
	Status UpdateAction `json:"status"`
	gorm.Model
}
