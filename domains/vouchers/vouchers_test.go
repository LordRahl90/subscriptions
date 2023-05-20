package vouchers

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit"
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

func TestCreateNewVoucher(t *testing.T) {
	ctx := context.Background()

	p := newVoucher(t)

	t.Cleanup(func() {
		println("cleaning up")
		db.Exec("DELETE FROM vouchers WHERE id = ?", p.ID)
	})

	ps, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, ps)

	err = ps.Create(ctx, p)
	require.NoError(t, err)
	require.NotEmpty(t, p.ID)
}

func TestFindVouchers(t *testing.T) {
	ctx := context.Background()
	ps, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, ps)

	pps := []*Voucher{}

	t.Cleanup(func() {
		ids := make([]string, len(pps))
		println("pps len", len(pps))
		for i := range pps {
			ids[i] = pps[i].ID
		}
		db.Exec("DELETE FROM vouchers WHERE id IN ?", ids)
	})

	for i := 1; i <= 3; i++ {
		p := newVoucher(t)
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

func TestValidateVoucher(t *testing.T) {
	table := []struct {
		name     string
		args     Voucher
		expError bool
		errMsg   string
	}{
		{
			name: "unknown voucher",
			args: Voucher{
				VoucherType: VoucherTypeUnknown,
				Amount:      100,
			},
			expError: true,
			errMsg:   "unknown voucher",
		},
		{
			name: "clean amount voucher",
			args: Voucher{
				VoucherType: VoucherTypeAmount,
				Amount:      100,
			},
			expError: false,
		},
		{
			name: "no amount voucher",
			args: Voucher{
				VoucherType: VoucherTypeAmount,
				Percentage:  100,
			},
			expError: true,
			errMsg:   "amount voucher should have amount",
		},
		{
			name: "clean percentage voucher",
			args: Voucher{
				VoucherType: VoucherTypePercentage,
				Percentage:  20,
			},
			expError: false,
		},
		{
			name: "no percentage voucher",
			args: Voucher{
				VoucherType: VoucherTypePercentage,
				Amount:      20,
			},
			expError: true,
			errMsg:   "percent voucher should have the percentage",
		},
		{
			name: "overage percentage voucher",
			args: Voucher{
				VoucherType: VoucherTypePercentage,
				Percentage:  101,
			},
			expError: true,
			errMsg:   "percentage cannot be more than 100",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.Validate()
			if tt.expError {
				require.EqualError(t, err, tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}

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

func newVoucher(t *testing.T) *Voucher {
	t.Helper()
	return &Voucher{
		ProductID:   primitive.NewObjectID().Hex(),
		VoucherType: VoucherTypePercentage,
		Code:        strings.ToUpper(gofakeit.BS()),
		Percentage:  20,
	}
}

func cleanup() {
	db.Exec("DELETE FROM vouchers")
}
