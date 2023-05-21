package subscription

import (
	"gorm.io/gorm"
)

// UpdateAction type of action that can happen on subscription updates
type Status int

const (
	// StatusCancelled cancelled subscription
	StatusCancelled Status = iota
	// StatusActive active subscription
	StatusActive
	// StatusPaused paused subscription
	StatusPaused
	// StatusEnded subscriptions that have been completed.
	// i.e the created_at and duration is after the current date
	StatusEnded
)

// Subscription contains the subscription details
type Subscription struct {
	ID                 string  `json:"id" gorm:"primaryKey;size:32"`
	UserID             string  `json:"user_id" gorm:"primaryKey;size:32"`
	ProductID          string  `json:"product_id" gorm:"primaryKey;size:32"`
	SubscriptionPlanID string  `json:"subscription_plan_id" gorm:"primaryKey;size:32"`
	VoucherID          string  `json:"voucher_id" gorm:"primaryKey;size:32"`
	Duration           int     `json:"duration"`
	Amount             float64 `json:"amount"`
	Discount           float64 `json:"discount"`
	Tax                float64 `json:"tax"`
	TrialPeriod        int     `json:"trial_period"`
	Status             Status  `json:"status"`
	gorm.Model
}

/*
plan has amount, duration, trials
voucher has the discount values
product has tax
*/
