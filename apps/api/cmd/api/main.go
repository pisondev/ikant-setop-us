package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/pisondev/ikant-setop-us/apps/api/internal/config"
	"github.com/pisondev/ikant-setop-us/apps/api/internal/database"
	"github.com/pisondev/ikant-setop-us/apps/api/internal/modules/dashboard"
	"github.com/pisondev/ikant-setop-us/apps/api/internal/modules/fish"
	"github.com/pisondev/ikant-setop-us/apps/api/internal/modules/stock"
	"github.com/pisondev/ikant-setop-us/apps/api/internal/modules/stockout"
	"github.com/pisondev/ikant-setop-us/apps/api/internal/modules/storage"
	"github.com/pisondev/ikant-setop-us/apps/api/internal/shared"
)

func main() {
	cfg := config.Load()
	log := shared.NewLogger(cfg.AppEnv)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.NewPool(ctx, cfg.DatabaseURL())
	if err != nil {
		log.WithError(err).Fatal("failed to connect database")
	}
	defer db.Close()

	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName,
		ErrorHandler: shared.ErrorHandler(log),
	})

	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSAllowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return shared.Success(c, fiber.StatusOK, "API is running", fiber.Map{
			"service": cfg.AppName,
			"version": cfg.AppVersion,
		})
	})

	v1 := app.Group("/api/v1")
	v1.Get("/health", func(c *fiber.Ctx) error {
		return shared.Success(c, fiber.StatusOK, "API is running", fiber.Map{
			"service": cfg.AppName,
			"version": cfg.AppVersion,
		})
	})

	fishRepo := fish.NewRepository(db)
	fish.NewHandler(fishRepo).RegisterRoutes(v1)

	storageRepo := storage.NewRepository(db)
	storage.NewHandler(storageRepo).RegisterRoutes(v1)

	stockRepo := stock.NewRepository(db)
	stockService := stock.NewService(stockRepo)
	stock.NewHandler(stockRepo, stockService).RegisterRoutes(v1)

	stockoutRepo := stockout.NewRepository(db)
	stockoutService := stockout.NewService(stockoutRepo)
	stockout.NewHandler(stockoutRepo, stockoutService).RegisterRoutes(v1)

	dashboardRepo := dashboard.NewRepository(db)
	dashboard.NewHandler(dashboardRepo).RegisterRoutes(v1)

	app.Use(func(c *fiber.Ctx) error {
		return shared.Error(c, fiber.StatusNotFound, "Route not found", nil)
	})

	go func() {
		if err := app.Listen(":" + cfg.AppPort); err != nil {
			log.WithError(err).Fatal("server stopped unexpectedly")
		}
	}()

	log.WithField("port", cfg.AppPort).Info("api server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.WithError(err).Error("failed to shutdown server cleanly")
	}
}
