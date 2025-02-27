/*
Copyright © 2025 Santiago Tamashiro <santiago.tamashiro@gmail.com>
*/
package config

type Cli struct {
	TrimLength   int
	UseOldRender bool
}

type DbConfig struct {
	Host         string
	Port         uint
	PgUser       string
	PgPassword   string
	DbName       string
	SeedDatabase bool
	SeedEntries  uint
}

type Configuration struct {
	DbConfig DbConfig
	Cli      Cli
}
