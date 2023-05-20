package vouchers

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

var _ Manager = (*VoucherService)(nil)

// VoucherService implements the Manager interface
type VoucherService struct {
	db *gorm.DB
	m  *sync.Mutex
}

// New returns a new instance of Voucher service
func New(db *gorm.DB) (*VoucherService, error) {
	if err := db.AutoMigrate(&Voucher{}); err != nil {
		return nil, err
	}
	return &VoucherService{
		db: db,
		m:  &sync.Mutex{},
	}, nil
}

// Create creates a new Voucher
func (vs *VoucherService) Create(ctx context.Context, v *Voucher) error {
	if err := v.validate(); err != nil {
		return err
	}
	v.ID = primitive.NewObjectID().Hex()
	v.CreatedAt = time.Now()

	return vs.db.WithContext(ctx).Save(&v).Error
}

// FindOne finds a voucher with its ID
func (ps *VoucherService) FindOne(ctx context.Context, id string) (*Voucher, error) {
	var result *Voucher
	err := ps.db.WithContext(ctx).
		Preload("Products").
		Where("id = ?", id).
		First(&result).Error
	return result, err
}
