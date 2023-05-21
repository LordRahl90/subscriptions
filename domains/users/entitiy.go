package users

import (
	"gorm.io/gorm"
)

type UserType string

const (
	// UserTypeAdmin admin user
	UserTypeAdmin UserType = "admin"

	// UserTypeRegular regular user or customer
	UserTypeRegular = "regular"
)

// User contains the base user structure
type User struct {
	ID       string   `json:"id" gorm:"primaryKey;size:100"`
	Email    string   `json:"email" gorm:"uniqueIndex;size:100"`
	UserType UserType `json:"user_type"`
	Name     string   `json:"name"`
	Password string   `json:"password"`
	gorm.Model
}
