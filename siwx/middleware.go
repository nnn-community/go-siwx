package siwx

import (
    "github.com/gofiber/fiber/v2"
)

func Middleware(c *fiber.Ctx) error {
    user := c.Locals("siwx").(User)

    if !user.IsLoggedIn {
        return c.SendStatus(fiber.StatusUnauthorized)
    }

    return c.Next()
}
