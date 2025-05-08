package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Sensrdt/coupon-system/internal/cache"
	"github.com/Sensrdt/coupon-system/internal/model"
)

type CouponService struct {
	repo  model.Repository
	cache cache.Cache
}

func NewCouponService(repo model.Repository, cache cache.Cache) *CouponService {
	return &CouponService{repo: repo, cache: cache}
}

func (s *CouponService) GetApplicableCoupons(ctx context.Context, cart []model.Cart, orderTotal int32, timeStamp string) []model.Coupon {
	// add lru cache
	cacheKey := fmt.Sprintf("applicable_coupons_%v_%v_%v", cart, orderTotal, timeStamp)
	coupons, found := s.cache.Get(cacheKey)
	if found {
		return coupons.([]model.Coupon)
	}
	all := s.repo.GetAllCoupons()
	var applicable []model.Coupon
	for _, c := range all {
		if s.meetsCriteria(c, cart, orderTotal, timeStamp) {
			applicable = append(applicable, c)
		}
	}
	s.cache.Set(cacheKey, applicable)
	return applicable
}

func (s *CouponService) ValidateCoupon(ctx context.Context, code string, cart []model.Cart, total int32, ts string) map[string]interface{} {
	cacheKey := fmt.Sprintf("validate_coupon_%v_%v_%v_%v", code, cart, total, ts)
	coupon, found := s.cache.Get(cacheKey)
	if found {
		return coupon.(map[string]interface{})
	}
	couponFromDB, foundinDB := s.repo.FindCouponByCode(code)
	if !foundinDB {
		return map[string]interface{}{"is_valid": false, "reason": "coupon not found"}
	}
	if !s.meetsCriteria(couponFromDB, cart, total, ts) {
		return map[string]interface{}{"is_valid": false, "reason": "coupon expired or not applicable"}
	}

	// dummy
	discount := map[string]interface{}{
		"items_discount":   50,
		"charges_discount": 20,
	}

	s.cache.Set(cacheKey, discount)
	return map[string]interface{}{"is_valid": true, "discount": discount, "message": "coupon applied successfully"}
}

func (s *CouponService) meetsCriteria(c model.Coupon, cart []model.Cart, total int32, ts string) bool {
	t, _ := time.Parse(time.RFC3339, ts)
	expiry, _ := time.Parse(time.RFC3339, c.ExpiryDate)
	if t.After(expiry) || total < int32(c.MinOrderValue) {
		return false
	}
	return true
}
