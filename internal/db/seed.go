package db

import (
	"github.com/ta-ma/prop-filter-app/internal/datagen"
	"github.com/ta-ma/prop-filter-app/internal/models"
	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB, entries uint) {
	// Migrate ammenities
	deleteTable("properties_ammenities")
	migrateTable(&models.Ammenity{})

	ammenities := make([]models.Ammenity, 0)
	for _, a := range models.GetAmmenityValues() {
		ammenities = append(ammenities, models.Ammenity{Description: a})
	}
	db.Create(&ammenities)

	// Migrate lightings
	migrateTable(&models.Lighting{})

	lightings := make([]models.Lighting, 0)
	for _, l := range models.GetLightingValues() {
		lightings = append(lightings, models.Lighting{Description: l})
	}
	db.Create(&lightings)

	// Migrate properties
	migrateTable(&models.Property{})

	props := datagen.GenerateMockProperties(entries)
	// Batch insert in slices of 1000 elements due to Postgres restrictions
	batches := (entries / 1000)
	if batches == 0 {
		batches = 1
	}

	for i := uint(0); i < batches; i++ {
		lower := 1000 * i
		upper := 1000 * (i + 1)

		if upper > entries {
			upper = entries
		}
		propsSlice := props[lower:upper]
		db.Create(&propsSlice)
	}
}

func migrateTable[T any](model *T) {
	if db.Migrator().HasTable(model) {
		db.Migrator().DropTable(model)
	}

	db.AutoMigrate(model)
}

func deleteTable(tableName string) {
	if db.Migrator().HasTable(tableName) {
		db.Migrator().DropTable(tableName)
	}
}
