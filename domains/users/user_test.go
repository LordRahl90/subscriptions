package users

import (
	"context"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		cleanup()
		os.Exit(code)
	}()
	db = setupTestDB()
	code = m.Run()
}

func TestCreateNewUser(t *testing.T) {
	ctx := context.Background()

	u := newUser(t)
	us, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, us)

	err = us.Create(ctx, u)
	require.NoError(t, err)
	require.NotEmpty(t, u.ID)
}

func TestCreateDuplicateUser(t *testing.T) {
	ctx := context.Background()
	us, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, us)

	u := newUser(t)
	err = us.Create(ctx, u)
	require.NoError(t, err)
	require.NotEmpty(t, u.ID)

	usr := newUser(t)
	usr.Email = u.Email

	err = us.Create(ctx, usr)
	require.NoError(t, err)
}

func TestFindUser(t *testing.T) {
	ctx := context.Background()
	var (
		values  []string
		records = make(map[string]*User)
	)

	us, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, us)

	for i := 0; i <= 3; i++ {
		u := newUser(t)
		err := us.Create(ctx, u)
		require.NoError(t, err)
		require.NotEmpty(t, u.ID)
		records[u.ID] = u
		values = append(values, u.ID)
	}

	id := values[2]
	res, err := us.Find(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestFindNonExistentUser(t *testing.T) {
	ctx := context.Background()

	us, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, us)

	res, err := us.Find(ctx, primitive.NewObjectID().Hex())
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	require.Empty(t, res)
}

func TestFindUserByEmail(t *testing.T) {
	ctx := context.Background()
	var (
		values  []string
		records = make(map[string]*User)
	)

	us, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, us)

	for i := 0; i <= 3; i++ {
		u := newUser(t)
		err := us.Create(ctx, u)
		require.NoError(t, err)
		require.NotEmpty(t, u.ID)
		records[u.ID] = u
		values = append(values, u.ID)
	}

	id := values[2]
	rec := records[id]
	res, err := us.FindByEmail(ctx, rec.Email)
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, rec.ID, res.ID)
	assert.Equal(t, rec.Email, res.Email)
	assert.Equal(t, rec.Name, res.Name)
}

func setupTestDB() *gorm.DB {
	env := os.Getenv("ENVIRONMENT")
	dsn := "root:@tcp(127.0.0.1:3306)/subscriptions?charset=utf8mb4&parseTime=True&loc=Local"
	if env == "cicd" {
		dsn = "test_user:password@tcp(127.0.0.1:33306)/subscriptions?charset=utf8mb4&parseTime=True&loc=Local"
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func cleanup() {
	db.Exec("DELETE FROM users")
}

func newUser(t *testing.T) *User {
	t.Helper()
	return &User{
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Password: "password",
	}
}
