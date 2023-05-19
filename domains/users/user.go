package users

import (
	"context"
	"time"

	"subscriptions/domains/core"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

// UserService contains the service for manipulating the user entity
type UserService struct {
	db *gorm.DB
}

// NewUserService returns a new instance of user service
func New(db *gorm.DB) (*UserService, error) {
	if err := db.AutoMigrate(&User{}); err != nil {
		return nil, err
	}
	return &UserService{db: db}, nil
}

// Create creates a new user record
func (us *UserService) Create(ctx context.Context, u *User) error {
	u.ID = primitive.NewObjectID().Hex()
	u.CreatedAt = time.Now()
	hash, err := core.GeneratePassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hash
	return us.db.Save(&u).Error
}

// Find finds a user with the given ID
func (us *UserService) Find(ctx context.Context, id string) (*User, error) {
	var user *User
	err := us.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return user, err
}

// FindByEmail finds a user by the given email
func (us *UserService) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user *User
	err := us.db.Where("email = ?", email).First(&user).Error
	return user, err
}
