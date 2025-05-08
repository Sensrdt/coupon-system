package model

import "context"

type Repository interface {
	CreateCoupon(ctx context.Context, c *Coupon) error
	GetAllCoupons(ctx context.Context) []Coupon
	FindCouponByCode(ctx context.Context, code string) (Coupon, bool)
}
