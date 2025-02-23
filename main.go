/*
Copyright © 2025 Santiago Tamashiro <santiago.tamashiro@gmail.com>
*/
package main

import (
	"fmt"

	"github.com/ta-ma/prop-filter-app/cmd"
	"github.com/ta-ma/prop-filter-app/internal/config"
	"github.com/ta-ma/prop-filter-app/internal/db"
)

func main() {
	config, err := config.Read("config.json")
	if err != nil {
		fmt.Println("ERROR: Could not read JSON configuration file!")
		panic(err)
	}

	db.Initialize(&config.DbConfig)
	cmd.Execute()
}
