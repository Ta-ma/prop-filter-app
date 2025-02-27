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

	maxAmmenities := len(models.GetAmmenityValues())
	ammenityIds := make([]int, maxAmmenities)
	for j := 0; j < maxAmmenities; j++ {
		ammenityIds[j] = j + 1
	}

	for i := uint(0); i < amount; i++ {
		addr = gofakeit.Address()

		// Shuffle ammenities list and pick a random amount of them
		gofakeit.ShuffleInts(ammenityIds)
		ammenities := make([]models.Ammenity, 0)
		ammenitiesCount := gofakeit.IntRange(0, maxAmmenities)
		for j := 0; j < ammenitiesCount; j++ {
			ammenities = append(ammenities, models.Ammenity{ID: uint(ammenityIds[j])})
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
			Ammenities:    ammenities,
		})
	}

	return properties
}
