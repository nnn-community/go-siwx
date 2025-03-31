package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/nnn-community/go-siwx/siwx"
)

func SiwxMiddleware(c *fiber.Ctx) error {
    user := c.Locals("siwx").(siwx.User)

    if !user.IsLoggedIn {
        return c.SendStatus(fiber.StatusUnauthorized)
    }

    return c.Next()
}
