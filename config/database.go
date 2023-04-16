package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseDialect string

const (
	DialectMySQL    DatabaseDialect = "mysql"
	DialectPostgres DatabaseDialect = "postgres"
)

type Database struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     int    `env:"DB_PORT" envDefault:"3306"`
	User     string `env:"DB_USER" envDefault:"root"`
	Password string `env:"DB_PASSWORD" envDefault:""`
	Name     string `env:"DB_NAME" envDefault:""`
	Charset  string `env:"DB_CHARSET" envDefault:"utf8mb4"`
	Debug    bool   `env:"DB_DEBUG" envDefault:"false"`
	Dialect  string `env:"DB_DIALECT" envDefault:"mysql"`
}

func (d Database) GetDialector() gorm.Dialector {
	switch d.Dialect {
	case string(DialectMySQL):
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			d.User,
			d.Password,
			d.Host,
			d.Port,
			d.Name,
		)
		return mysql.Open(dsn)
	case string(DialectPostgres):
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
			d.Host,
			d.User,
			d.Password,
			d.Name,
			d.Port,
		)
		return postgres.Open(dsn)
	default:
		return nil
	}
}
