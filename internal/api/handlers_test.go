package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Sensrdt/coupon-system/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCouponService is a mock implementation of the CouponService interface
type MockCouponService struct {
	mock.Mock
}

func (m *MockCouponService) GetApplicableCoupons(ctx context.Context, cart *model.Cart) ([]*model.Coupon, error) {
	args := m.Called(ctx, cart)
	return args.Get(0).([]*model.Coupon), args.Error(1)
}

func (m *MockCouponService) ValidateCoupon(ctx context.Context, code string, cart *model.Cart) (bool, error) {
	args := m.Called(ctx, code, cart)
	return args.Bool(0), args.Error(1)
}

func (m *MockCouponService) CreateCoupon(ctx context.Context, coupon *model.Coupon) error {
	args := m.Called(ctx, coupon)
	return args.Error(0)
}

func setupTestRouter() (*gin.Engine, *MockCouponService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := new(MockCouponService)
	handler := NewHandler(mockService)

	// Apply middleware to routes
	router.POST("/applicable", requireJSON(), handler.GetApplicableCouponsHandler)
	router.POST("/validate", requireJSON(), handler.ValidateCouponHandler)
	router.POST("/", requireJSON(), handler.CreateCouponHandler)

	return router, mockService
}

func TestGetApplicableCouponsHandler(t *testing.T) {
	router, mockService := setupTestRouter()

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
	mockService.On("GetApplicableCoupons", mock.Anything, mock.AnythingOfType("*model.Cart")).Return(coupons, nil)

	// Test data
	request := GetApplicableCouponsRequest{
		Items: []model.CartItem{
			{ID: "item1", Price: 150},
		},
		Total: 150,
	}

	// Create request
	body, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/applicable", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
	var response []*model.Coupon
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, "TEST10", response[0].Code)

	mockService.AssertExpectations(t)
}

func TestValidateCouponHandler(t *testing.T) {
	router, mockService := setupTestRouter()

	// Setup expectations
	mockService.On("ValidateCoupon", mock.Anything, "TEST10", mock.AnythingOfType("*model.Cart")).Return(true, nil)

	// Test data
	request := ValidateCouponRequest{
		Code: "TEST10",
		Cart: model.Cart{
			Items: []model.CartItem{
				{ID: "item1", Price: 150},
			},
			Total: 150,
		},
	}

	// Create request
	body, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/validate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
	var response ValidateCouponResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Valid)

	mockService.AssertExpectations(t)
}

func TestCreateCouponHandler(t *testing.T) {
	router, mockService := setupTestRouter()

	// Test data
	request := CreateCouponRequest{
		Code:            "TEST10",
		DiscountType:    "percentage",
		DiscountValue:   10,
		MinOrderValue:   100,
		MaxDiscount:     50,
		StartDate:       time.Now(),
		EndDate:         time.Now().Add(24 * time.Hour),
		UsageLimit:      100,
		IsActive:        true,
		ApplicableItems: []string{"item1"},
	}

	// Setup expectations
	mockService.On("CreateCoupon", mock.Anything, mock.AnythingOfType("*model.Coupon")).Return(nil)

	// Create request
	body, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusCreated, w.Code)

	mockService.AssertExpectations(t)
}

func TestInvalidRequest(t *testing.T) {
	router, _ := setupTestRouter()

	// Test invalid JSON
	req, _ := http.NewRequest("POST", "/applicable", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMissingContentType(t *testing.T) {
	router, mockService := setupTestRouter()

	// Setup expectations - even though we expect a bad request due to content type,
	// the handler might still try to call the service
	mockService.On("GetApplicableCoupons", mock.Anything, mock.AnythingOfType("*model.Cart")).
		Return([]*model.Coupon{}, nil)

	// Test missing content type
	req, _ := http.NewRequest("POST", "/applicable", bytes.NewBufferString("{}"))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
