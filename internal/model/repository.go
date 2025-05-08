package model

type Repository interface {
	CreateCoupon(c *Coupon) error
	GetAllCoupons() []Coupon
	FindCouponByCode(code string) (Coupon, bool)
}
