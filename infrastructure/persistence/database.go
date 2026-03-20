package persistence_infrastructure

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(dialector gorm.Dialector, isDebug bool) (*gorm.DB, error) {
	gormConf := new(gorm.Config)

	if isDebug {
		gormConf.Logger = logger.New(
			&log.Logger,
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Silent,
				Colorful:      true,
			},
		)
	}

	if dialector == nil {
		return nil, fmt.Errorf("unsupported database dialect, %s", dialector.Name())
	}

	instance, err := gorm.Open(dialector, gormConf)
	if err != nil {
		return nil, err
	}

	if isDebug {
		return instance.Debug(), nil
	}

	return instance, nil
}

func NewDialector(dsn string, dialect string) (gorm.Dialector, error) {
	switch dialect {
	case "mysql":
		return mysql.Open(dsn), nil
	case "postgres":
		return postgres.Open(dsn), nil
	default:
		return nil, fmt.Errorf("unsupported database dialect, %s", dialect)
	}
}
