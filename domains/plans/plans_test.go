package plans

import (
	"context"
	"os"
	"testing"

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

func TestCreateNewSubscriptionPlan(t *testing.T) {
	ctx := context.Background()

	p := newSubscriptionPlan(t)

	t.Cleanup(func() {
		println("cleaning up")
		db.Exec("DELETE FROM subscriptionplans WHERE id = ?", p.ID)
	})

	ps, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, ps)

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

	t.Cleanup(func() {
		ids := make([]string, len(pps))
		println("pps len", len(pps))
		for i := range pps {
			ids[i] = pps[i].ID
		}
		db.Exec("DELETE FROM plans WHERE id IN ?", ids)
	})

	for i := 1; i <= 3; i++ {
		p := newSubscriptionPlan(t)
		require.NoError(t, ps.Create(ctx, p))
		pps = append(pps, p)
	}

	single, err := ps.FindOne(ctx, pps[1].ID)
	require.NoError(t, err)
	require.NotEmpty(t, single)

	// assert.Equal(t, pps[1].ID, single.ID)
	// assert.Equal(t, pps[1].Name, single.Name)
	// assert.Equal(t, pps[1].Description, single.Description)
	// assert.Equal(t, pps[1].Tax, single.Tax)
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

func newSubscriptionPlan(t *testing.T) *SubscriptionPlan {
	t.Helper()
	return &SubscriptionPlan{
		ProductID: primitive.NewObjectID().Hex(),
	}
}

func cleanup() {
	db.Exec("DELETE FROM subscription_plans")
}
