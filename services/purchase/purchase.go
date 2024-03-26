package purchase

import (
	"context"
	"fmt"
	"time"

	"subscriptions/domains/subscription"
	"subscriptions/domains/vouchers"
)

// Service contains the sales logic handling
type Service struct {
	productService      ProductManager
	planService         PlanManager
	voucherService      VoucherManager
	subscriptionService SubscriptionManager
}

// New initializes a new purchase service
func New(
	ps ProductManager,
	pls PlanManager,
	vs VoucherManager,
	sbs SubscriptionManager,
) *Service {
	return &Service{
		productService:      ps,
		planService:         pls,
		voucherService:      vs,
		subscriptionService: sbs,
	}
}

// Process takes the purchase DTO and uses it to create a new subscription
func (ps Service) Process(ctx context.Context, p *Purchase) (*subscription.Subscription, error) {
	var voucher vouchers.Voucher
	product, err := ps.productService.FindOne(ctx, p.ProductID)
	if err != nil {
		return nil, err
	}

	plan, err := ps.planService.FindOne(ctx, p.PlanID)
	if err != nil {
		return nil, err
	}

	if p.Voucher != "" {
		v, err := ps.voucherService.FindByCode(ctx, p.ProductID, p.Voucher)
		if err != nil {
			return nil, err
		}

		if v.ProductID != p.ProductID {
			return nil, fmt.Errorf("this voucher is not for this product")
		}

		if v.ExpiresOn.Before(time.Now()) {
			return nil, fmt.Errorf("voucher has expired")
		}

		voucher = *v
	}

	taxRate := product.TaxRate
	amount := plan.Amount
	discount := voucher.Calculate(amount)
	subTotal := amount - discount
	tax := taxRate / 100 * subTotal
	total := subTotal + tax

	sub := &subscription.Subscription{
		UserID:             p.UserID,
		ProductID:          p.ProductID,
		SubscriptionPlanID: p.PlanID,
		VoucherID:          voucher.ID,
		Duration:           plan.Duration,
		Amount:             amount,
		Discount:           discount,
		Tax:                tax,
		Total:              total,
		TrialPeriod:        plan.TrialDuration,
		// All new subscriptions are active
		Status: subscription.StatusActive,
	}
	if err := ps.subscriptionService.Create(ctx, sub); err != nil {
		return nil, err
	}

	return sub, nil
}
