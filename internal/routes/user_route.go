package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/thanadonexe/go-shorter/internal/adapters/handlers"
	"github.com/thanadonexe/go-shorter/internal/adapters/repositories"
	"github.com/thanadonexe/go-shorter/internal/core/ports"
	"github.com/thanadonexe/go-shorter/internal/core/services"
	"github.com/thanadonexe/go-shorter/internal/middlewares"
)

func SetupUserRoute(app *fiber.App, db *sqlx.DB, cache ports.CacheRepository) {
	repo := repositories.NewUserRepository(db, cache)
	service := services.NewUserService(repo)
	handler := handlers.NewUserHandler(service)

	r := app.Group("/v1")
	r.Get("/users/:id", middlewares.AuthMiddleware, handler.Get)
	r.Get("/users", middlewares.AuthMiddleware, handler.GetAll)
	r.Post("/users", handler.Create)
	r.Put("/users/:id", middlewares.AuthMiddleware, handler.Update)
	r.Delete("/users/:id", middlewares.AuthMiddleware, handler.Delete)
}
