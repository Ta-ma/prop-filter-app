package config

type DbConfig struct {
	Host         string
	Port         uint
	PgUser       string
	PgPassword   string
	DbName       string
	SeedDatabase bool
}

type Configuration struct {
	DbConfig DbConfig
}
