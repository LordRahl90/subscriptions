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
	"subscriptions/domains/vouchers"
	"subscriptions/requests"
	"subscriptions/responses"
	"subscriptions/services/purchase/mocks"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

	res = requestHelperWithMockedServer(t, svr, http.MethodGet, "/subscriptions", token, nil)
	require.Equal(t, http.StatusOK, res.Code)

	var results []responses.Subscription
	err = json.Unmarshal(res.Body.Bytes(), &results)
	require.NoError(t, err)
	require.Len(t, results, 1)

	assert.Equal(t, response.ID, results[0].ID)
	assert.Equal(t, response.Total, results[0].Total)
	assert.Equal(t, response.Tax, results[0].Tax)
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
	assert.Equal(t, uint(3), response.Duration)
	assert.Equal(t, uint(1), response.TrialDuration)
	assert.Equal(t, 400.0, response.Tax)
	assert.Equal(t, 0.0, response.Discount)
	assert.Equal(t, 2400.0, response.Total)
}

func createToken(t *testing.T) string {
	t.Helper()
	td := core.TokenData{
		UserID: uuid.NewString(),
		Email:  gofakeit.Email(),
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
