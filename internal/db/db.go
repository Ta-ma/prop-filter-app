package db

import (
	"fmt"

	"github.com/ta-ma/prop-filter-app/internal/config"
	"github.com/ta-ma/prop-filter-app/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Initialize(dbConfig *config.DbConfig) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dbConfig.Host, dbConfig.PgUser, dbConfig.PgPassword, dbConfig.DbName, dbConfig.Port)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("ERROR: Could not connect to Postgres:", err)
		panic("Failed to connect to Postgres database!")
	}

	if dbConfig.SeedDatabase {
		SeedDatabase(db)
	}
}

func QueryPropertiesWithAmmenities(selector string) ([]models.Property, error) {
	var properties []models.Property
	if db == nil {
		return properties, fmt.Errorf("database connection has not been initialized")
	}

	err := db.Where(selector).Preload("Ammenities").Find(&properties).Error
	if err != nil {
		return properties, err
	}
	return properties, nil
}
