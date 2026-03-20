package config

import "fmt"

type Application struct {
	Host            string `env:"APP_HOST,required" envDefault:"localhost"`
	Port            int    `env:"APP_PORT,required" envDefault:"8080"`
	Name            string `env:"APP_NAME" envDefault:"Created By Aldiwildan"`
	TZ              string `env:"APP_TIMEZONE" envDefault:"Asia/Jakarta"`
	Timeout         int    `env:"APP_TIMEOUT" envDefault:"10"`
	Debug           bool   `env:"APP_DEBUG" envDefault:"false"`
	ShutdownTimeout int    `env:"APP_SHUTDOWN_TIMEOUT" envDefault:"10"`
}

func (a Application) GetAddress() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}
