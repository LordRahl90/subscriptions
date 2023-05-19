package responses

import (
	"time"
)

// User user response struct
type User struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Gender    string     `json:"gender"`
	Age       uint16     `json:"age"`
	Token     string     `json:"token"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
