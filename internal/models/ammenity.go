/*
Copyright © 2025 Santiago Tamashiro <santiago.tamashiro@gmail.com>
*/
package models

type Amenity struct {
	ID          uint
	Description string
}

func GetAmenityValues() []string {
	return []string{"yard", "pool", "garage", "rooftop", "waterfront"}
}
