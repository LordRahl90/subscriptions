package servers

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	server        *Server
	db            *gorm.DB
	initErr       error
	signingSecret = "hello-world"
)

func TestMain(m *testing.M) {
	var code = 1
	defer func() {
		cleanup()
		os.Exit(code)
	}()
	db, initErr = setupTestDB()
	if initErr != nil {
		log.Fatal(initErr)
	}
	s, err := NewWithDefaults(db)
	if err != nil {
		log.Fatal(err)
	}
	server = s
	server.SigningSecret = signingSecret
	code = m.Run()
}

func requestHelper(t *testing.T, method, path, token string, payload []byte) *httptest.ResponseRecorder {
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
	}

	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	server.Router.ServeHTTP(w, req)
	require.NotNil(t, w)
	return w
}

func setupTestDB() (*gorm.DB, error) {
	env := os.Getenv("ENVIRONMENT")
	dsn := "root:@tcp(127.0.0.1:3306)/subscriptions?charset=utf8mb4&parseTime=True&loc=Local"
	if env == "cicd" {
		dsn = "test_user:password@tcp(127.0.0.1:33306)/subscriptions?charset=utf8mb4&parseTime=True&loc=Local"
	}
	dbase, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := dbase.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxIdleConns(0)

	return dbase, err
}

func cleanup() {
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM products")
	db.Exec("DELETE FROM subscription_plans")
	db.Exec("DELETE FROM vouchers")
	db.Exec("DELETE FROM subscriptions")
}
