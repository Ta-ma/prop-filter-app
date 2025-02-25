package models

type Property struct {
	ID            uint
	SquareFootage float32
	Price         float32
	Rooms         uint
	Bathrooms     uint
	Latitude      float64
	Longitude     float64
	Description   string
	Lighting      Lighting
	LightingID    uint
	Ammenities    []Ammenity `gorm:"many2many:properties_ammenities;"`
}
