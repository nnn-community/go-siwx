package siwx

import (
    "github.com/gofiber/fiber/v2"
    "github.com/nnn-community/go-siwx/siwx/cors"
    "github.com/nnn-community/go-siwx/siwx/db"
    "github.com/nnn-community/go-siwx/siwx/redis"
    "github.com/nnn-community/go-siwx/siwx/session"
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
    redis.Register(cfg.Redis)

    app.Use(cors.Register())
    app.Use(localSiwx())

    setRoutes(app, cfg)

    return app
}

func localSiwx() fiber.Handler {
    return func(c *fiber.Ctx) error {
        s := session.Get(c)

        c.Locals("siwx.user", s.User)

        return c.Next()
    }
}
