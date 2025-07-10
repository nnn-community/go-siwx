package siwx

import (
    "github.com/gofiber/fiber/v2"
    "github.com/nnn-community/go-siwx/siwx/user"
)

func Authenticated(c *fiber.Ctx) error {
    u := user.Get(c)

    if !u.IsLoggedIn {
        return c.SendStatus(fiber.StatusUnauthorized)
    }

    return c.Next()
}

func Can(permissions []string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        u := user.Get(c)
        valid := u.IsLoggedIn && u.Can(permissions)

        if !valid {
            return c.SendStatus(fiber.StatusUnauthorized)
        }

        return c.Next()
    }
}

func CanAll(permissions []string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        u := user.Get(c)
        valid := u.IsLoggedIn && u.CanAll(permissions)

        if !valid {
            return c.SendStatus(fiber.StatusUnauthorized)
        }

        return c.Next()
    }
}
