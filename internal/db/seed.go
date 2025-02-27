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

	// Create haversine function
	db.Exec(`create function fn_spheric_distance(x1 float, y1 float, x2 float, y2 float) returns float 
as
$$
declare 
	x1_radians float := x1 * PI() / 180; 
	y1_radians float := y1 * PI() / 180;
	x2_radians float := x2 * PI() / 180;
	y2_radians float := y2 * PI() / 180;
	earth_radius_miles float := 3958.939; 
	hav_theta float := (1 - COS(x1_radians - x2_radians)) / 2;
	hav_phi float := (1 - COS(y1_radians - y2_radians)) / 2;
	hav_alpha float := hav_theta + COS(x1_radians) * COS(x2_radians) * hav_phi;
begin
return 2 * earth_radius_miles * ASIN(SQRT(hav_alpha));
end;
$$
language plpgsql;`)
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
