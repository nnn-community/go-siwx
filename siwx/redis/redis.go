package redis

import (
    "fmt"
    "github.com/gofiber/fiber/v2"
    redisp "github.com/gofiber/storage/redis/v3"
    "os"
)

type Config struct {
    // Provide Url string to the Redis Server, omit DB from the string as it will be added from the DB option.
    //
    // Optional, default: os.Getenv("REDIS_URL")
    Url string `json:"url"`

    // DB number where the session is stored.
    //
    // Optional, default: os.Getenv("REDIS_DB")
    DB string `json:"db"`
}

func Register(cfg Config) fiber.Handler {
    return func(c *fiber.Ctx) error {
        redisUrl := os.Getenv("REDIS_URL")
        redisDb := os.Getenv("REDIS_DB")

        if cfg.Url != "" {
            redisUrl = cfg.Url
        }

        if cfg.DB != "" {
            redisDb = cfg.DB
        }

        redisStorage := redisp.New(redisp.Config{
            URL: fmt.Sprintf("%s/%s", redisUrl, redisDb),
        })

        c.Locals("redis", redisStorage)

        return c.Next()
    }
}

func Get(c *fiber.Ctx) *redisp.Storage {
    return c.Locals("redis").(*redisp.Storage)
}
