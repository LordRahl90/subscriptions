package subscription

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		cleanup()
		os.Exit(code)
	}()
	db = setupTestDB()
	code = m.Run()
}

func TestCreateNewSubscription(t *testing.T) {
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()

	s := &Subscription{
		UserID:             userID,
		ProductID:          primitive.NewObjectID().Hex(),
		SubscriptionPlanID: primitive.NewObjectID().Hex(),
		Duration:           4,
		Amount:             2000,
		Discount:           500,
		Tax:                200,
		TrialPeriod:        1,
		Status:             StatusActive,
	}

	t.Cleanup(func() {
		println("cleaning up")
		db.Exec("DELETE FROM subscriptions WHERE id = ?", s.ID)
	})

	ps, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, ps)

	err = ps.Create(ctx, s)
	require.NoError(t, err)
	require.NotEmpty(t, s.ID)

	res, err := ps.Find(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestFindSubscriptions(t *testing.T) {
	ctx := context.Background()
	ps, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, ps)

	pps := []*Subscription{}
	userID := primitive.NewObjectID().Hex()

	t.Cleanup(func() {
		ids := make([]string, len(pps))
		for i := range pps {
			ids[i] = pps[i].ID
		}
		db.Exec("DELETE FROM subscriptions WHERE id IN ?", ids)
	})

	for i := 1; i <= 3; i++ {
		p := newSubscription(t, userID, StatusActive)
		require.NoError(t, ps.Create(ctx, p))
		pps = append(pps, p)
	}

	res, err := ps.Find(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, res, 3)

	single, err := ps.FindOne(ctx, pps[1].ID)
	require.NoError(t, err)
	require.NotEmpty(t, single)

	assert.Equal(t, pps[1].ID, single.ID)
	assert.Equal(t, pps[1].UserID, single.UserID)
	assert.Equal(t, pps[1].Tax, single.Tax)
	assert.Equal(t, pps[1].Duration, single.Duration)
	assert.Equal(t, pps[1].Discount, single.Discount)
}

func TestFindSubscriptionsByStatus(t *testing.T) {
	ctx := context.Background()
	ps, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, ps)

	pps := []*Subscription{}
	userID := primitive.NewObjectID().Hex()
	status := []Status{StatusActive, StatusPaused, StatusCancelled, StatusEnded, StatusActive, StatusEnded}

	t.Cleanup(func() {
		ids := make([]string, len(pps))
		for i := range pps {
			ids[i] = pps[i].ID
		}
		db.Exec("DELETE FROM subscriptions WHERE id IN ?", ids)
	})

	for i := range status {
		p := newSubscription(t, userID, status[i])
		require.NoError(t, ps.Create(ctx, p))
		pps = append(pps, p)
	}

	res, err := ps.FindByStatus(ctx, userID, StatusActive)
	require.NoError(t, err)
	require.Len(t, res, 2)

	res, err = ps.FindByStatus(ctx, userID, StatusPaused)
	require.NoError(t, err)
	require.Len(t, res, 1)

	res, err = ps.FindByStatus(ctx, userID, StatusCancelled)
	require.NoError(t, err)
	require.Len(t, res, 1)

	res, err = ps.FindByStatus(ctx, userID, StatusEnded)
	require.NoError(t, err)
	require.Len(t, res, 2)
}

func TestUpdateStatus(t *testing.T) {
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()

	s := &Subscription{
		UserID:             userID,
		ProductID:          primitive.NewObjectID().Hex(),
		SubscriptionPlanID: primitive.NewObjectID().Hex(),
		Duration:           4,
		Amount:             2000,
		Discount:           500,
		Tax:                200,
		TrialPeriod:        1,
		Status:             StatusActive,
	}

	t.Cleanup(func() {
		println("cleaning up")
		db.Exec("DELETE FROM subscriptions WHERE id = ?", s.ID)
	})

	ps, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, ps)

	err = ps.Create(ctx, s)
	require.NoError(t, err)
	require.NotEmpty(t, s.ID)

	err = ps.UpdateStatus(ctx, s.ID, StatusPaused)
	require.NoError(t, err)

	res, err := ps.FindOne(ctx, s.ID)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	assert.Equal(t, StatusPaused, res.Status)

	results, err := ps.FindByStatus(ctx, userID, StatusPaused)
	require.NoError(t, err)
	require.Len(t, results, 1)
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

func newSubscription(t *testing.T, userID string, status Status) *Subscription {
	t.Helper()
	return &Subscription{
		UserID:             userID,
		ProductID:          primitive.NewObjectID().Hex(),
		SubscriptionPlanID: primitive.NewObjectID().Hex(),
		Duration:           3,
		Amount:             2000,
		Discount:           200,
		Tax:                200,
		Status:             status,
	}
}

func cleanup() {
	db.Exec("DELETE FROM subscriptions")
}
