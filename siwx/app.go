package siwx

import (
    "github.com/gofiber/fiber/v2"
)

func New(config ...Config) *fiber.App {
    cfg := Config{}

    if len(config) > 0 {
        cfg = config[0]
    }

    app := fiber.New(cfg.Fiber)

    app.Use(registerCors())
    app.Use(registerRedis(cfg.Redis))
    app.Use(registerSiwx())

    return app
}
