package service

import (
	"context"
	"testing"
	"time"

	"github.com/Sensrdt/coupon-system/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

// MockCache is a mock implementation of the Cache interface
type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(key string) (interface{}, bool) {
	args := m.Called(key)
	return args.Get(0), args.Bool(1)
}

func (m *MockCache) Set(key string, value interface{}) {
	m.Called(key, value)
}

func (m *MockCache) Delete(key string) {
	m.Called(key)
}

func (m *MockRepository) CreateCoupon(ctx context.Context, coupon *model.Coupon) error {
	args := m.Called(ctx, coupon)
	return args.Error(0)
}

func (m *MockRepository) GetAllCoupons(ctx context.Context) ([]*model.Coupon, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Coupon), args.Error(1)
}

func (m *MockRepository) FindCouponByCode(ctx context.Context, code string) (*model.Coupon, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Coupon), args.Error(1)
}

func (m *MockRepository) UpdateCoupon(ctx context.Context, coupon *model.Coupon) error {
	args := m.Called(ctx, coupon)
	return args.Error(0)
}

func setupTestService(t *testing.T) (*CouponService, *MockRepository, *MockCache) {
	mockRepo := new(MockRepository)
	mockCache := new(MockCache)
	service := NewCouponService(mockRepo, mockCache)
	return service, mockRepo, mockCache
}

func TestGetApplicableCoupons(t *testing.T) {
	service, mockRepo, mockCache := setupTestService(t)
	ctx := context.Background()

	// Mock data
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
	}

	// Setup expectations
	mockRepo.On("GetAllCoupons", ctx).Return(coupons, nil)
	mockCache.On("Get", mock.Anything).Return(nil, false)
	mockCache.On("Set", mock.Anything, mock.Anything).Return()

	// Test data
	cart := &model.Cart{
		Items: []model.CartItem{
			{ID: "item1", Price: 150},
		},
		Total: 150,
	}

	// Execute test
	applicableCoupons, err := service.GetApplicableCoupons(ctx, cart)
	assert.NoError(t, err)
	assert.Len(t, applicableCoupons, 1)
	assert.Equal(t, "TEST10", applicableCoupons[0].Code)

	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestValidateCoupon(t *testing.T) {
	service, mockRepo, mockCache := setupTestService(t)
	ctx := context.Background()

	// Mock data
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

	// Setup expectations
	mockRepo.On("FindCouponByCode", ctx, "TEST10").Return(coupon, nil)
	mockRepo.On("UpdateCoupon", ctx, mock.AnythingOfType("*model.Coupon")).Return(nil)
	mockCache.On("Get", mock.Anything).Return(nil, false)
	mockCache.On("Set", mock.Anything, true).Return()

	// Test data
	cart := &model.Cart{
		Items: []model.CartItem{
			{ID: "item1", Price: 150},
		},
		Total: 150,
	}

	// Execute test
	valid, err := service.ValidateCoupon(ctx, "TEST10", cart)
	assert.NoError(t, err)
	assert.True(t, valid)

	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestCreateCoupon(t *testing.T) {
	service, mockRepo, mockCache := setupTestService(t)
	ctx := context.Background()

	// Test data
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

	// Setup expectations
	mockRepo.On("CreateCoupon", ctx, coupon).Return(nil)
	mockCache.On("Delete", mock.Anything).Return()

	// Execute test
	err := service.CreateCoupon(ctx, coupon)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestConcurrentOperations(t *testing.T) {
	service, mockRepo, mockCache := setupTestService(t)
	ctx := context.Background()

	// Mock data
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

	// Setup expectations
	mockRepo.On("FindCouponByCode", ctx, "TEST10").Return(coupon, nil).Times(10)
	mockRepo.On("UpdateCoupon", ctx, mock.AnythingOfType("*model.Coupon")).Return(nil).Times(10)
	mockCache.On("Get", mock.Anything).Return(nil, false).Times(10)
	mockCache.On("Set", mock.Anything, true).Return().Times(10)

	// Test data
	cart := &model.Cart{
		Items: []model.CartItem{
			{ID: "item1", Price: 150},
		},
		Total: 150,
	}

	// Test concurrent operations
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			valid, err := service.ValidateCoupon(ctx, "TEST10", cart)
			assert.NoError(t, err)
			assert.True(t, valid)
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}
