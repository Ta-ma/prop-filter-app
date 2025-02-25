package models

type Ammenity struct {
	ID          uint
	Description string
}

func GetAmmenityValues() []string {
	return []string{"yard", "pool", "garage", "rooftop", "waterfront"}
}
