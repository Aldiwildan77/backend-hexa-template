package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Aldiwildan77/backend-hexa-template/config"
	"github.com/Aldiwildan77/backend-hexa-template/core/module"
	handlers "github.com/Aldiwildan77/backend-hexa-template/handler/http"
	logger_infrastructure "github.com/Aldiwildan77/backend-hexa-template/infrastructure/logger"
	middleware_infrastructure "github.com/Aldiwildan77/backend-hexa-template/infrastructure/middleware"
	persistence_infrastructure "github.com/Aldiwildan77/backend-hexa-template/infrastructure/persistence"
	report_repository "github.com/Aldiwildan77/backend-hexa-template/repository/report"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

func StartHTTPServer(cmd *cobra.Command, args []string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

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

	stats := middleware_infrastructure.NewStatistic()

	metricsMiddleware := middleware_infrastructure.NewMetrics()

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

	dialector, err := persistence_infrastructure.NewDialector(cfg.Database.DSN, cfg.Database.Dialect)
	if err != nil {
		log.Fatal().Msgf("error creating database dialector: %v", err)
	}

	db, err := persistence_infrastructure.NewDatabase(dialector, cfg.Database.Debug)
	if err != nil {
		log.Fatal().Msgf("error connecting to database: %v", err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Debug = cfg.Application.Debug

	e.Logger.SetOutput(log.Logger)
	e.Use(middleware.Secure())
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(cors))
	e.Use(middleware.RequestID())
	e.Use(middleware.RateLimiterWithConfig(rl))
	e.Use(stats.Process)
	e.Use(middleware_infrastructure.ServerHeader)

	if cfg.Middleware.Metrics.Enabled {
		e.Use(metricsMiddleware.Process)
	}

	e.GET("/stats", stats.Handle)

	if cfg.Middleware.Metrics.Enabled {
		e.GET("/metrics", metricsMiddleware.Handle)
	}

	// REMOVE THIS
	reportRepo := report_repository.New(db)
	reportUsecase := module.NewReportUsecase(reportRepo)
	reportHandler := handlers.NewReportHandler(reportUsecase)
	reportRouter := e.Group("/report")

	reportRouter.GET("/", reportHandler.GetReports)
	reportRouter.POST("/", reportHandler.CreateReport)
	reportRouter.GET("/:id", reportHandler.GetReport)
	reportRouter.PUT("/:id", reportHandler.UpdateReport)
	reportRouter.DELETE("/:id", reportHandler.DeleteReport)
	// END REMOVE THIS

	srv := &http.Server{
		Addr:    cfg.Application.GetAddress(),
		Handler: e,
	}

	// Start server
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("error starting server: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("Shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(ctx, time.Duration(cfg.Application.ShutdownTimeout)*time.Second)
		defer cancel()

		return srv.Shutdown(shutdownCtx)
	})

	log.Info().Msgf("Server is running on %s", cfg.Application.GetAddress())

	return eg.Wait()
}
