package subscription

import (
	"time"

	"gorm.io/gorm"
)

// Status subscription status
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

	month = 24 * 60 * 30 * time.Minute
)

// String returns the string representation of status
func (s Status) String() string {
	switch s {
	case StatusCancelled:
		return "cancelled"
	case StatusActive:
		return "active"
	case StatusPaused:
		return "paused"
	case StatusEnded:
		return "ended"
	default:
		return "unknown"
	}
}

// FromString return a status representation from a given string
func FromString(s string) Status {
	switch s {
	case StatusCancelled.String():
		return StatusCancelled
	case StatusActive.String():
		return StatusActive
	case StatusPaused.String():
		return StatusPaused
	case StatusEnded.String():
		return StatusEnded
	default:
		return -1
	}
}

// Subscription contains the subscription details
type Subscription struct {
	ID                 string  `json:"id" gorm:"primaryKey;size:32"`
	UserID             string  `json:"user_id" gorm:"primaryKey;size:100"`
	ProductID          string  `json:"product_id" gorm:"primaryKey;size:100"`
	SubscriptionPlanID string  `json:"subscription_plan_id" gorm:"primaryKey;size:100"`
	VoucherID          string  `json:"voucher_id" gorm:"primaryKey;size:100"`
	Duration           uint    `json:"duration"`
	Amount             float64 `json:"amount" gorm:"type:double(30,2)"`
	Discount           float64 `json:"discount" gorm:"type:double(30,2)"`
	Tax                float64 `json:"tax" gorm:"type:double(30,2)"`
	Total              float64 `json:"total" gorm:"type:double(30,2)"`
	TrialPeriod        uint    `json:"trial_period"`
	Status             Status  `json:"status"`
	gorm.Model
}

// EndDate returns the expected end date of the subscription
func (s Subscription) EndDate() time.Time {
	if s.CreatedAt.IsZero() {
		return time.Now().Add(-1 * 24 * 60 * time.Minute)
	}

	duration := s.Duration + s.TrialPeriod

	return s.CreatedAt.Add(time.Duration(duration * uint(month)))
}

// TrialEnds computes date that this subscription ends
func (s Subscription) TrialEnds() time.Time {
	return s.CreatedAt.Add(time.Duration(s.TrialPeriod * uint(month)))
}

// IsTrial determines if the subscription is still on trial or not
func (s Subscription) IsTrial() bool {
	if s.CreatedAt.IsZero() || s.TrialPeriod == 0 {
		return false
	}

	return s.TrialEnds().After(time.Now())
}
