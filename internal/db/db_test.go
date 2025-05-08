package db

import (
	"context"
	"testing"
	"time"

	"github.com/Sensrdt/coupon-system/internal/model"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *DB {
	db := NewDB()
	_ = NewRepository(db.DB)
	return db
}

func TestCreateCoupon(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	coupon := &model.Coupon{
		Code:            "TEST10",
		DiscountType:    "percentage",
		DiscountValue:   10,
		MinOrderValue:   100,
		MaxDiscount:     50,
		StartDate:       time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
		UsageLimit:      100,
		UsageCount:      0,
		IsActive:        true,
		ApplicableItems: []string{"item1", "item2"},
	}

	err := db.CreateCoupon(ctx, coupon)
	assert.NoError(t, err)
	assert.NotZero(t, coupon.ID)
}

func TestGetAllCoupons(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// Create test coupons
	coupons := []*model.Coupon{
		{
			Code:            "TEST10",
			DiscountType:    "percentage",
			DiscountValue:   10,
			MinOrderValue:   100,
			MaxDiscount:     50,
			StartDate:       time.Now(),
			EndDate:         time.Now().Add(24 * time.Hour),
			UsageLimit:      100,
			UsageCount:      0,
			IsActive:        true,
			ApplicableItems: []string{"item1"},
		},
		{
			Code:            "TEST20",
			DiscountType:    "percentage",
			DiscountValue:   20,
			MinOrderValue:   200,
			MaxDiscount:     100,
			StartDate:       time.Now(),
			EndDate:         time.Now().Add(24 * time.Hour),
			UsageLimit:      50,
			UsageCount:      0,
			IsActive:        true,
			ApplicableItems: []string{"item2"},
		},
	}

	for _, coupon := range coupons {
		err := db.CreateCoupon(ctx, coupon)
		assert.NoError(t, err)
	}

	// Test GetAllCoupons
	retrievedCoupons, err := db.GetAllCoupons(ctx)
	assert.NoError(t, err)
	assert.Len(t, retrievedCoupons, 2)
}

func TestFindCouponByCode(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// Create test coupon
	coupon := &model.Coupon{
		Code:            "TEST10",
		DiscountType:    "percentage",
		DiscountValue:   10,
		MinOrderValue:   100,
		MaxDiscount:     50,
		StartDate:       time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
		UsageLimit:      100,
		UsageCount:      0,
		IsActive:        true,
		ApplicableItems: []string{"item1"},
	}

	err := db.CreateCoupon(ctx, coupon)
	assert.NoError(t, err)

	// Test FindCouponByCode
	foundCoupon, err := db.FindCouponByCode(ctx, "TEST10")
	assert.NoError(t, err)
	assert.NotNil(t, foundCoupon)
	assert.Equal(t, coupon.Code, foundCoupon.Code)

	// Test non-existent coupon
	notFoundCoupon, err := db.FindCouponByCode(ctx, "NONEXISTENT")
	assert.NoError(t, err)
	assert.Nil(t, notFoundCoupon)
}

func TestUpdateCoupon(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// Create test coupon
	coupon := &model.Coupon{
		Code:            "TEST10",
		DiscountType:    "percentage",
		DiscountValue:   10,
		MinOrderValue:   100,
		MaxDiscount:     50,
		StartDate:       time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
		UsageLimit:      100,
		UsageCount:      0,
		IsActive:        true,
		ApplicableItems: []string{"item1"},
	}

	err := db.CreateCoupon(ctx, coupon)
	assert.NoError(t, err)

	// Update coupon
	coupon.DiscountValue = 20
	coupon.MaxDiscount = 100
	err = db.UpdateCoupon(ctx, coupon)
	assert.NoError(t, err)

	// Verify update
	updatedCoupon, err := db.FindCouponByCode(ctx, "TEST10")
	assert.NoError(t, err)
	assert.Equal(t, float64(20), updatedCoupon.DiscountValue)
	assert.Equal(t, float64(100), updatedCoupon.MaxDiscount)
}

func TestConcurrentOperations(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// Create test coupon
	coupon := &model.Coupon{
		Code:            "TEST10",
		DiscountType:    "percentage",
		DiscountValue:   10,
		MinOrderValue:   100,
		MaxDiscount:     50,
		StartDate:       time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
		UsageLimit:      100,
		UsageCount:      0,
		IsActive:        true,
		ApplicableItems: []string{"item1"},
	}

	err := db.CreateCoupon(ctx, coupon)
	assert.NoError(t, err)

	// Test concurrent reads
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			_, err := db.FindCouponByCode(ctx, "TEST10")
			assert.NoError(t, err)
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}
