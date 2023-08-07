package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/thanadonexe/go-shorter/internal/adapters/caches"
	"github.com/thanadonexe/go-shorter/internal/config"
	"github.com/thanadonexe/go-shorter/internal/routes"
	"github.com/thanadonexe/go-shorter/migrations"
)

func init() {
	if err := config.InitAppConfig(); err != nil {
		panic(fmt.Sprintf("Failed to init app config - Error: %v", err))
	}
}

func setupRoutes(app *fiber.App, db *sqlx.DB, cache *caches.RedisCache) {
	routes.SetupAuthRoute(app, db, cache)
	routes.SetupUserRoute(app, db, cache)
	routes.SetupUrlRoute(app, db, cache)
}

func main() {
	app := fiber.New()

	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.AppConfig.DBHost, config.AppConfig.DBPort, config.AppConfig.DBUsername, config.AppConfig.DBPassword, config.AppConfig.DBName,
	))
	if err != nil {
		panic(fmt.Sprintf("Failed to connect DB - Error: %v", err))
	}

	cache, err := caches.NewRedisCache(config.AppConfig.REDISHost, config.AppConfig.REDISPort, config.AppConfig.REDISPassword, config.AppConfig.REDISDb)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect Redis - Error: %v", err))
	}

	migrations.Run(db)
	setupRoutes(app, db, cache)

	go func() {
		if err := app.Listen(fmt.Sprintf(":%d", config.AppConfig.Port)); err != nil {
			panic(fmt.Sprintf("Failed to start the api server - Error: %v", err))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	fmt.Println("Gracefully shutdown...")
	_ = app.Shutdown()

	fmt.Println("Running clean up tasks...")
	// clean up tasks
	db.Close()
	cache.Close()
}
