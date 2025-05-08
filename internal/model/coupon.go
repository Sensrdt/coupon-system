package model

import (
	"time"
)

// Coupon represents a discount coupon
type Coupon struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Code            string    `json:"code" gorm:"uniqueIndex"`
	DiscountType    string    `json:"discount_type"`
	DiscountValue   float64   `json:"discount_value"`
	MinOrderValue   float64   `json:"min_order_value"`
	MaxDiscount     float64   `json:"max_discount"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	UsageLimit      int       `json:"usage_limit"`
	UsageCount      int       `json:"usage_count"`
	IsActive        bool      `json:"is_active"`
	ApplicableItems []string  `json:"applicable_items" gorm:"type:text;serializer:json"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Cart represents a shopping cart
type Cart struct {
	Items []CartItem `json:"items"`
	Total float64    `json:"total"`
}

// CartItem represents an item in the cart
type CartItem struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
}

// CreateCouponRequest represents the request for creating a coupon
type CreateCouponRequest struct {
	Code            string    `json:"code"`
	DiscountType    string    `json:"discount_type"`
	DiscountValue   float64   `json:"discount_value"`
	MinOrderValue   float64   `json:"min_order_value"`
	MaxDiscount     float64   `json:"max_discount"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	UsageLimit      int       `json:"usage_limit"`
	IsActive        bool      `json:"is_active"`
	ApplicableItems []string  `json:"applicable_items"`
}
