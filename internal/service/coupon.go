package service

import (
	"context"
	"sync"
	"time"

	"github.com/Sensrdt/coupon-system/internal/cache"
	"github.com/Sensrdt/coupon-system/internal/model"
)

type CouponService struct {
	repo  model.Repository
	cache cache.Cache
	mu    sync.RWMutex
}

func NewCouponService(repo model.Repository, cache cache.Cache) *CouponService {
	return &CouponService{
		repo:  repo,
		cache: cache,
	}
}

func (s *CouponService) GetApplicableCoupons(ctx context.Context, cart *model.Cart) ([]*model.Coupon, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cacheKey := generateCacheKey("applicable", cart)
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.([]*model.Coupon), nil
	}

	coupons, err := s.repo.GetAllCoupons(ctx)
	if err != nil {
		return nil, err
	}

	applicableCoupons := make([]*model.Coupon, 0)
	now := time.Now()

	for _, coupon := range coupons {
		if !coupon.IsActive {
			continue
		}

		if now.Before(coupon.StartDate) || now.After(coupon.EndDate) {
			continue
		}

		if cart.Total < coupon.MinOrderValue {
			continue
		}

		if coupon.UsageCount >= coupon.UsageLimit {
			continue
		}

		hasApplicableItem := false
		for _, item := range cart.Items {
			for _, applicableItem := range coupon.ApplicableItems {
				if item.ID == applicableItem {
					hasApplicableItem = true
					break
				}
			}
			if hasApplicableItem {
				break
			}
		}

		if hasApplicableItem {
			applicableCoupons = append(applicableCoupons, coupon)
		}
	}

	// Cache the result
	s.cache.Set(cacheKey, applicableCoupons)

	return applicableCoupons, nil
}

func (s *CouponService) ValidateCoupon(ctx context.Context, code string, cart *model.Cart) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cacheKey := generateCacheKey("validate", code, cart)
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.(bool), nil
	}

	coupon, err := s.repo.FindCouponByCode(ctx, code)
	if err != nil {
		return false, err
	}

	if coupon == nil {
		return false, nil
	}

	now := time.Now()
	if !coupon.IsActive {
		return false, nil
	}

	if now.Before(coupon.StartDate) || now.After(coupon.EndDate) {
		return false, nil
	}

	if cart.Total < coupon.MinOrderValue {
		return false, nil
	}

	if coupon.UsageCount >= coupon.UsageLimit {
		return false, nil
	}

	hasApplicableItem := false
	for _, item := range cart.Items {
		for _, applicableItem := range coupon.ApplicableItems {
			if item.ID == applicableItem {
				hasApplicableItem = true
				break
			}
		}
		if hasApplicableItem {
			break
		}
	}

	if !hasApplicableItem {
		return false, nil
	}

	coupon.UsageCount++
	if err := s.repo.UpdateCoupon(ctx, coupon); err != nil {
		return false, err
	}

	s.cache.Set(cacheKey, true)

	return true, nil
}

func (s *CouponService) CreateCoupon(ctx context.Context, coupon *model.Coupon) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if coupon.Code == "" {
		return ErrInvalidCouponCode
	}

	if coupon.DiscountValue <= 0 {
		return ErrInvalidDiscountValue
	}

	if coupon.MinOrderValue < 0 {
		return ErrInvalidMinOrderValue
	}

	if coupon.MaxDiscount < 0 {
		return ErrInvalidMaxDiscount
	}

	if coupon.UsageLimit <= 0 {
		return ErrInvalidUsageLimit
	}

	if coupon.StartDate.After(coupon.EndDate) {
		return ErrInvalidDateRange
	}

	// Create coupon
	if err := s.repo.CreateCoupon(ctx, coupon); err != nil {
		return err
	}

	// Invalidate cache
	s.cache.Delete(generateCacheKey("applicable", nil))

	return nil
}

// dummy
func generateCacheKey(prefix string, params ...interface{}) string {
	return prefix
}

// Error types
var (
	ErrInvalidCouponCode    = NewError("invalid coupon code")
	ErrInvalidDiscountValue = NewError("invalid discount value")
	ErrInvalidMinOrderValue = NewError("invalid minimum order value")
	ErrInvalidMaxDiscount   = NewError("invalid maximum discount")
	ErrInvalidUsageLimit    = NewError("invalid usage limit")
	ErrInvalidDateRange     = NewError("invalid date range")
)

// Error represents a service error
type Error struct {
	message string
}

// NewError creates a new Error
func NewError(message string) *Error {
	return &Error{message: message}
}

// Error returns the error message
func (e *Error) Error() string {
	return e.message
}
