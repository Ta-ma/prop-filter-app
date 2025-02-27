/*
Copyright © 2025 Santiago Tamashiro <santiago.tamashiro@gmail.com>
*/
package models

type PropertyViewModel struct {
	Description    string
	Price          float32
	Square_footage float32
	Rooms          uint
	Bathrooms      uint
	Latitude       float64
	Longitude      float64
	Lighting       string
	Amenities      string
	Dist           float32
}
