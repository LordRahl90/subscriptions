package purchase

import (
	"context"
	"os"
	"testing"
	"time"

	"subscriptions/domains/plans"
	"subscriptions/domains/products"
	"subscriptions/domains/subscription"
	"subscriptions/domains/vouchers"
	"subscriptions/services/purchase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		os.Exit(code)
	}()
	code = m.Run()
}

func TestProcessWithoutDiscount(t *testing.T) {
	ctx := context.Background()
	p := &Purchase{
		ProductID: primitive.NewObjectID().Hex(),
		PlanID:    primitive.NewObjectID().Hex(),
	}
	ps := &mocks.ProductMock{
		FindOneFunc: func(ctx context.Context, id string) (*products.Product, error) {
			return &products.Product{
				TaxRate: 20,
			}, nil
		},
	}
	pls := &mocks.SubscritionPlanMock{
		FindOneFunc: func(ctx context.Context, id string) (*plans.SubscriptionPlan, error) {
			return &plans.SubscriptionPlan{
				Amount:        2000,
				Duration:      3,
				TrialDuration: 1,
			}, nil
		},
	}

	vs := &mocks.VoucherMocks{
		FindByCodeFunc: func(ctx context.Context, productID, code string) (*vouchers.Voucher, error) {
			return &vouchers.Voucher{
				ID:          primitive.NewObjectID().Hex(),
				VoucherType: vouchers.VoucherTypePercentage,
				Percentage:  10,
				ExpiresOn:   time.Now().Add(24 * time.Hour),
			}, nil
		},
	}

	sbs := &mocks.SubscriptionMock{
		CreateFunc: func(ctx context.Context, s *subscription.Subscription) error {
			s.ID = primitive.NewObjectID().Hex()
			return nil
		},
	}

	pps := New(ps, pls, vs, sbs)
	sub, err := pps.Process(ctx, p)
	require.NoError(t, err)
	require.NotEmpty(t, sub)

	assert.NotEmpty(t, sub.ID)
	assert.Equal(t, 0.0, sub.Discount)
	assert.Equal(t, 400.0, sub.Tax)
	assert.Equal(t, 2400.0, sub.Total)
}

func TestProcessWithDiscount(t *testing.T) {
	ctx := context.Background()
	p := &Purchase{
		ProductID: primitive.NewObjectID().Hex(),
		PlanID:    primitive.NewObjectID().Hex(),
		Voucher:   "HELLO_ONE_TWO",
	}
	ps := &mocks.ProductMock{
		FindOneFunc: func(ctx context.Context, id string) (*products.Product, error) {
			return &products.Product{
				TaxRate: 20,
			}, nil
		},
	}
	pls := &mocks.SubscritionPlanMock{
		FindOneFunc: func(ctx context.Context, id string) (*plans.SubscriptionPlan, error) {
			return &plans.SubscriptionPlan{
				Amount:        2000,
				Duration:      3,
				TrialDuration: 1,
			}, nil
		},
	}

	vs := &mocks.VoucherMocks{
		FindByCodeFunc: func(ctx context.Context, productID, code string) (*vouchers.Voucher, error) {
			return &vouchers.Voucher{
				ID:          primitive.NewObjectID().Hex(),
				ProductID:   p.ProductID,
				VoucherType: vouchers.VoucherTypePercentage,
				Percentage:  10,
				ExpiresOn:   time.Now().Add(24 * time.Hour),
			}, nil
		},
	}

	sbs := &mocks.SubscriptionMock{
		CreateFunc: func(ctx context.Context, s *subscription.Subscription) error {
			s.ID = primitive.NewObjectID().Hex()
			return nil
		},
	}

	pps := New(ps, pls, vs, sbs)
	sub, err := pps.Process(ctx, p)
	require.NoError(t, err)
	require.NotEmpty(t, sub)

	assert.NotEmpty(t, sub.ID)
	assert.Equal(t, 200.0, sub.Discount)
	assert.Equal(t, 360.0, sub.Tax)
	assert.Equal(t, 2160.0, sub.Total)
}

func TestProcessWithExpiredDiscount(t *testing.T) {
	ctx := context.Background()
	p := &Purchase{
		ProductID: primitive.NewObjectID().Hex(),
		PlanID:    primitive.NewObjectID().Hex(),
		Voucher:   "HELLO_ONE_TWO",
	}
	ps := &mocks.ProductMock{
		FindOneFunc: func(ctx context.Context, id string) (*products.Product, error) {
			return &products.Product{
				TaxRate: 20,
			}, nil
		},
	}
	pls := &mocks.SubscritionPlanMock{
		FindOneFunc: func(ctx context.Context, id string) (*plans.SubscriptionPlan, error) {
			return &plans.SubscriptionPlan{
				Amount:        2000,
				Duration:      3,
				TrialDuration: 1,
			}, nil
		},
	}

	vs := &mocks.VoucherMocks{
		FindByCodeFunc: func(ctx context.Context, productID, code string) (*vouchers.Voucher, error) {
			return &vouchers.Voucher{
				ID:          primitive.NewObjectID().Hex(),
				ProductID:   p.ProductID,
				VoucherType: vouchers.VoucherTypePercentage,
				Percentage:  10,
				ExpiresOn:   time.Now().Add(time.Duration(-48 * time.Hour)),
			}, nil
		},
	}

	sbs := &mocks.SubscriptionMock{
		CreateFunc: func(ctx context.Context, s *subscription.Subscription) error {
			s.ID = primitive.NewObjectID().Hex()
			return nil
		},
	}

	pps := New(ps, pls, vs, sbs)
	sub, err := pps.Process(ctx, p)
	require.EqualError(t, err, "voucher has expired")
	require.Nil(t, sub)
}
