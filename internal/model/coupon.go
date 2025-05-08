package model

type Coupon struct {
	ID                    uint    `json:"-" gorm:"primaryKey"`
	Code                  string  `json:"coupon_code" gorm:"uniqueIndex"`
	ExpiryDate            string  `json:"expiry_date"`
	UsageType             string  `json:"usage_type"`
	ApplicableMedicineIDs string  `json:"applicable_medicine_ids"`
	ApplicableCategories  string  `json:"applicable_categories"`
	MinOrderValue         float64 `json:"min_order_value"`
	DiscountType          string  `json:"discount_type"`
	DiscountValue         float64 `json:"discount_value"`
	MaxUsagePerUser       int     `json:"max_usage_per_user"`
	TermsAndConditions    string  `json:"terms_and_conditions"`
}

type CreateCouponRequest struct {
	CouponCode            string `json:"coupon_code"`
	ExpiryDate            string `json:"expiry_date"`
	UsageType             string `json:"usage_type"`
	ApplicableMedicineIDs string `json:"applicable_medicine_ids"`
	ApplicableCategories  string `json:"applicable_categories"`
	MinOrderValue         int32  `json:"min_order_value"`
	ValidTimeWindow       string `json:"valid_time_window"`
	TermsAndConditions    string `json:"terms_and_conditions"`
	DiscountType          string `json:"discount_type"`
	DiscountValue         int32  `json:"discount_value"`
	MaxUsagePerUser       int32  `json:"max_usage_per_user"`
}
