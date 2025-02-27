/*
Copyright © 2025 Santiago Tamashiro <santiago.tamashiro@gmail.com>
*/
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
	Dist           float32
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

func QueryProperties(queryFilter string, limit int, offset int, calcDist bool, distX string, distY string) ([]QueryResult, error) {
	if db == nil {
		return []QueryResult{}, fmt.Errorf("database connection has not been initialized")
	}

	var queryResult []QueryResult
	var queryBuilder *gorm.DB
	if calcDist {
		queryBuilder = getDistanceQuery(queryFilter, distX, distY)
	} else {
		queryBuilder = getStandardQuery(queryFilter)
	}
	err := queryBuilder.Limit(limit).Offset(offset).Scan(&queryResult).Error

	if err != nil {
		return []QueryResult{}, err
	}

	return queryResult, nil
}

func GetPropertiesCount(queryFilter string, calcDist bool, distX string, distY string) (int, error) {
	var count int64
	var queryBuilder *gorm.DB
	if calcDist {
		queryBuilder = getDistanceQuery(queryFilter, distX, distY)
	} else {
		queryBuilder = getStandardQuery(queryFilter)
	}

	err := queryBuilder.Count(&count).Error

	if err != nil {
		return int(count), err
	}
	return int(count), nil
}

func getStandardQuery(queryFilter string) *gorm.DB {
	selectStatement :=
		"p.description, p.price, p.square_footage, p.rooms, p.bathrooms, p.latitude, p.longitude," +
			"l.description as lighting, STRING_AGG(a.description, ',') ammenities"

	return db.Table("properties as p").
		Select(selectStatement).
		Joins("join lightings l on p.lighting_id = l.id").
		Joins("left join properties_ammenities pa on p.id = pa.property_id").
		Joins("left join ammenities a on a.id = pa.ammenity_id").
		Group("p.id, l.description").
		Where(queryFilter)
}

func getDistanceQuery(queryFilter string, distX string, distY string) *gorm.DB {
	selectStatement :=
		"p.description, p.price, p.square_footage, p.rooms, p.bathrooms, p.latitude, p.longitude," +
			"l.description as lighting, STRING_AGG(a.description, ',') ammenities, d.dist"

	distStatement :=
		fmt.Sprintf("join (select id, fn_spheric_distance(%s, %s, latitude, longitude) as dist from properties) d on p.id = d.id", distX, distY)

	return db.Table("properties as p").
		Select(selectStatement).
		Joins("join lightings l on p.lighting_id = l.id").
		Joins(distStatement).
		Joins("left join properties_ammenities pa on p.id = pa.property_id").
		Joins("left join ammenities a on a.id = pa.ammenity_id").
		Group("p.id, l.description, d.dist").
		Where(queryFilter)
}
