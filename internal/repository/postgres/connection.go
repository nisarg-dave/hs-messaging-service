package postgres

import (
	"fmt"

	"hs-messaging-service/internal/config"
	"hs-messaging-service/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(config *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// {} means creates an empty/zero value of the struct
	err = db.AutoMigrate(&domain.Message{})
	if err != nil {
		return nil, err
	}

	fmt.Println("Database connected successfully")

	return db, nil
}