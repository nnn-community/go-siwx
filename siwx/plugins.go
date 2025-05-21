package siwx

import (
    "fmt"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/storage/redis/v3"
    "os"
)

func registerCors() fiber.Handler {
    return cors.New(cors.Config{
        AllowOriginsFunc: func(_ string) bool {
            return true
        },
        AllowMethods:     "GET,HEAD,POST,PUT,DELETE,OPTIONS,PATCH",
        AllowCredentials: true,
    })
}

func registerRedis(cfg Redis) fiber.Handler {
    return func(c *fiber.Ctx) error {
        redisUrl := os.Getenv("REDIS_URL")
        redisDb := os.Getenv("REDIS_DB")

        if cfg.Url != "" {
            redisUrl = cfg.Url
        }

        if cfg.DB != "" {
            redisDb = cfg.DB
        }

        redisStorage := redis.New(redis.Config{
            URL: fmt.Sprintf("%s/%s", redisUrl, redisDb),
        })

        c.Locals("redis", redisStorage)

        return c.Next()
    }
}

func registerSiwx() fiber.Handler {
    return func(c *fiber.Ctx) error {
        rs := c.Locals("redis").(*redis.Storage)
        sid := c.Cookies("sessionId")
        user, _ := getSessionUser(rs, sid)

        c.Locals("siwx", user)

        return c.Next()
    }
}
