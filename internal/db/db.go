package db

import (
	"fmt"

	"github.com/ta-ma/prop-filter-app/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

type QueryResult struct {
	Description    string
	Price          float32
	Square_footage float32
	Rooms          uint
	Bathrooms      uint
	Latitude       float64
	Longitude      float64
	Lighting       string
	Ammenities     string
}

func Initialize(dbConfig *config.DbConfig) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dbConfig.Host, dbConfig.PgUser, dbConfig.PgPassword, dbConfig.DbName, dbConfig.Port)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		fmt.Println("ERROR: Could not connect to Postgres:", err)
		panic("Failed to connect to Postgres database!")
	}

	if dbConfig.SeedDatabase {
		SeedDatabase(db, dbConfig.SeedEntries)
	}
}

func QueryProperties(queryFilter string, limit int, offset int) ([]QueryResult, error) {
	if db == nil {
		return []QueryResult{}, fmt.Errorf("database connection has not been initialized")
	}

	var queryResult []QueryResult
	err := db.
		Table("properties as p").
		Select("p.description, p.price, p.square_footage, p.rooms, p.bathrooms, p.latitude, p.longitude, l.description as lighting, STRING_AGG(a.description, ',') ammenities").
		Joins("join lightings l on p.lighting_id = l.id").
		Joins("left join properties_ammenities pa on p.id = pa.property_id").
		Joins("left join ammenities a on a.id = pa.ammenity_id").
		Group("p.id, l.description").
		Where(queryFilter).
		Limit(limit).
		Offset(offset).
		Scan(&queryResult).
		Error

	if err != nil {
		return []QueryResult{}, err
	}

	return queryResult, nil
}

func GetPropertiesCount(queryFilter string) (int, error) {
	var count int64
	err := db.
		Table("properties as p").
		Select("p.description, p.price, p.square_footage, p.rooms, p.bathrooms, p.latitude, p.longitude, l.description as lighting, STRING_AGG(a.description, ',')").
		Joins("join lightings l on p.lighting_id = l.id").
		Joins("left join properties_ammenities pa on p.id = pa.property_id").
		Joins("left join ammenities a on a.id = pa.ammenity_id").
		Group("p.id, l.description").
		Where(queryFilter).
		Count(&count).
		Error

	if err != nil {
		return int(count), err
	}
	return int(count), nil
}
