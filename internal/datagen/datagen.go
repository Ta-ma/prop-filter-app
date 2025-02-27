/*
Copyright © 2025 Santiago Tamashiro <santiago.tamashiro@gmail.com>
*/
package datagen

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/ta-ma/prop-filter-app/internal/models"
)

func GenerateMockProperties(amount uint) []models.Property {
	var properties []models.Property = make([]models.Property, 0)
	var addr *gofakeit.AddressInfo

	maxAmenities := len(models.GetAmenityValues())
	amenityIds := make([]int, maxAmenities)
	for j := 0; j < maxAmenities; j++ {
		amenityIds[j] = j + 1
	}

	for i := uint(0); i < amount; i++ {
		addr = gofakeit.Address()

		// Shuffle amenities list and pick a random amount of them
		gofakeit.ShuffleInts(amenityIds)
		amenities := make([]models.Amenity, 0)
		amenitiesCount := gofakeit.IntRange(0, maxAmenities)
		for j := 0; j < amenitiesCount; j++ {
			amenities = append(amenities, models.Amenity{ID: uint(amenityIds[j])})
		}

		properties = append(properties, models.Property{
			ID:            i + 1,
			SquareFootage: gofakeit.Float32Range(100, 2000),
			Lighting:      models.Lighting{ID: uint(gofakeit.IntRange(1, len(models.GetLightingValues())))},
			Price:         gofakeit.Float32Range(10000, 999999),
			Rooms:         gofakeit.UintRange(1, 12),
			Bathrooms:     gofakeit.UintRange(1, 9),
			Latitude:      addr.Latitude,
			Longitude:     addr.Longitude,
			Description:   fmt.Sprintf("%s %s, %s", addr.Street, addr.City, addr.State),
			Amenities:     amenities,
		})
	}

	return properties
}
