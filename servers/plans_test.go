package servers

import (
	"context"
	"encoding/json"
	"net/http"
	"subscriptions/domains/plans"
	"subscriptions/requests"
	"subscriptions/responses"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateSubscriptionPlan(t *testing.T) {
	req := requests.SubscriptionPlan{
		ProductID: primitive.NewObjectID().Hex(),
		Amount:    200,
		Duration:  3,
	}

	b, err := json.Marshal(req)
	require.NoError(t, err)
	require.NotNil(t, b)

	token := createToken(t)

	res := requestHelper(t, http.MethodPost, "/plans", token, b)
	require.Equal(t, http.StatusCreated, res.Code)

	var response responses.SubscriptionPlan
	err = json.Unmarshal(res.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.NotEmpty(t, response.ID)
	assert.NotEmpty(t, response.CreatedAt)
	assert.NotEmpty(t, response.UpdatedAt)

	assert.Equal(t, req.ProductID, response.ProductID)
	assert.Equal(t, req.Amount, response.Amount)
	assert.Equal(t, req.Duration, response.Duration)
	assert.Equal(t, req.TrialDuration, response.TrialDuration)
}

func TestCreateSubscriptionPlanWithBadJSON(t *testing.T) {
	b := []byte(`
	{
		"product_id": "6468cf5b8900f1f8a43766d1",
		"amount": 200,
		"duration": 3,
		"trial_duration": 0
	  `)

	token := createToken(t)

	res := requestHelper(t, http.MethodPost, "/plans", token, b)
	require.Equal(t, http.StatusBadRequest, res.Code)
}

func TestPlanDetails(t *testing.T) {
	ctx := context.Background()
	productID := primitive.NewObjectID().Hex()
	pPlans := make([]*plans.SubscriptionPlan, 3)

	for i := 0; i < 3; i++ {
		v := &plans.SubscriptionPlan{
			ProductID:     productID,
			Amount:        200 * float64(i),
			Duration:      3,
			TrialDuration: 0,
		}
		require.NoError(t, planService.Create(ctx, v))
		pPlans[i] = v
	}

	id := pPlans[1].ID

	res := requestHelper(t, http.MethodGet, "/plans/"+id, "", nil)
	require.Equal(t, http.StatusOK, res.Code)

	var response responses.SubscriptionPlan
	err := json.Unmarshal(res.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.NotEmpty(t, response.ID)
	assert.NotEmpty(t, response.CreatedAt)
	assert.NotEmpty(t, response.UpdatedAt)

	assert.Equal(t, pPlans[1].ProductID, response.ProductID)
	assert.Equal(t, pPlans[1].Amount, response.Amount)
	assert.Equal(t, pPlans[1].Duration, response.Duration)
	assert.Equal(t, pPlans[1].TrialDuration, response.TrialDuration)
}
