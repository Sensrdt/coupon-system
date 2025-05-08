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

	// Drop and recreate tables to ensure schema is up to date
	if db.Migrator().HasTable(&model.Coupon{}) {
		log.Println("Dropping existing coupons table...")
		if err := db.Migrator().DropTable(&model.Coupon{}); err != nil {
			return fmt.Errorf("failed to drop tables: %v", err)
		}
	}

	log.Println("Creating coupons table...")
	if err := db.AutoMigrate(&model.Coupon{}); err != nil {
		return fmt.Errorf("failed to create coupons table: %v", err)
	}
	log.Println("Coupons table created successfully")

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

func (db *DB) GetAllCoupons(ctx context.Context) ([]*model.Coupon, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var coupons []*model.Coupon
	if err := db.WithContext(ctx).Find(&coupons).Error; err != nil {
		return nil, err
	}
	return coupons, nil
}

func (db *DB) FindCouponByCode(ctx context.Context, code string) (*model.Coupon, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var coupon model.Coupon
	if err := db.WithContext(ctx).Where("code = ?", code).First(&coupon).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &coupon, nil
}

func (db *DB) UpdateCoupon(ctx context.Context, coupon *model.Coupon) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.WithContext(ctx).Save(coupon).Error
}
