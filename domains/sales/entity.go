package sales

// PurchaseRequest request for a subscription purchase
type PurchaseRequest struct {
	ID                 string `json:"id" gorm:"primaryKey;size:32"`
	ProductID          string `json:"product_id" gorm:"size:32"`
	SubscriptionPlanID string `json:"subscription_plan_id" gorm:"size:32"`
	VoucherID          string `json:"voucher_id" gorm:"size:32"`
}
