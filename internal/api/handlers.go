package api

import (
	"net/http"

	"github.com/Sensrdt/coupon-system/internal/model"
	"github.com/Sensrdt/coupon-system/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.CouponService
}

func NewHandler(svc *service.CouponService) *Handler {
	return &Handler{svc: svc}
}

// @Summary Get applicable coupons
// @Description Get applicable coupons for a given cart and order total
// @Accept json
// @Produce json
// @Param cart_items body []model.Cart true "Cart items"
// @Param order_total int true "Order total"
// @Param time_stamp string true "Timestamp"
// @Success 200 {array} model.Coupon "Applicable coupons"
// @Failure 400 {object} gin.H "Bad Request"
func (h *Handler) GetApplicableCouponsHandler(c *gin.Context) {
	r := c.Request
	var request struct {
		CartItems  []model.Cart `json:"cart_items"`
		OrderTotal int32        `json:"order_total"`
		TimeStamp  string       `json:"time_stamp"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// add context
	ctx := r.Context()

	coupons := h.svc.GetApplicableCoupons(ctx, request.CartItems, request.OrderTotal, request.TimeStamp)
	c.JSON(http.StatusOK, coupons)
}

// @Summary Validate coupon
// @Description Validate a coupon for a given cart and order total
// @Accept json
// @Produce json
// @Param coupon_code body string true "Coupon code"
// @Param cart_items body []model.Cart true "Cart items"
// @Param order_total int true "Order total"
// @Param timestamp string true "Timestamp"
// @Success 200 {object} map[string]interface{} "Validation result"
// @Failure 400 {object} gin.H "Bad Request"
func (h *Handler) ValidateCouponHandler(c *gin.Context) {
	r := c.Request
	var req struct {
		CouponCode string       `json:"coupon_code"`
		CartItems  []model.Cart `json:"cart_items"`
		OrderTotal int32        `json:"order_total"`
		Timestamp  string       `json:"timestamp"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := r.Context()

	res := h.svc.ValidateCoupon(ctx, req.CouponCode, req.CartItems, req.OrderTotal, req.Timestamp)
	c.JSON(http.StatusOK, res)
}
