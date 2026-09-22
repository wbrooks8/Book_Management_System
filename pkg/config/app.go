package config

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Connect opens the database configured by DATABASE_DSN.
func Connect() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "kyle:changeMe12@tcp(127.0.0.1:3306)/simplerest?charset=utf8mb4&parseTime=True&loc=Local"
	}
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return d, nil
}

// Close releases the underlying SQL connection pool.
func Close(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
