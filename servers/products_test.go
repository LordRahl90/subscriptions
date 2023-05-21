package servers

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"subscriptions/domains/plans"
	"subscriptions/domains/products"
	"subscriptions/requests"
	"subscriptions/responses"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateNewProduct(t *testing.T) {
	req := requests.Product{
		Name:        gofakeit.BuzzWord(),
		Description: gofakeit.Word(),
		Tax:         gofakeit.Float64Range(10, 25),
		TrialExists: true,
	}

	b, err := json.Marshal(req)
	require.NoError(t, err)
	require.NotEmpty(t, b)

	res := requestHelper(t, http.MethodPost, "/products", "", b)
	require.Equal(t, http.StatusCreated, res.Code)

	var response responses.Product
	err = json.Unmarshal(res.Body.Bytes(), &response)
	require.NoError(t, err)
	require.NotEmpty(t, response)

	assert.NotEmpty(t, response.ID)
	assert.NotEmpty(t, response.CreatedAt)
	assert.NotEmpty(t, response.UpdatedAt)
	assert.Equal(t, req.Name, response.Name)
	assert.Equal(t, req.Description, response.Description)
	assert.Equal(t, req.Tax, response.Tax)
	assert.Equal(t, req.TrialExists, response.TrialExists)

	db.Exec("DELETE FROM products WHERE id = ?", response.ID)
}

func TestCreateWithInvalidJSON(t *testing.T) {
	b := []byte(`
	{
		"name": "system engine",
		"description": "enim",
		"tax": 12.981642025655608,
		"trial_exists": true
	  peos`)

	res := requestHelper(t, http.MethodPost, "/products", "", b)
	require.Equal(t, http.StatusBadRequest, res.Code)
}

func TestGetAllProducts(t *testing.T) {
	ctx := context.Background()
	for i := 1; i <= 3; i++ {
		require.NoError(t, productService.Create(ctx, &products.Product{
			Name:        gofakeit.BuzzWord(),
			Description: gofakeit.Word(),
			TrialExists: true,
			TaxRate:     10,
		}))
	}

	res := requestHelper(t, http.MethodGet, "/products", "", nil)
	require.Equal(t, http.StatusOK, res.Code)
	var results []responses.Product
	err := json.Unmarshal(res.Body.Bytes(), &results)
	require.NoError(t, err)

	assert.Len(t, results, 3)

	id := results[1].ID
	res = requestHelper(t, http.MethodGet, "/products/"+id, "", nil)
	require.Equal(t, http.StatusOK, res.Code)

	var response responses.Product
	err = json.Unmarshal(res.Body.Bytes(), &response)
	require.NoError(t, err)
	require.NotEmpty(t, response)

	assert.NotEmpty(t, response.ID)
	assert.NotEmpty(t, response.CreatedAt)
	assert.NotEmpty(t, response.UpdatedAt)
	assert.Equal(t, results[1].Name, response.Name)
	assert.Equal(t, results[1].Description, response.Description)
	assert.Equal(t, results[1].Tax, response.Tax)
	assert.Equal(t, results[1].TrialExists, response.TrialExists)
}

func TestGetProductPlans(t *testing.T) {
	ctx := context.Background()
	p := &products.Product{
		Name:        gofakeit.BuzzWord(),
		Description: gofakeit.Word(),
		TrialExists: true,
		TaxRate:     10,
	}
	require.NoError(t, productService.Create(ctx, p))
	for i := 0; i < 3; i++ {
		v := &plans.SubscriptionPlan{
			ProductID:     p.ID,
			Amount:        200 * float64(i),
			Duration:      3,
			TrialDuration: 0,
		}
		require.NoError(t, planService.Create(ctx, v))
	}

	res := requestHelper(t, http.MethodGet, "/products/"+p.ID+"/plans", "", nil)
	require.Equal(t, http.StatusOK, res.Code)

	var response []responses.SubscriptionPlan

	err := json.Unmarshal(res.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Len(t, response, 3)
}
