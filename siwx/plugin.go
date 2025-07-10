package siwx

import (
    "github.com/gofiber/fiber/v2"
    "github.com/nnn-community/go-siwx/siwx/session"
)

func registerSiwx() fiber.Handler {
    return func(c *fiber.Ctx) error {
        s := session.Get(c)

        c.Locals("siwx.user", s.User)

        return c.Next()
    }
}
