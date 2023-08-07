package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/thanadonexe/go-shorter/internal/config"
	"github.com/thanadonexe/go-shorter/internal/utils/token"
)

func AuthMiddleware(c *fiber.Ctx) error {
	var tokenString string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	}

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid token"})
	}

	claims, err := token.ParseToken(tokenString, config.AppConfig.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
	}

	id := claims.ID

	c.Locals("X-User-ID", id)

	return c.Next()
}
