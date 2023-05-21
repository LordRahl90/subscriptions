package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"subscriptions/domains/plans"
	"subscriptions/domains/products"
	"subscriptions/domains/users"
	"subscriptions/domains/vouchers"

	"github.com/brianvoe/gofakeit"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	ctx := context.Background()
	db, err := setupDB()
	if err != nil {
		panic(err)
	}
	if err := loadProductAndPlans(ctx, db); err != nil {
		panic(err)
	}

	if err := loadUsers(ctx, db); err != nil {
		panic(err)
	}
	fmt.Println("seeding completed")
}

func loadProductAndPlans(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("db not initialized")
	}
	productService, err := products.New(db)
	if err != nil {
		return err
	}

	planService, err := plans.New(db)
	if err != nil {
		return err
	}

	voucherService, err := vouchers.New(db)
	if err != nil {
		return err
	}

	res, err := productService.Find(ctx)
	if err != nil {
		return err
	}
	if len(res) > 0 {
		panic("there's some data populated already")
	}

	for i := 1; i <= 10; i++ {
		p := &products.Product{
			Name:        gofakeit.BeerName(),
			Description: gofakeit.BS(),
			TaxRate:     gofakeit.Float64Range(20, 40),
		}
		if err := productService.Create(ctx, p); err != nil {
			return err
		}

		for j := 1; j < 3; j++ {
			plan := &plans.SubscriptionPlan{
				ProductID:     p.ID,
				Amount:        gofakeit.Float64Range(1000, 2000),
				Duration:      uint(gofakeit.Float32Range(3, 6)),
				TrialDuration: uint(gofakeit.Float64Range(0, 2)),
			}

			if err := planService.Create(ctx, plan); err != nil {
				return err
			}
		}

		v := []*vouchers.Voucher{
			{
				ProductID:   p.ID,
				Code:        "WELCOME_FIRST_100",
				VoucherType: vouchers.VoucherTypeAmount,
				Amount:      200,
				Limit:       100,
				ExpiresOn:   time.Now().Add(24 * 3 * time.Hour),
			},
			{
				ProductID:   p.ID,
				Code:        "WELCOME_FIRST_101",
				VoucherType: vouchers.VoucherTypePercentage,
				Percentage:  20,
				Limit:       100,
				ExpiresOn:   time.Now().Add(24 * 3 * time.Hour),
			},
			{
				ProductID:   p.ID,
				Code:        "WELCOME_STAFF_ONLY",
				VoucherType: vouchers.VoucherTypePercentage,
				Percentage:  15,
				Limit:       30,
				ExpiresOn:   time.Now().Add(24 * 3 * time.Hour),
			},
		}

		for j := range v {
			if err := voucherService.Create(ctx, v[j]); err != nil {
				return err
			}
		}
	}

	return nil
}

func loadUsers(ctx context.Context, db *gorm.DB) error {
	userService, err := users.New(db)
	if err != nil {
		return err
	}

	u := []*users.User{
		{
			Email:    gofakeit.Email(),
			UserType: users.UserTypeAdmin,
			Name:     gofakeit.Name(),
			Password: "p@assword",
		},
		{
			Email:    gofakeit.Email(),
			UserType: users.UserTypeRegular,
			Name:     gofakeit.Name(),
			Password: "p@assword",
		},
		{
			Email:    gofakeit.Email(),
			UserType: users.UserTypeRegular,
			Name:     gofakeit.Name(),
			Password: "password",
		},
		{
			Email:    gofakeit.Email(),
			UserType: users.UserTypeRegular,
			Name:     gofakeit.Name(),
			Password: "password",
		},
	}

	for i := range u {
		if err := userService.Create(ctx, u[i]); err != nil {
			return err
		}
	}

	return nil
}

func setupDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
