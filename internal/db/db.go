package db

import (
	"fmt"

	"github.com/ta-ma/prop-filter-app/internal/config"
	"github.com/ta-ma/prop-filter-app/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Initialize(dbConfig *config.DbConfig) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dbConfig.Host, dbConfig.PgUser, dbConfig.PgPassword, dbConfig.DbName, dbConfig.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("ERROR: Could not connect to Postgres:", err)
		panic("Failed to connect to Postgres database!")
	}

	if dbConfig.SeedDatabase {
		SeedDatabase(db)
	}

	var properties []models.Property
	db.Preload("Ammenities").Find(&properties)

	for _, prop := range properties {
		ammenities := ""
		for _, a := range prop.Ammenities {
			ammenities += a.Description + ", "
		}
		fmt.Println(prop.ID, prop.Description, prop.Price, prop.SquareFootage, ammenities)
	}
}
