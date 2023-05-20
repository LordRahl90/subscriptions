package subscription

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

var _ Manager = (*SubscriptionService)(nil)

// SubscriptionService implements the Manager interface
type SubscriptionService struct {
	db *gorm.DB
	m  *sync.Mutex
}

// New returns a new instance of subscription service
func New(db *gorm.DB) (*SubscriptionService, error) {
	if err := db.AutoMigrate(&Subscription{}); err != nil {
		return nil, err
	}
	return &SubscriptionService{
		db: db,
		m:  &sync.Mutex{},
	}, nil
}

// Create creates a new subscription
func (ss *SubscriptionService) Create(ctx context.Context, p *Subscription) error {
	p.ID = primitive.NewObjectID().Hex()
	p.CreatedAt = time.Now()

	return ss.db.WithContext(ctx).Save(&p).Error
}

// Find returns all the subscriptions available to a user
// we might need a flag here to only return active subscriptions
func (ss *SubscriptionService) Find(ctx context.Context, userID string) (result []Subscription, err error) {
	err = ss.db.WithContext(ctx).Find(&result).Error
	return
}

// FindOne finds a subscription with its ID
func (ss *SubscriptionService) FindOne(ctx context.Context, id string) (*Subscription, error) {
	var result *Subscription
	err := ss.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	return result, err
}

// Update implements Manager
func (*SubscriptionService) Update(ctx context.Context, subID string) {
	panic("unimplemented")
}
