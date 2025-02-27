/*
Copyright © 2025 Santiago Tamashiro <santiago.tamashiro@gmail.com>
*/
package models

type Lighting struct {
	ID          uint
	Description string
}

func GetLightingValues() []string {
	return []string{"low", "medium", "high"}
}
