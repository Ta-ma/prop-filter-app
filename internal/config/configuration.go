package config

type Cli struct {
	TrimLength int
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
