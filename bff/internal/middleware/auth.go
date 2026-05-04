package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

const demoToken = "demo-token"

func Auth() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token := ctx.Get("Authorization")
		if token == "Bearer "+demoToken {
			ctx.Locals("userID", "me")
			return ctx.Next()
		}

		if !strings.HasPrefix(token, "Bearer user:") {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid or missing bearer token",
			})
		}

		userID := strings.TrimPrefix(token, "Bearer user:")
		if strings.TrimSpace(userID) == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid or missing bearer token",
			})
		}

		ctx.Locals("userID", userID)
		return ctx.Next()
	}
}
