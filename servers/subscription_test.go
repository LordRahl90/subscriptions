package servers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"subscriptions/domains/core"
	"subscriptions/domains/plans"
	"subscriptions/domains/products"
	"subscriptions/domains/subscription"
	"subscriptions/domains/users"
	"subscriptions/domains/vouchers"
	"subscriptions/requests"
	"subscriptions/responses"
	"subscriptions/services/purchase/mocks"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

const month = 30 * 24

func TestCreateSubscriptionWithDiscount(t *testing.T) {
	mockDB := make(map[string][]subscription.Subscription)
	req := requests.Purchase{
		ProductID: primitive.NewObjectID().Hex(),
		PlanID:    primitive.NewObjectID().Hex(),
		Voucher:   "HELLO_WORLD",
	}

	um := &userMock{}
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
		FindByCodeFunc: func(ctx context.Context, id string) (*vouchers.Voucher, error) {
			return &vouchers.Voucher{
				ID:          primitive.NewObjectID().Hex(),
				VoucherType: vouchers.VoucherTypePercentage,
				Percentage:  20,
			}, nil
		},
	}
	sbs := &mocks.SubscriptionMock{
		CreateFunc: func(ctx context.Context, s *subscription.Subscription) error {
			s.ID = primitive.NewObjectID().Hex()
			s.CreatedAt = time.Now()
			mockDB[s.UserID] = append(mockDB[s.UserID], *s)
			return nil
		},
		FindFunc: func(ctx context.Context, userID string) ([]subscription.Subscription, error) {
			return mockDB[userID], nil
		},
	}

	b, err := json.Marshal(req)
	require.NoError(t, err)
	require.NotEmpty(t, b)

	token := createToken(t)

	svr := New(um, ps, vs, pls, sbs)
	res := requestHelperWithMockedServer(t, svr, http.MethodPost, "/subscriptions", token, b)
	require.Equal(t, http.StatusCreated, res.Code)

	var response responses.Subscription
	err = json.Unmarshal(res.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.NotEmpty(t, response.ID)
	assert.NotEmpty(t, response.EndDate)
	assert.True(t, response.EndDate.After(response.StartDate))
	assert.Equal(t, uint(3), response.Duration)
	assert.Equal(t, uint(1), response.TrialDuration)
	assert.Equal(t, 320.0, response.Tax)
	assert.Equal(t, 400.0, response.Discount)
	assert.Equal(t, 1920.0, response.Total)
	assert.Equal(t, "active", response.Status)

	res = requestHelperWithMockedServer(t, svr, http.MethodGet, "/subscriptions", token, nil)
	require.Equal(t, http.StatusOK, res.Code)

	var results []responses.Subscription
	err = json.Unmarshal(res.Body.Bytes(), &results)
	require.NoError(t, err)
	require.Len(t, results, 1)

	assert.Equal(t, response.ID, results[0].ID)
	assert.Equal(t, response.Total, results[0].Total)
	assert.Equal(t, response.Tax, results[0].Tax)
	assert.Equal(t, response.Status, results[0].Status)
}

func TestCreateSubscriptionWithOutDiscount(t *testing.T) {
	mockDB := make(map[string][]*subscription.Subscription)
	req := requests.Purchase{
		ProductID: primitive.NewObjectID().Hex(),
		PlanID:    primitive.NewObjectID().Hex(),
	}

	um := &userMock{}
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
		FindByCodeFunc: func(ctx context.Context, id string) (*vouchers.Voucher, error) {
			return &vouchers.Voucher{
				ID:          primitive.NewObjectID().Hex(),
				VoucherType: vouchers.VoucherTypePercentage,
				Percentage:  20,
			}, nil
		},
	}
	sbs := &mocks.SubscriptionMock{
		CreateFunc: func(ctx context.Context, s *subscription.Subscription) error {
			s.ID = primitive.NewObjectID().Hex()
			s.CreatedAt = time.Now()
			mockDB[s.UserID] = append(mockDB[s.UserID], s)
			return nil
		},
	}

	b, err := json.Marshal(req)
	require.NoError(t, err)
	require.NotEmpty(t, b)

	token := createToken(t)

	svr := New(um, ps, vs, pls, sbs)
	res := requestHelperWithMockedServer(t, svr, http.MethodPost, "/subscriptions", token, b)
	require.Equal(t, http.StatusCreated, res.Code)

	var response responses.Subscription
	err = json.Unmarshal(res.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.NotEmpty(t, response.ID)
	assert.NotEmpty(t, response.EndDate)
	assert.True(t, response.EndDate.After(response.StartDate))
	assert.Equal(t, "active", response.Status)
	assert.Equal(t, uint(3), response.Duration)
	assert.Equal(t, uint(1), response.TrialDuration)
	assert.Equal(t, 400.0, response.Tax)
	assert.Equal(t, 0.0, response.Discount)
	assert.Equal(t, 2400.0, response.Total)
}

func TestGetUserSubscriptions(t *testing.T) {
	td := core.TokenData{
		UserID:   uuid.NewString(),
		Email:    gofakeit.Email(),
		UserType: string(users.UserTypeAdmin),
	}
	token, err := td.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, token)

	sbs := &mocks.SubscriptionMock{
		FindFunc: func(ctx context.Context, userID string) ([]subscription.Subscription, error) {
			return []subscription.Subscription{
				{
					ID:                 primitive.NewObjectID().Hex(),
					UserID:             td.UserID,
					ProductID:          primitive.NewObjectID().Hex(),
					SubscriptionPlanID: primitive.NewObjectID().Hex(),
					Duration:           3,
					Amount:             2000,
					Discount:           500,
					Tax:                400,
					Total:              2400,
					Model: gorm.Model{
						CreatedAt: time.Now(),
					},
				},
				{
					ID:                 primitive.NewObjectID().Hex(),
					UserID:             td.UserID,
					ProductID:          primitive.NewObjectID().Hex(),
					SubscriptionPlanID: primitive.NewObjectID().Hex(),
					Duration:           5,
					TrialPeriod:        1,
					Amount:             2500,
					Tax:                450,
					Total:              2950,
					Model: gorm.Model{
						CreatedAt: time.Now(),
					},
				},
			}, nil
		},
	}

	svr := New(nil, nil, nil, nil, sbs)
	res := requestHelperWithMockedServer(t, svr, http.MethodGet, "/subscriptions", token, nil)

	require.Equal(t, http.StatusOK, res.Code)

	var results []responses.Subscription
	err = json.Unmarshal(res.Body.Bytes(), &results)
	require.NoError(t, err)
	require.Len(t, results, 2)

	assert.NotEmpty(t, results[0].EndDate)
	assert.True(t, results[0].EndDate.After(results[0].StartDate))
	diff := results[0].EndDate.Sub(results[0].StartDate)
	assert.Equal(t, 3.0, diff.Hours()/month)
	assert.Equal(t, 2400.0, results[0].Total)
	assert.Equal(t, 2000.0, results[0].Price)
	assert.Equal(t, 400.0, results[0].Tax)
	assert.Equal(t, 500.0, results[0].Discount)
	assert.Equal(t, uint(3), results[0].Duration)
	assert.Equal(t, uint(0), results[0].TrialDuration)
	assert.Equal(t, uint(3), results[0].TotalDuration)

	assert.Equal(t, 2950.0, results[1].Total)
	assert.NotEmpty(t, results[1].EndDate)
	assert.True(t, results[1].EndDate.After(results[1].StartDate))
	diff = results[1].EndDate.Sub(results[1].StartDate)
	assert.Equal(t, 6.0, diff.Hours()/month)
	assert.Equal(t, 2950.0, results[1].Total)
	assert.Equal(t, 2500.0, results[1].Price)
	assert.Equal(t, 450.0, results[1].Tax)
	assert.Equal(t, 0.0, results[1].Discount)
	assert.Equal(t, uint(5), results[1].Duration)
	assert.Equal(t, uint(1), results[1].TrialDuration)
	assert.Equal(t, uint(6), results[1].TotalDuration)

}

func TestGetUserSubscriptions_WithError(t *testing.T) {
	td := core.TokenData{
		UserID:   uuid.NewString(),
		Email:    gofakeit.Email(),
		UserType: string(users.UserTypeAdmin),
	}
	token, err := td.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, token)

	sbs := &mocks.SubscriptionMock{
		FindFunc: func(ctx context.Context, userID string) ([]subscription.Subscription, error) {
			return nil, fmt.Errorf("cannot connect to db")
		},
	}

	svr := New(nil, nil, nil, nil, sbs)
	res := requestHelperWithMockedServer(t, svr, http.MethodGet, "/subscriptions", token, nil)

	require.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestGetSubscriptionDetails(t *testing.T) {
	td := core.TokenData{
		UserID: primitive.NewObjectID().Hex(),
		Email:  gofakeit.Email(),
	}
	token, err := td.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, token)

	sbs := &mocks.SubscriptionMock{
		FindOneFunc: func(ctx context.Context, subID string) (*subscription.Subscription, error) {
			return &subscription.Subscription{
				ID:                 subID,
				UserID:             td.UserID,
				ProductID:          primitive.NewObjectID().Hex(),
				SubscriptionPlanID: primitive.NewObjectID().Hex(),
				Duration:           3,
				Amount:             2000,
				Discount:           500,
				Tax:                400,
				Total:              2400,
				Status:             subscription.StatusActive,
				Model: gorm.Model{
					CreatedAt: time.Now(),
				},
			}, nil
		},
	}

	id := primitive.NewObjectID().Hex()
	svr := New(nil, nil, nil, nil, sbs)
	res := requestHelperWithMockedServer(t, svr, http.MethodGet, "/subscriptions/"+id, token, nil)

	require.Equal(t, http.StatusOK, res.Code)

	var result responses.Subscription
	err = json.Unmarshal(res.Body.Bytes(), &result)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	assert.NotEmpty(t, result.EndDate)
	assert.True(t, result.EndDate.After(result.StartDate))
	diff := result.EndDate.Sub(result.StartDate)
	assert.Equal(t, 3.0, diff.Hours()/month)
	assert.Equal(t, 2400.0, result.Total)
	assert.Equal(t, 2000.0, result.Price)
	assert.Equal(t, 400.0, result.Tax)
	assert.Equal(t, 500.0, result.Discount)
	assert.Equal(t, uint(3), result.Duration)
	assert.Equal(t, uint(0), result.TrialDuration)
	assert.Equal(t, uint(3), result.TotalDuration)
	assert.Equal(t, subscription.StatusActive.String(), result.Status)
}

func TestGetSubscriptionDetails_WithError(t *testing.T) {
	td := core.TokenData{
		UserID: primitive.NewObjectID().Hex(),
		Email:  gofakeit.Email(),
	}
	token, err := td.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, token)

	sbs := &mocks.SubscriptionMock{
		FindOneFunc: func(ctx context.Context, subID string) (*subscription.Subscription, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	id := primitive.NewObjectID().Hex()
	svr := New(nil, nil, nil, nil, sbs)
	res := requestHelperWithMockedServer(t, svr, http.MethodGet, "/subscriptions/"+id, token, nil)

	require.Equal(t, http.StatusNotFound, res.Code)
	exp := `{"error":"not found","success":false}`
	assert.Equal(t, exp, res.Body.String())
}

func TestGetSubscriptionDetails_WithRandomToken(t *testing.T) {
	td := core.TokenData{
		UserID: primitive.NewObjectID().Hex(),
		Email:  gofakeit.Email(),
	}
	token, err := td.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, token)

	sbs := &mocks.SubscriptionMock{
		FindOneFunc: func(ctx context.Context, subID string) (*subscription.Subscription, error) {
			return &subscription.Subscription{
				ID:                 subID,
				UserID:             td.UserID,
				ProductID:          primitive.NewObjectID().Hex(),
				SubscriptionPlanID: primitive.NewObjectID().Hex(),
				Duration:           3,
				Amount:             2000,
				Discount:           500,
				Tax:                400,
				Total:              2400,
				Status:             subscription.StatusActive,
				Model: gorm.Model{
					CreatedAt: time.Now(),
				},
			}, nil
		},
	}

	id := primitive.NewObjectID().Hex()
	svr := New(nil, nil, nil, nil, sbs)
	res := requestHelperWithMockedServer(t, svr, http.MethodGet, "/subscriptions/"+id, createToken(t), nil)

	require.Equal(t, http.StatusUnauthorized, res.Code)
	exp := `{"error":"you are not allowed to access this subscription details","success":false}`
	assert.Equal(t, exp, res.Body.String())
}

func TestPauseSubscriptionDuringTrialPeriod(t *testing.T) {
	td := core.TokenData{
		UserID:   uuid.NewString(),
		Email:    gofakeit.Email(),
		UserType: string(users.UserTypeAdmin),
	}
	token, err := td.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, token)

	id := primitive.NewObjectID().Hex()

	sbs := &mocks.SubscriptionMock{
		FindOneFunc: func(ctx context.Context, subID string) (*subscription.Subscription, error) {
			return &subscription.Subscription{
				ID:                 id,
				UserID:             td.UserID,
				ProductID:          primitive.NewObjectID().Hex(),
				SubscriptionPlanID: primitive.NewObjectID().Hex(),
				Duration:           5,
				TrialPeriod:        1,
				Amount:             2500,
				Tax:                450,
				Total:              2950,
				Model: gorm.Model{
					CreatedAt: time.Now(),
				},
			}, nil
		},
	}
	svr := New(nil, nil, nil, nil, sbs)
	res := requestHelperWithMockedServer(t, svr, http.MethodPatch, "/subscriptions/"+id+"/pause", token, nil)
	require.Equal(t, http.StatusBadRequest, res.Code)
	exp := `{"error":"you cannot pause during trial period","success":false}`
	assert.Equal(t, exp, res.Body.String())
}

func TestPauseSubscription_FailedUpdate(t *testing.T) {
	td := core.TokenData{
		UserID:   uuid.NewString(),
		Email:    gofakeit.Email(),
		UserType: string(users.UserTypeAdmin),
	}
	token, err := td.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, token)

	id := primitive.NewObjectID().Hex()

	sbs := &mocks.SubscriptionMock{
		FindOneFunc: func(ctx context.Context, subID string) (*subscription.Subscription, error) {
			return &subscription.Subscription{
				ID:                 id,
				UserID:             td.UserID,
				ProductID:          primitive.NewObjectID().Hex(),
				SubscriptionPlanID: primitive.NewObjectID().Hex(),
				Duration:           5,
				Amount:             2500,
				Tax:                450,
				Total:              2950,
				Model: gorm.Model{
					CreatedAt: time.Now(),
				},
			}, nil
		},
		UpdateStatusFunc: func(ctx context.Context, subID string, status subscription.Status) error {
			return fmt.Errorf("cannot update successfully")
		},
	}

	svr := New(nil, nil, nil, nil, sbs)
	res := requestHelperWithMockedServer(t, svr, http.MethodPatch, "/subscriptions/"+id+"/pause", token, nil)
	require.Equal(t, http.StatusInternalServerError, res.Code)
	exp := `{"error":"cannot update successfully","success":false}`
	assert.Equal(t, exp, res.Body.String())
}

func TestPauseSubscription_RandomToken(t *testing.T) {
	td := core.TokenData{
		UserID:   uuid.NewString(),
		Email:    gofakeit.Email(),
		UserType: string(users.UserTypeAdmin),
	}
	token, err := td.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, token)

	id := primitive.NewObjectID().Hex()

	sbs := &mocks.SubscriptionMock{
		FindOneFunc: func(ctx context.Context, subID string) (*subscription.Subscription, error) {
			return &subscription.Subscription{
				ID:                 id,
				UserID:             td.UserID,
				ProductID:          primitive.NewObjectID().Hex(),
				SubscriptionPlanID: primitive.NewObjectID().Hex(),
				Duration:           5,
				Amount:             2500,
				Tax:                450,
				Total:              2950,
				Status:             subscription.StatusActive,
				Model: gorm.Model{
					CreatedAt: time.Now(),
				},
			}, nil
		},
		UpdateStatusFunc: func(ctx context.Context, subID string, status subscription.Status) error {
			return fmt.Errorf("cannot update successfully")
		},
	}

	svr := New(nil, nil, nil, nil, sbs)
	res := requestHelperWithMockedServer(t, svr, http.MethodPatch, "/subscriptions/"+id+"/pause", createToken(t), nil)
	require.Equal(t, http.StatusUnauthorized, res.Code)
	exp := `{"error":"you are not allowed to access this subscription details","success":false}`
	assert.Equal(t, exp, res.Body.String())
}

func TestPauseSubscription(t *testing.T) {
	ctx := context.Background()
	td := core.TokenData{
		UserID:   primitive.NewObjectID().Hex(),
		Email:    gofakeit.Email(),
		UserType: string(users.UserTypeAdmin),
	}
	token, err := td.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, token)

	sub := &subscription.Subscription{
		UserID:             td.UserID,
		ProductID:          primitive.NewObjectID().Hex(),
		SubscriptionPlanID: primitive.NewObjectID().Hex(),
		Duration:           5,
		Amount:             2500,
		Tax:                450,
		Total:              2950,
		Status:             subscription.StatusActive,
	}
	sbs, err := subscription.New(db)
	require.NoError(t, err)
	require.NoError(t, sbs.Create(ctx, sub))

	svr := New(nil, nil, nil, nil, sbs)
	res := requestHelperWithMockedServer(t, svr, http.MethodPatch, "/subscriptions/"+sub.ID+"/pause", token, nil)
	require.Equal(t, http.StatusOK, res.Code)

	result, err := sbs.FindByStatus(ctx, td.UserID, subscription.StatusPaused)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Len(t, result, 1)
}

func TestUnPauseSubscription_FailedUpdate(t *testing.T) {
	td := core.TokenData{
		UserID:   uuid.NewString(),
		Email:    gofakeit.Email(),
		UserType: string(users.UserTypeAdmin),
	}
	token, err := td.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, token)

	id := primitive.NewObjectID().Hex()

	sbs := &mocks.SubscriptionMock{
		FindOneFunc: func(ctx context.Context, subID string) (*subscription.Subscription, error) {
			return &subscription.Subscription{
				ID:                 id,
				UserID:             td.UserID,
				ProductID:          primitive.NewObjectID().Hex(),
				SubscriptionPlanID: primitive.NewObjectID().Hex(),
				Duration:           5,
				Amount:             2500,
				Tax:                450,
				Total:              2950,
				Model: gorm.Model{
					CreatedAt: time.Now(),
				},
			}, nil
		},
		UpdateStatusFunc: func(ctx context.Context, subID string, status subscription.Status) error {
			return fmt.Errorf("cannot update successfully")
		},
	}

	svr := New(nil, nil, nil, nil, sbs)
	res := requestHelperWithMockedServer(t, svr, http.MethodPatch, "/subscriptions/"+id+"/unpause", token, nil)
	require.Equal(t, http.StatusInternalServerError, res.Code)
	exp := `{"error":"cannot update successfully","success":false}`
	assert.Equal(t, exp, res.Body.String())
}

func TestUnPauseSubscription_RandomToken(t *testing.T) {
	td := core.TokenData{
		UserID:   uuid.NewString(),
		Email:    gofakeit.Email(),
		UserType: string(users.UserTypeAdmin),
	}
	token, err := td.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, token)

	id := primitive.NewObjectID().Hex()

	sbs := &mocks.SubscriptionMock{
		FindOneFunc: func(ctx context.Context, subID string) (*subscription.Subscription, error) {
			return &subscription.Subscription{
				ID:                 id,
				UserID:             td.UserID,
				ProductID:          primitive.NewObjectID().Hex(),
				SubscriptionPlanID: primitive.NewObjectID().Hex(),
				Duration:           5,
				Amount:             2500,
				Tax:                450,
				Total:              2950,
				Model: gorm.Model{
					CreatedAt: time.Now(),
				},
			}, nil
		},
		UpdateStatusFunc: func(ctx context.Context, subID string, status subscription.Status) error {
			return fmt.Errorf("cannot update successfully")
		},
	}

	svr := New(nil, nil, nil, nil, sbs)
	res := requestHelperWithMockedServer(t, svr, http.MethodPatch, "/subscriptions/"+id+"/unpause", createToken(t), nil)
	require.Equal(t, http.StatusUnauthorized, res.Code)
	exp := `{"error":"you are not allowed to access this subscription details","success":false}`
	assert.Equal(t, exp, res.Body.String())
}

func TestUnPauseSubscription(t *testing.T) {
	ctx := context.Background()
	td := core.TokenData{
		UserID:   primitive.NewObjectID().Hex(),
		Email:    gofakeit.Email(),
		UserType: string(users.UserTypeAdmin),
	}
	token, err := td.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, token)

	sub := &subscription.Subscription{
		UserID:             td.UserID,
		ProductID:          primitive.NewObjectID().Hex(),
		SubscriptionPlanID: primitive.NewObjectID().Hex(),
		Duration:           5,
		Amount:             2500,
		Tax:                450,
		Total:              2950,
		Status:             subscription.StatusPaused,
	}
	sbs, err := subscription.New(db)
	require.NoError(t, err)
	require.NoError(t, sbs.Create(ctx, sub))

	svr := New(nil, nil, nil, nil, sbs)
	res := requestHelperWithMockedServer(t, svr, http.MethodPatch, "/subscriptions/"+sub.ID+"/unpause", token, nil)
	require.Equal(t, http.StatusOK, res.Code)

	result, err := sbs.FindByStatus(ctx, td.UserID, subscription.StatusPaused)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Len(t, result, 0)

	result, err = sbs.FindByStatus(ctx, td.UserID, subscription.StatusActive)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Len(t, result, 1)
}

func TestCancelSubscription_FailedUpdate(t *testing.T) {
	td := core.TokenData{
		UserID:   uuid.NewString(),
		Email:    gofakeit.Email(),
		UserType: string(users.UserTypeAdmin),
	}
	token, err := td.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, token)

	id := primitive.NewObjectID().Hex()

	sbs := &mocks.SubscriptionMock{
		FindOneFunc: func(ctx context.Context, subID string) (*subscription.Subscription, error) {
			return &subscription.Subscription{
				ID:                 id,
				UserID:             td.UserID,
				ProductID:          primitive.NewObjectID().Hex(),
				SubscriptionPlanID: primitive.NewObjectID().Hex(),
				Duration:           5,
				Amount:             2500,
				Tax:                450,
				Total:              2950,
				Model: gorm.Model{
					CreatedAt: time.Now(),
				},
			}, nil
		},
		UpdateStatusFunc: func(ctx context.Context, subID string, status subscription.Status) error {
			return fmt.Errorf("cannot update successfully")
		},
	}

	svr := New(nil, nil, nil, nil, sbs)
	res := requestHelperWithMockedServer(t, svr, http.MethodPatch, "/subscriptions/"+id+"/pause", token, nil)
	require.Equal(t, http.StatusInternalServerError, res.Code)
	exp := `{"error":"cannot update successfully","success":false}`
	assert.Equal(t, exp, res.Body.String())
}
func TestCancelSubscription(t *testing.T) {
	ctx := context.Background()
	td := core.TokenData{
		UserID:   primitive.NewObjectID().Hex(),
		Email:    gofakeit.Email(),
		UserType: string(users.UserTypeAdmin),
	}
	token, err := td.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, token)

	sub := &subscription.Subscription{
		UserID:             td.UserID,
		ProductID:          primitive.NewObjectID().Hex(),
		SubscriptionPlanID: primitive.NewObjectID().Hex(),
		Duration:           5,
		Amount:             2500,
		Tax:                450,
		Total:              2950,
		Status:             subscription.StatusActive,
	}
	sbs, err := subscription.New(db)
	require.NoError(t, err)
	require.NoError(t, sbs.Create(ctx, sub))

	svr := New(nil, nil, nil, nil, sbs)
	res := requestHelperWithMockedServer(t, svr, http.MethodDelete, "/subscriptions/"+sub.ID, token, nil)
	require.Equal(t, http.StatusOK, res.Code)

	result, err := sbs.FindByStatus(ctx, td.UserID, subscription.StatusPaused)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Len(t, result, 0)

	result, err = sbs.FindByStatus(ctx, td.UserID, subscription.StatusCancelled)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Len(t, result, 1)
}

func TestCancelSubscription_WithRandomToken(t *testing.T) {
	ctx := context.Background()
	td := core.TokenData{
		UserID:   primitive.NewObjectID().Hex(),
		Email:    gofakeit.Email(),
		UserType: string(users.UserTypeAdmin),
	}
	token, err := td.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, token)

	randomToken := createToken(t)

	sub := &subscription.Subscription{
		UserID:             td.UserID,
		ProductID:          primitive.NewObjectID().Hex(),
		SubscriptionPlanID: primitive.NewObjectID().Hex(),
		Duration:           5,
		Amount:             2500,
		Tax:                450,
		Total:              2950,
		Status:             subscription.StatusActive,
	}
	sbs, err := subscription.New(db)
	require.NoError(t, err)
	require.NoError(t, sbs.Create(ctx, sub))

	svr := New(nil, nil, nil, nil, sbs)
	res := requestHelperWithMockedServer(t, svr, http.MethodDelete, "/subscriptions/"+sub.ID, randomToken, nil)
	require.Equal(t, http.StatusUnauthorized, res.Code)
	exp := `{"error":"you are not allowed to access this subscription details","success":false}`
	require.Equal(t, exp, res.Body.String())
}

func createToken(t *testing.T) string {
	t.Helper()
	td := core.TokenData{
		UserID:   uuid.NewString(),
		Email:    gofakeit.Email(),
		UserType: string(users.UserTypeAdmin),
	}
	token, err := td.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, token)

	return token
}

func requestHelperWithMockedServer(t *testing.T, s *Server, method, path, token string, payload []byte) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	var (
		req *http.Request
		err error
	)

	if len(payload) == 0 {
		req, err = http.NewRequest(method, path, nil)
	} else {
		req, err = http.NewRequest(method, path, bytes.NewBuffer(payload))
		fmt.Printf("\n\nRequest: %s\n\n", payload)
	}

	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	s.Router.ServeHTTP(w, req)
	require.NotNil(t, w)

	fmt.Printf("\n\nResponse: %s\n\n", w.Body.String())
	return w
}
