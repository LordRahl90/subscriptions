package products

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

var _ Manager = (*ProductService)(nil)

// ProductService implements the Manager interface
type ProductService struct {
	db *gorm.DB
	m  *sync.Mutex
}

// New returns a new instance of product service
func New(db *gorm.DB) (*ProductService, error) {
	if err := db.AutoMigrate(&Product{}); err != nil {
		return nil, err
	}
	return &ProductService{
		db: db,
		m:  &sync.Mutex{},
	}, nil
}

// Create creates a new product
func (ps *ProductService) Create(ctx context.Context, p *Product) error {
	p.ID = primitive.NewObjectID().Hex()
	p.CreatedAt = time.Now()

	return ps.db.WithContext(ctx).Save(&p).Error
}

// Find returns all the products within the system
func (ps *ProductService) Find(ctx context.Context) (result []Product, err error) {
	err = ps.db.WithContext(ctx).Find(&result).Error
	return
}

// FindOne finds a product with its ID
func (ps *ProductService) FindOne(ctx context.Context, id string) (*Product, error) {
	var result *Product
	err := ps.db.WithContext(ctx).Where("id = ?", id).First(&result).Error
	return result, err
}
