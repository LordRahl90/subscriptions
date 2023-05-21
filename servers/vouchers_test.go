package servers

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"subscriptions/requests"
	"subscriptions/responses"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateVoucherWithAmount(t *testing.T) {
	req := requests.Voucher{
		VoucherType: "amount",
		ProductID:   primitive.NewObjectID().Hex(),
		Code:        "1234567890",
		Amount:      200,
		ExpiresOn:   time.Now().Add(48 * time.Hour),
	}

	b, err := json.Marshal(req)
	require.NoError(t, err)
	require.NotEmpty(t, b)

	token := createToken(t)

	res := requestHelper(t, http.MethodPost, "/vouchers", token, b)
	require.Equal(t, http.StatusCreated, res.Code)

	var response responses.Voucher
	err = json.Unmarshal(res.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotEmpty(t, response.ID)
	assert.Equal(t, req.Code, response.Code)
	assert.Equal(t, req.ProductID, response.ProductID)
	assert.Equal(t, req.VoucherType, response.VoucherType)
	assert.Equal(t, req.Amount, response.Amount)
}

func TestCreateVoucherWithPercentage(t *testing.T) {
	req := requests.Voucher{
		VoucherType: "percentage",
		ProductID:   primitive.NewObjectID().Hex(),
		Code:        "1234567890",
		Percentage:  20,
		ExpiresOn:   time.Now().Add(24 * time.Hour),
	}

	b, err := json.Marshal(req)
	require.NoError(t, err)
	require.NotEmpty(t, b)

	token := createToken(t)

	res := requestHelper(t, http.MethodPost, "/vouchers", token, b)
	require.Equal(t, http.StatusCreated, res.Code)

	var response responses.Voucher
	err = json.Unmarshal(res.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotEmpty(t, response.ID)
	assert.Equal(t, req.Code, response.Code)
	assert.Equal(t, req.ProductID, response.ProductID)
	assert.Equal(t, req.VoucherType, response.VoucherType)
	assert.Equal(t, req.Percentage, response.Percentage)
}

func TestCreateVoucherWithInvalidPercentage(t *testing.T) {
	req := requests.Voucher{
		VoucherType: "percentage",
		ProductID:   primitive.NewObjectID().Hex(),
		Code:        "1234567890",
		Percentage:  200,
	}

	b, err := json.Marshal(req)
	require.NoError(t, err)
	require.NotEmpty(t, b)

	token := createToken(t)

	res := requestHelper(t, http.MethodPost, "/vouchers", token, b)
	require.Equal(t, http.StatusBadRequest, res.Code)
	println(res.Body.String())
}
