package users

import (
	"gorm.io/gorm"
)

// User contains the base user structure
type User struct {
	ID       string `json:"id" gorm:"primaryKey;size:32"`
	Email    string `json:"email" gorm:"uniqueIndex;size:100"`
	Name     string `json:"name"`
	Password string `json:"password"`
	gorm.Model
}
