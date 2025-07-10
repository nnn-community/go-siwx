package siwx

import (
    "github.com/gofiber/fiber/v2"
    "github.com/nnn-community/go-siwx/siwx/cors"
    "github.com/nnn-community/go-siwx/siwx/db"
    "github.com/nnn-community/go-siwx/siwx/redis"
    "github.com/nnn-community/go-utils/env"
    "os"
)

// Creates a new Fiber instance with GORM database with all required plugins for authentication process.
func New(config ...Config) *fiber.App {
    env.Load()
    cfg := Config{}

    if len(config) > 0 {
        cfg = config[0]
    }

    if cfg.DatabaseUrl == "" {
        cfg.DatabaseUrl = os.Getenv("DATABASE_URL")
    }

    app := fiber.New(cfg.Fiber)
    db.Connect(cfg.DatabaseUrl)

    app.Use(cors.Register())
    app.Use(redis.Register(cfg.Redis))
    app.Use(registerSiwx())

    setRoutes(app, cfg)

    return app
}
