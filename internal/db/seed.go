package db

import (
	"github.com/ta-ma/prop-filter-app/internal/models"
	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) {
	// Ammenities
	db.AutoMigrate(&models.Ammenity{})

	ammenities := []models.Ammenity{
		{Description: "pool"},
		{Description: "garage"},
		{Description: "yard"},
	}
	db.Create(&ammenities)

	// Properties
	db.AutoMigrate(&models.Property{})

	propertiesM := []models.Property{
		{
			SquareFootage: 500,
			Lighting:      "low",
			Price:         600000,
			Rooms:         6,
			Bathrooms:     2,
			LocationX:     150,
			LocationY:     250,
			Description:   "Ample place",
			Ammenities: []models.Ammenity{
				{Description: "yard"}, {Description: "pool"}, {Description: "garage"},
			},
		},
		{
			SquareFootage: 300,
			Lighting:      "high",
			Price:         450700,
			Rooms:         4,
			Bathrooms:     1,
			LocationX:     300,
			LocationY:     800,
			Description:   "Comfy",
			Ammenities: []models.Ammenity{
				{Description: "yard"},
			},
		},
		{
			SquareFootage: 200,
			Lighting:      "low",
			Price:         300000,
			Rooms:         3,
			Bathrooms:     1,
			LocationX:     65.9,
			LocationY:     75.7,
			Description:   "Haunted",
		},
		{
			SquareFootage: 675,
			Lighting:      "low",
			Price:         78050.5,
			Rooms:         3,
			Bathrooms:     1,
			LocationX:     500.2,
			LocationY:     600,
			Description:   "Nice place",
			Ammenities: []models.Ammenity{
				{Description: "pool"}, {Description: "garage"},
			},
		},
		{
			SquareFootage: 333,
			Lighting:      "low",
			Price:         190532.976,
			Rooms:         3,
			Bathrooms:     1,
			LocationX:     90,
			LocationY:     40,
			Description:   "Could be better",
			Ammenities: []models.Ammenity{
				{Description: "garage"}, {Description: "yard"},
			},
		},
	}
	db.Create(&propertiesM)
}
