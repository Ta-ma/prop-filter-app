package models

type Property struct {
	ID            uint
	SquareFootage float32
	Lighting      string
	Price         float32
	Rooms         uint
	Bathrooms     uint
	LocationX     float64
	LocationY     float64
	Description   string
	Ammenities    []Ammenity `gorm:"many2many:properties_ammenities;"`
}
