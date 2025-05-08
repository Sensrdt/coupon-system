package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Sensrdt/coupon-system/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
	mu sync.RWMutex
}

// ValidateTables checks if required tables exist and creates them if they don't
func (db *DB) ValidateTables() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check and create coupons table
	if !db.Migrator().HasTable(&model.Coupon{}) {
		log.Println("Creating coupons table...")
		if err := db.AutoMigrate(&model.Coupon{}); err != nil {
			return fmt.Errorf("failed to create coupons table: %v", err)
		}
		log.Println("Coupons table created successfully")
	} else {
		log.Println("Coupons table already exists")
	}

	if !db.Migrator().HasTable(&model.Cart{}) {
		log.Println("Creating cart table...")
		if err := db.AutoMigrate(&model.Cart{}); err != nil {
			return fmt.Errorf("failed to create cart table: %v", err)
		}
		log.Println("Cart table created successfully")
	} else {
		log.Println("Cart table already exists")
	}

	return nil
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

	dbInstance := &DB{DB: db}

	// Validate and create tables if needed
	if err := dbInstance.ValidateTables(); err != nil {
		log.Fatalf("Failed to validate/create tables: %v", err)
	}

	return dbInstance
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

// CreateCoupon creates a new coupon within a transaction
func (db *DB) CreateCoupon(ctx context.Context, c *model.Coupon) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existingCoupon model.Coupon
		if err := tx.Where("code = ?", c.Code).First(&existingCoupon).Error; err == nil {
			return fmt.Errorf("coupon code already exists")
		}

		if err := tx.Create(c).Error; err != nil {
			return fmt.Errorf("failed to create coupon: %v", err)
		}

		return nil
	})
}

func (db *DB) GetAllCoupons(ctx context.Context) []model.Coupon {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var coupons []model.Coupon
	db.WithContext(ctx).Find(&coupons)
	return coupons
}

func (db *DB) FindCouponByCode(ctx context.Context, code string) (model.Coupon, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var coupon model.Coupon
	result := db.WithContext(ctx).Where("code = ?", code).First(&coupon)
	return coupon, result.Error == nil
}
