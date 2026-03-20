package config

type DatabaseDialect string

const (
	DialectMySQL    DatabaseDialect = "mysql"
	DialectPostgres DatabaseDialect = "postgres"
)

type Database struct {
	DSN     string `env:"DB_DSN" envDefault:""`
	Debug   bool   `env:"DB_DEBUG" envDefault:"false"`
	Dialect string `env:"DB_DIALECT" envDefault:"mysql"`
}
