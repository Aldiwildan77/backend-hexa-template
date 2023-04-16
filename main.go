package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Aldiwildan77/backend-hexa-template/config"
	"github.com/Aldiwildan77/backend-hexa-template/core/module"
	report_handler "github.com/Aldiwildan77/backend-hexa-template/handler/http"
	logger_infrastructure "github.com/Aldiwildan77/backend-hexa-template/infrastructure/logger"
	middleware_infrastructure "github.com/Aldiwildan77/backend-hexa-template/infrastructure/middleware"
	persistence_infrastructure "github.com/Aldiwildan77/backend-hexa-template/infrastructure/persistence"
	report_repository "github.com/Aldiwildan77/backend-hexa-template/repository/report"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal().Msgf("error loading config: %v", err)
	}

	log.Logger = logger_infrastructure.NewLogger(*cfg)

	rls := middleware.NewRateLimiterMemoryStoreWithConfig(
		middleware.RateLimiterMemoryStoreConfig{
			Rate:      rate.Limit(cfg.Middleware.RateLimit.Rate),
			Burst:     cfg.Middleware.RateLimit.Burst,
			ExpiresIn: cfg.Middleware.RateLimit.GetExpireTime(),
		},
	)

	rl := middleware_infrastructure.NewRateLimiter(
		middleware_infrastructure.WithRateLimiterStore(rls),
	)

	cors := middleware_infrastructure.NewCORS(
		middleware_infrastructure.WithAllowOrigins(cfg.Middleware.CORS.AllowOrigins),
		middleware_infrastructure.WithAllowMethods(cfg.Middleware.CORS.AllowMethods),
	)

	to := middleware_infrastructure.NewTimeout(
		middleware_infrastructure.WithTimeout(cfg.Middleware.Timeout.Duration),
	)

	stats := middleware_infrastructure.NewStatistic()

	rd := persistence_infrastructure.NewRedis(redis.Options{
		Addr:         cfg.Cache.GetAddress(),
		ClientName:   cfg.Application.Name,
		Username:     cfg.Cache.Username,
		Password:     cfg.Cache.Password,
		MaxIdleConns: cfg.Cache.MaxIdle,
		MaxRetries:   cfg.Cache.MaxRetry,
	})
	if err := rd.Ping(context.Background()).Err(); err != nil {
		log.Fatal().Msgf("error connecting to redis: %v", err)
	}

	db, err := persistence_infrastructure.NewDatabase(cfg.Database.GetDialector(), cfg.Database.Debug)
	if err != nil {
		log.Fatal().Msgf("error connecting to database: %v", err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Debug = cfg.Application.Debug

	e.Logger.SetOutput(log.Logger)
	e.Use(middleware.Secure())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(cors))
	e.Use(middleware.RequestID())
	e.Use(middleware.RateLimiterWithConfig(rl))
	e.Use(middleware.TimeoutWithConfig(to))
	e.Use(stats.Process)
	e.Use(middleware_infrastructure.ServerHeader)

	e.GET("/stats", stats.Handle)

	// REMOVE THIS
	reportRepo := report_repository.New(db)
	reportUsecase := module.NewReportUsecase(reportRepo)
	reportHandler := report_handler.NewReportHandler(reportUsecase)
	reportRouter := e.Group("/report")

	reportRouter.GET("/", reportHandler.GetReports)
	reportRouter.POST("/", reportHandler.CreateReport)
	reportRouter.GET("/:id", reportHandler.GetReport)
	reportRouter.PUT("/:id", reportHandler.UpdateReport)
	reportRouter.DELETE("/:id", reportHandler.DeleteReport)
	// END REMOVE THIS

	// Start server
	go func() {
		if err := e.Start(cfg.Application.GetAddress()); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("shutting down the server: %v", err)
		}
	}()

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Application.Timeout)*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal().Msgf("error shutting down server: %v", err)
	}

	log.Info().Msg("server gracefully stopped")
}
