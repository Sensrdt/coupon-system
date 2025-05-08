package model

import (
	"context"
)

type Repository interface {
	// CreateCoupon creates a new coupon
	CreateCoupon(ctx context.Context, coupon *Coupon) error

	// GetAllCoupons returns all coupons
	GetAllCoupons(ctx context.Context) ([]*Coupon, error)

	FindCouponByCode(ctx context.Context, code string) (*Coupon, error)

	UpdateCoupon(ctx context.Context, coupon *Coupon) error
}
