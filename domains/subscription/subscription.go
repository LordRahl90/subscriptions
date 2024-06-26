package subscription

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

// Service implements the Manager interface
type Service struct {
	db *gorm.DB
}

// New returns a new instance of subscription service
func New(db *gorm.DB) (*Service, error) {
	if err := db.AutoMigrate(&Subscription{}); err != nil {
		return nil, err
	}

	return &Service{
		db: db,
	}, nil
}

// Create creates a new subscription.
func (ss *Service) Create(ctx context.Context, p *Subscription) error {
	p.ID = primitive.NewObjectID().Hex()
	p.CreatedAt = time.Now()

	return ss.db.WithContext(ctx).Save(&p).Error
}

// Find returns all the subscriptions available to a user
func (ss *Service) Find(ctx context.Context, userID string) (result []Subscription, err error) {
	err = ss.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&result).Error
	return
}

// FindOne finds a subscription with its ID
func (ss *Service) FindOne(ctx context.Context, id string) (*Subscription, error) {
	var result *Subscription
	err := ss.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	return result, err
}

// FindByStatus finds a user's subscription based on the status
func (ss *Service) FindByStatus(ctx context.Context, userID string, status Status) (result []Subscription, err error) {
	err = ss.db.
		WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, status).
		Find(&result).Error
	return
}

// UpdateStatus updates the subscritpion status
func (ss *Service) UpdateStatus(ctx context.Context, subID string, status Status) error {
	return ss.db.WithContext(ctx).
		Model(&Subscription{}).
		Where("id = ?", subID).
		Update("status", status).Error
}
