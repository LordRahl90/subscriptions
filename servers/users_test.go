package servers

import (
	"encoding/json"
	"net/http"
	"testing"

	"subscriptions/requests"
	"subscriptions/responses"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateNewUser(t *testing.T) {
	u := createUser(t)

	b, err := json.Marshal(u)
	require.NoError(t, err)
	require.NotNil(t, b)

	res := requestHelper(t, http.MethodPost, "/user/create", "", b)
	require.Equal(t, http.StatusCreated, res.Code)

	var r responses.User
	err = json.Unmarshal(res.Body.Bytes(), &r)
	require.NoError(t, err)
	require.NotNil(t, r)

	assert.NotEmpty(t, r.ID)
	assert.Equal(t, u.Email, r.Email)
	assert.Equal(t, u.Name, r.Name)
	assert.NotEmpty(t, r.Token)
	assert.NotEmpty(t, r.CreatedAt)
}

func TestCreateUserWithBadJSON(t *testing.T) {
	b := []byte(`{
		"name": "Bart Beatty",
		"email": "cordiajacobi@carroll.net",
		"password": "password",
		"gender": "male",
		"age": 49,
	}`)
	res := requestHelper(t, http.MethodPost, "/user/create", "", b)
	require.Equal(t, http.StatusBadRequest, res.Code)
}

func TestAuthenticateUser(t *testing.T) {
	u := createUser(t)

	b, err := json.Marshal(u)
	require.NoError(t, err)
	require.NotNil(t, b)

	res := requestHelper(t, http.MethodPost, "/user/create", "", b)
	require.Equal(t, http.StatusCreated, res.Code)

	authReq := requests.Authenticate{
		Email:    u.Email,
		Password: "password",
	}

	b, err = json.Marshal(authReq)
	require.NoError(t, err)
	require.NotNil(t, b)

	res = requestHelper(t, http.MethodPost, "/login", "", b)
	require.Equal(t, http.StatusOK, res.Code)

	var r responses.User
	err = json.Unmarshal(res.Body.Bytes(), &r)
	require.NoError(t, err)
	require.NotNil(t, r)

	assert.NotEmpty(t, r.ID)
	assert.Equal(t, u.Email, r.Email)
	assert.Equal(t, u.Name, r.Name)
	assert.NotEmpty(t, r.Token)
	assert.NotEmpty(t, r.CreatedAt)
}

func TestAuthenticateWithNonExistingEmail(t *testing.T) {
	authReq := requests.Authenticate{
		Email:    gofakeit.Email(),
		Password: "password",
	}

	b, err := json.Marshal(authReq)
	require.NoError(t, err)
	require.NotNil(t, b)

	res := requestHelper(t, http.MethodPost, "/login", "", b)
	require.Equal(t, http.StatusBadRequest, res.Code)
	exp := `{"data":null,"message":"email not registered","success":false}`
	require.Equal(t, exp, res.Body.String())
}

func TestAuthenticationWithBadPassword(t *testing.T) {
	u := createUser(t)

	b, err := json.Marshal(u)
	require.NoError(t, err)
	require.NotNil(t, b)

	res := requestHelper(t, http.MethodPost, "/user/create", "", b)
	require.Equal(t, http.StatusCreated, res.Code)

	authReq := requests.Authenticate{
		Email:    u.Email,
		Password: "passwordss",
	}

	b, err = json.Marshal(authReq)
	require.NoError(t, err)
	require.NotNil(t, b)

	res = requestHelper(t, http.MethodPost, "/login", "", b)
	println(res.Body.String())
	require.Equal(t, http.StatusBadRequest, res.Code)
	exp := `{"data":null,"message":"password does not match","success":false}`
	require.Equal(t, exp, res.Body.String())
}

func TestAuthenticateWithBadJSON(t *testing.T) {
	b := []byte(`{
		"email": "cordiajacobi@carroll.net",
		"password": "password",
	}`)

	res := requestHelper(t, http.MethodPost, "/login", "", b)
	require.Equal(t, http.StatusBadRequest, res.Code)
}

func createUser(t *testing.T) *requests.CreateUser {
	t.Helper()
	return &requests.CreateUser{
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Password: "password",
	}
}
