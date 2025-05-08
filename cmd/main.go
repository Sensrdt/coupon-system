package main

import (
	"os"

	_ "github.com/Sensrdt/coupon-system/docs/swagger" // swagger docs
	"github.com/Sensrdt/coupon-system/internal/api"
	"github.com/Sensrdt/coupon-system/internal/cache"
	"github.com/Sensrdt/coupon-system/internal/db"
	"github.com/Sensrdt/coupon-system/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files" // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Coupon System API
// @version 1.0
// @description This is a coupon system API server.
// @host localhost:8080
// @BasePath /
func main() {
	cfg := os.Getenv("PORT")
	if cfg == "" {
		cfg = "8080"
	}
	dbConn := db.NewDB()
	repo := db.NewRepository(dbConn.DB)
	cache := cache.NewLRU(100)
	couponService := service.NewCouponService(repo, cache)
	apiHandler := api.NewHandler(couponService)
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router := r.Group("/coupons")
	{
		router.POST("/applicable", apiHandler.GetApplicableCouponsHandler)
		router.POST("/validate", apiHandler.ValidateCouponHandler)
		router.POST("/", apiHandler.CreateCouponHandler)
	}

	r.Run(":" + cfg)
}
