package models

type Lighting struct {
	ID          uint
	Description string
}

func GetLightingValues() []string {
	return []string{"low", "medium", "high"}
}
