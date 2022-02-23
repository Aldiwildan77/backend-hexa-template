package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Aldiwildan77/backend-hexa-template/config"
	"github.com/Aldiwildan77/backend-hexa-template/core/module"
	report_handler "github.com/Aldiwildan77/backend-hexa-template/handler/http"
	report_repository "github.com/Aldiwildan77/backend-hexa-template/repository/report"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", cfg)

	db, err := NewDatabaseInstance(cfg)
	if err != nil {
		panic(err)
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	e := echo.New()
	e.Logger.SetLevel(log.LstdFlags)

	reportRepo := report_repository.New(db)
	reportUsecase := module.NewReportUsecase(reportRepo)
	reportHandler := report_handler.NewReportHandler(reportUsecase)
	reportRouter := e.Group("/report")

	reportRouter.GET("/", reportHandler.GetReports)
	reportRouter.POST("/", reportHandler.CreateReport)
	reportRouter.GET("/:id", reportHandler.GetReport)
	reportRouter.PUT("/:id", reportHandler.UpdateReport)
	reportRouter.DELETE("/:id", reportHandler.DeleteReport)

	// Start server
	go func() {
		if err := e.Start(addr); err != nil {
			log.Fatalf("error starting server: %v", err)
		}
	}()

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	e.Shutdown(ctx)
}

func NewDatabaseInstance(conf *config.Config) (*gorm.DB, error) {
	var dsn string
	gormConf := &gorm.Config{}

	if conf.DBDebug {
		gormConf.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Silent,
				Colorful:      true,
			},
		)
	}

	var instance *gorm.DB
	var err error

	switch conf.DBDialect {
	case "mysql":
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			conf.DBUsername,
			conf.DBPassword,
			conf.DBHost,
			conf.DBPort,
			conf.DBSchema,
		)
		instance, err = gorm.Open(mysql.Open(dsn), gormConf)
	case "postgres":
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
			conf.DBHost,
			conf.DBUsername,
			conf.DBPassword,
			conf.DBSchema,
			conf.DBPort,
		)
		instance, err = gorm.Open(postgres.Open(dsn), gormConf)
	default:
		err = fmt.Errorf("unsupported database dialect: %s", conf.DBDialect)
	}

	if err != nil {
		return nil, err
	}

	if conf.DBDebug {
		return instance.Debug(), nil
	}

	return instance, nil
}
