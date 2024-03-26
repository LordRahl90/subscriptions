package products

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	if db == nil {
		log.Fatal("empty db")
	}
	code = m.Run()
}

func TestCreateNewProduct(t *testing.T) {
	ctx := context.Background()

	p := newProduct(t)
	p.TaxRate = 34.5

	t.Cleanup(func() {
		if err := db.Exec("DELETE FROM products WHERE id = ?", p.ID).Error; err != nil {
			log.Fatal(err)
		}
	})

	ps, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, ps)

	err = ps.Create(ctx, p)
	require.NoError(t, err)
	require.NotEmpty(t, p.ID)

	res, err := ps.Find(ctx)
	require.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestFindProducts(t *testing.T) {
	ctx := context.Background()
	ps, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, ps)

	pps := []*Product{}

	t.Cleanup(func() {
		ids := make([]string, len(pps))
		for i := range pps {
			ids[i] = pps[i].ID
		}
		if err := db.Exec("DELETE FROM products WHERE id IN ?", ids).Error; err != nil {
			log.Fatal(err)
		}
	})

	for i := 1; i <= 3; i++ {
		p := newProduct(t)
		require.NoError(t, ps.Create(ctx, p))
		pps = append(pps, p)
	}

	res, err := ps.Find(ctx)
	require.NoError(t, err)
	assert.Len(t, res, 3)

	single, err := ps.FindOne(ctx, pps[1].ID)
	require.NoError(t, err)
	require.NotEmpty(t, single)

	assert.Equal(t, pps[1].ID, single.ID)
	assert.Equal(t, pps[1].Name, single.Name)
	assert.Equal(t, pps[1].Description, single.Description)
	assert.Equal(t, pps[1].TaxRate, single.TaxRate)

	result, err := ps.FindByIDs(ctx, pps[1].ID, pps[2].ID)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	assert.Len(t, result, 2)
}

func FuzzCreateProduct(f *testing.F) {
	ctx := context.Background()
	ps, err := New(db)
	if err != nil || ps == nil {
		log.Fatal(err)
	}

	f.Fuzz(func(t *testing.T, name, description string, tax float64) {
		p := &Product{
			Name:        name,
			Description: description,
			TaxRate:     tax,
		}

		err := ps.Create(ctx, p)
		if err != nil {
			log.Fatal(err)
		}
	})
}

func setupTestDB() *gorm.DB {
	env := os.Getenv("ENVIRONMENT")
	dsn := "root:password@tcp(127.0.0.1:3306)/subscriptions?charset=utf8mb4&parseTime=True&loc=Local"
	if env == "cicd" {
		dsn = "test_user:password@tcp(127.0.0.1:33306)/subscriptions?charset=utf8mb4&parseTime=True&loc=Local"
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func newProduct(t *testing.T) *Product {
	t.Helper()
	return &Product{
		Name:        gofakeit.Name(),
		Description: gofakeit.Word(),
		TaxRate:     25.0,
	}
}

func cleanup() {
	if err := db.Exec("DELETE FROM products where id!=''").Error; err != nil {
		log.Fatal(err)
	}
}
