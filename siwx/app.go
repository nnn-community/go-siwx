package siwx

import (
    "github.com/gofiber/fiber/v2"
    "github.com/nnn-community/go-siwx/siwx/cors"
    "github.com/nnn-community/go-siwx/siwx/db"
    "github.com/nnn-community/go-siwx/siwx/redis"
    "github.com/nnn-community/go-siwx/siwx/session"
    "os"
)

// Creates a new Fiber instance with GORM database with all required plugins for authentication process.
func NewFiber(config ...Config) *fiber.App {
    cfg := Config{}

    if len(config) > 0 {
        cfg = config[0]
    }

    if cfg.DatabaseUrl == "" {
        cfg.DatabaseUrl = os.Getenv("DATABASE_URL")
    }

    app := fiber.New(cfg.Fiber)

    db.Connect(cfg.DatabaseUrl)
    redis.Register(cfg.Redis)

    app.Use(cors.Register())
    app.Use(session.Register())

    setRoutes(app, cfg)

    return app
}
