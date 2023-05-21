package purchase

import (
	"subscriptions/domains/plans"
	"subscriptions/domains/products"
	"subscriptions/domains/vouchers"

	"gorm.io/gorm"
)

// PurchaseService contains the sales logic handling
type PurchaseService struct {
	db             *gorm.DB
	productService products.Manager
	planService    plans.Manager
	voucherService vouchers.Manager
}

// this service takes the purchase DTO and makes it into a subscription
