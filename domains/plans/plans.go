package plans

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

var _ Manager = (*SubscriptionPlanService)(nil)

// SubscriptionPlanService implements the Manager interface
type SubscriptionPlanService struct {
	db *gorm.DB
}

// New returns a new instance of SubscriptionPlan service
func New(db *gorm.DB) (*SubscriptionPlanService, error) {
	if err := db.AutoMigrate(&SubscriptionPlan{}); err != nil {
		return nil, err
	}
	return &SubscriptionPlanService{
		db: db,
	}, nil
}

// Create creates a new SubscriptionPlan
func (vs *SubscriptionPlanService) Create(ctx context.Context, v *SubscriptionPlan) error {
	v.ID = primitive.NewObjectID().Hex()
	v.CreatedAt = time.Now()

	return vs.db.WithContext(ctx).Save(&v).Error
}

// Find returns all the subscription plans for a given product
func (ps *SubscriptionPlanService) Find(ctx context.Context, productID string) (result []SubscriptionPlan, err error) {
	err = ps.db.WithContext(ctx).Where("product_id = ?", productID).Find(&result).Error
	return
}

// FindOne finds a subscriptionplan with its ID
func (ps *SubscriptionPlanService) FindOne(ctx context.Context, id string) (*SubscriptionPlan, error) {
	var result *SubscriptionPlan
	err := ps.db.WithContext(ctx).
		Where("id = ?", id).
		First(&result).Error
	return result, err
}
