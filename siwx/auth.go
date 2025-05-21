package siwx

import (
    "database/sql"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/storage/redis/v3"
)

func addAuthRoutes(app *fiber.App, db *sql.DB) {
    app.Get("/siwx/session", func(c *fiber.Ctx) error {
        rs := c.Locals("redis").(*redis.Storage)
        sid := c.Cookies("sessionId")
        siwx, err := getSessionSiwx(rs, sid)

        if err != nil {
            return c.Status(fiber.StatusOK).JSON(nil)
        }

        return c.Status(fiber.StatusOK).JSON(siwx)
    })

    app.Post("/siwx/verify", func(c *fiber.Ctx) error {
        if db == nil {
            return c.SendStatus(fiber.StatusNotFound)
        }

        payload := struct {
            Message   string `json:"message"`
            Signature string `json:"signature"`
        }{}

        if err := c.BodyParser(&payload); err != nil {
            return c.SendStatus(fiber.StatusBadRequest)
        }

        verification, err := VerifyUser(db, payload.Message, payload.Signature)

        if err != nil || !verification.Success {
            return c.SendStatus(fiber.StatusUnauthorized)
        }

        rs := c.Locals("redis").(*redis.Storage)
        sid := c.Cookies("sessionId")
        sessId, err := parseSessionId(sid)

        return c.Status(fiber.StatusOK).JSON(true)
    })

    app.Delete("/siwx/session", func(c *fiber.Ctx) error {
        rs := c.Locals("redis").(*redis.Storage)
        sid := c.Cookies("sessionId")
        sessId, err := parseSessionId(sid)

        if err == nil {
            rs.Delete(*sessId)
        }

        c.Locals("siwx", nil)

        return c.Status(fiber.StatusOK).JSON(true)
    })

    app.Get("/siwx/profile", func(c *fiber.Ctx) error {
        return c.Status(fiber.StatusOK).JSON(GetUser(c))
    })
}
