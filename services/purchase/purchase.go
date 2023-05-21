package purchase

import (
	"context"
	"fmt"
	"subscriptions/domains/plans"
	"subscriptions/domains/products"
	"subscriptions/domains/subscription"
	"subscriptions/domains/vouchers"
	"time"
)

// Service contains the sales logic handling
type Service struct {
	productService      products.Manager
	planService         plans.Manager
	voucherService      vouchers.Manager
	subscriptionService subscription.Manager
}

// New initializes a new purchase service
func New(ps products.Manager, pls plans.Manager, vs vouchers.Manager, sbs subscription.Manager) *Service {
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
		v, err := ps.voucherService.FindByCode(ctx, p.Voucher)
		if err != nil {
			return nil, err
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

	// All new subscriptions are active
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
		Status:             subscription.StatusActive,
	}
	if err := ps.subscriptionService.Create(ctx, sub); err != nil {
		return nil, err
	}

	return sub, nil
}
