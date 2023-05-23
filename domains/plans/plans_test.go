package plans

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db        *gorm.DB
	initError error
)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		cleanup()
		os.Exit(code)
	}()
	db, initError = setupTestDB()
	if initError != nil {
		log.Fatal(initError)
	}
	code = m.Run()
}

func TestCreateNewSubscriptionPlan(t *testing.T) {
	ctx := context.Background()

	p := &SubscriptionPlan{
		ProductID:     primitive.NewObjectID().Hex(),
		Amount:        100,
		Duration:      3,
		TrialDuration: 1,
	}

	ps, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, ps)

	t.Cleanup(func() {
		if err := db.Exec("DELETE FROM subscription_plans WHERE id = ?", p.ID).Error; err != nil {
			log.Fatal(err)
		}
	})

	err = ps.Create(ctx, p)
	require.NoError(t, err)
	require.NotEmpty(t, p.ID)
}

func TestFindSubscriptionPlans(t *testing.T) {
	ctx := context.Background()
	ps, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, ps)

	pps := []*SubscriptionPlan{}
	productID := primitive.NewObjectID().Hex()

	t.Cleanup(func() {
		ids := make([]string, len(pps))
		for i := range pps {
			ids[i] = pps[i].ID
		}
		if err := db.Exec("DELETE FROM subscription_plans WHERE id IN ?", ids).Error; err != nil {
			log.Fatal(err)
		}
	})

	for i := 1; i <= 3; i++ {
		p := newSubscriptionPlan(t, productID)
		require.NoError(t, ps.Create(ctx, p))
		pps = append(pps, p)
	}

	res, err := ps.Find(ctx, productID)
	require.NoError(t, err)
	assert.Len(t, res, 3)

	single, err := ps.FindOne(ctx, pps[1].ID)
	require.NoError(t, err)
	require.NotEmpty(t, single)

	assert.Equal(t, pps[1].ID, single.ID)
	assert.Equal(t, pps[1].ProductID, single.ProductID)
	assert.Equal(t, pps[1].Amount, single.Amount)
	assert.Equal(t, pps[1].Duration, single.Duration)
	assert.Equal(t, pps[1].TrialDuration, single.TrialDuration)
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

func newSubscriptionPlan(t *testing.T, productID string) *SubscriptionPlan {
	t.Helper()
	return &SubscriptionPlan{
		ProductID: productID,
		Amount:    200,
		Duration:  4,
	}
}

func cleanup() {
	db.Exec("DELETE FROM subscription_plans")
}
