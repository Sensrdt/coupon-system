package db

import (
	"log"
	"os"

	"github.com/Sensrdt/coupon-system/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewDB() *DB {
	var err error

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "coupon.db"
	}

	db, err := gorm.Open(sqlite.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	return &DB{db}
}

func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func NewRepository(db *gorm.DB) model.Repository {
	return &DB{DB: db}
}

func (db *DB) CreateCoupon(c *model.Coupon) error {
	return db.Create(c).Error
}

func (db *DB) GetAllCoupons() []model.Coupon {
	var coupons []model.Coupon
	db.Find(&coupons)
	return coupons
}

func (db *DB) FindCouponByCode(code string) (model.Coupon, bool) {
	var coupon model.Coupon
	result := db.Where("code = ?", code).First(&coupon)
	return coupon, result.Error == nil
}
