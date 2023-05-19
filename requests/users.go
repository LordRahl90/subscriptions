package requests

// CreateUser request format for creating user
type CreateUser struct {
	Name     string `json:"name"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	UserType string `json:"user_type"`
}

// Authenticate requests format for authenticating user
type Authenticate struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
