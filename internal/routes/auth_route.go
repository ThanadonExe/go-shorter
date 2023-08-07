package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/thanadonexe/go-shorter/internal/adapters/handlers"
	"github.com/thanadonexe/go-shorter/internal/adapters/repositories"
	"github.com/thanadonexe/go-shorter/internal/core/ports"
	"github.com/thanadonexe/go-shorter/internal/core/services"
)

func SetupAuthRoute(app *fiber.App, db *sqlx.DB, cache ports.CacheRepository) {
	repo := repositories.NewUserRepository(db, cache)
	service := services.NewAuthService(repo)
	handler := handlers.NewAuthHandler(service)

	r := app.Group("/v1")

	r.Post("/auth", handler.Login)
	r.Post("/refresh", handler.Refresh)
}
