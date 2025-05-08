package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Sensrdt/coupon-system/internal/model"
	"github.com/gin-gonic/gin"
)

// CouponService defines the interface for coupon-related operations
type CouponService interface {
	GetApplicableCoupons(ctx context.Context, cart *model.Cart) ([]*model.Coupon, error)
	ValidateCoupon(ctx context.Context, code string, cart *model.Cart) (bool, error)
	CreateCoupon(ctx context.Context, coupon *model.Coupon) error
}

// Handler handles HTTP requests
type Handler struct {
	couponService CouponService
}

// NewHandler creates a new Handler
func NewHandler(couponService CouponService) *Handler {
	return &Handler{
		couponService: couponService,
	}
}

// GetApplicableCouponsRequest represents the request body for getting applicable coupons
type GetApplicableCouponsRequest struct {
	Items []model.CartItem `json:"items"`
	Total float64          `json:"total"`
}

// ValidateCouponRequest represents the request body for validating a coupon
type ValidateCouponRequest struct {
	Code string     `json:"code"`
	Cart model.Cart `json:"cart"`
}

// ValidateCouponResponse represents the response for coupon validation
type ValidateCouponResponse struct {
	Valid bool `json:"valid"`
}

// CreateCouponRequest represents the request body for creating a coupon
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

// GetApplicableCouponsHandler handles requests to get applicable coupons
// @Summary Get applicable coupons
// @Description Get coupons applicable to the given cart
// @Tags coupons
// @Accept json
// @Produce json
// @Param request body GetApplicableCouponsRequest true "Cart items and total"
// @Success 200 {array} model.Coupon
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /coupons/applicable [post]
func (h *Handler) GetApplicableCouponsHandler(c *gin.Context) {
	var req GetApplicableCouponsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request format"})
		return
	}

	cart := &model.Cart{
		Items: req.Items,
		Total: req.Total,
	}

	coupons, err := h.couponService.GetApplicableCoupons(c.Request.Context(), cart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get applicable coupons"})
		return
	}

	c.JSON(http.StatusOK, coupons)
}

// ValidateCouponHandler handles requests to validate a coupon
// @Summary Validate coupon
// @Description Validate a coupon against cart items
// @Tags coupons
// @Accept json
// @Produce json
// @Param request body ValidateCouponRequest true "Coupon code and cart"
// @Success 200 {object} ValidateCouponResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /coupons/validate [post]
func (h *Handler) ValidateCouponHandler(c *gin.Context) {
	var req ValidateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request format"})
		return
	}

	valid, err := h.couponService.ValidateCoupon(c.Request.Context(), req.Code, &req.Cart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to validate coupon"})
		return
	}

	c.JSON(http.StatusOK, ValidateCouponResponse{Valid: valid})
}

// CreateCouponHandler handles requests to create a coupon
// @Summary Create coupon
// @Description Create a new coupon
// @Tags coupons
// @Accept json
// @Produce json
// @Param request body CreateCouponRequest true "Coupon details"
// @Success 201 {object} model.Coupon
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /coupons/create [post]
func (h *Handler) CreateCouponHandler(c *gin.Context) {
	var req CreateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request format"})
		return
	}

	coupon := &model.Coupon{
		Code:            req.Code,
		DiscountType:    req.DiscountType,
		DiscountValue:   req.DiscountValue,
		MinOrderValue:   req.MinOrderValue,
		MaxDiscount:     req.MaxDiscount,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		UsageLimit:      req.UsageLimit,
		IsActive:        req.IsActive,
		ApplicableItems: req.ApplicableItems,
	}

	if err := h.couponService.CreateCoupon(c.Request.Context(), coupon); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create coupon"})
		return
	}

	c.JSON(http.StatusCreated, coupon)
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// requireJSON is a middleware that checks if the Content-Type header is set to application/json
func requireJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")
		if contentType != "application/json" {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: "Content-Type must be application/json"})
			return
		}
		c.Next()
	}
}
