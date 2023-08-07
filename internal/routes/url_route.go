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

func SetupUrlRoute(app *fiber.App, db *sqlx.DB, cache ports.CacheRepository) {
	repo := repositories.NewUrlRepository(db, cache)
	service := services.NewUrlService(repo)
	handler := handlers.NewUrlHandler(service)

	app.Get("/:code", handler.Resolve)

	v1 := app.Group("/v1")
	v1.Get("/urls/:id", handler.Get)
	v1.Get("/urls", middlewares.AuthMiddleware, handler.GetAll)
	v1.Post("/urls", middlewares.AuthMiddleware, handler.Create)
	v1.Delete("/urls/:id", middlewares.AuthMiddleware, handler.Delete)
}
