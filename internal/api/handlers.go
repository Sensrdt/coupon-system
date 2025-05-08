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
// @Param request body GetApplicableCouponsRequest true "Request body"
// @Success 200 {array} model.Coupon "Applicable coupons"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Router /coupons/applicable [post]
func (h *Handler) GetApplicableCouponsHandler(c *gin.Context) {
	var request GetApplicableCouponsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	coupons := h.svc.GetApplicableCoupons(c.Request.Context(), request.CartItems, request.OrderTotal, request.TimeStamp)
	c.JSON(http.StatusOK, coupons)
}

// @Summary Validate coupon
// @Description Validate a coupon for a given cart and order total
// @Accept json
// @Produce json
// @Param request body ValidateCouponRequest true "Request body"
// @Success 200 {object} ValidateCouponResponse "Validation result"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Router /coupons/validate [post]
func (h *Handler) ValidateCouponHandler(c *gin.Context) {
	var req ValidateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	res := h.svc.ValidateCoupon(c.Request.Context(), req.CouponCode, req.CartItems, req.OrderTotal, req.Timestamp)
	c.JSON(http.StatusOK, res)
}

// @Summary Create coupon
// @Description Create a new coupon
// @Accept json
// @Produce json
// @Param request body model.CreateCouponRequest true "Request body"
// @Success 200 {object} CreateCouponResponse "Coupon created"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Router /coupons/create [post]
func (h *Handler) CreateCouponHandler(c *gin.Context) {
	var req model.CreateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.svc.CreateCoupon(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, CreateCouponResponse{Message: "Coupon created successfully"})
}

type GetApplicableCouponsRequest struct {
	CartItems  []model.Cart `json:"cart_items"`
	OrderTotal int32        `json:"order_total"`
	TimeStamp  string       `json:"time_stamp"`
}

type ValidateCouponRequest struct {
	CouponCode string       `json:"coupon_code"`
	CartItems  []model.Cart `json:"cart_items"`
	OrderTotal int32        `json:"order_total"`
	Timestamp  string       `json:"timestamp"`
}
type ValidateCouponResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message,omitempty"`
}

type CreateCouponResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
