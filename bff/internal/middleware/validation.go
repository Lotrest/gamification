package middleware

import "github.com/gofiber/fiber/v2"

func RequireRouteParam(param string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if ctx.Params(param) == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "required route param is missing",
				"param":   param,
			})
		}

		return ctx.Next()
	}
}
