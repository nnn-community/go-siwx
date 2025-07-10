package redis

import (
    "fmt"
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

var Store *redisp.Storage

func Register(cfg Config) {
    redisUrl := os.Getenv("REDIS_URL")
    redisDb := os.Getenv("REDIS_DB")

    if cfg.Url != "" {
        redisUrl = cfg.Url
    }

    if cfg.DB != "" {
        redisDb = cfg.DB
    }

    Store = redisp.New(redisp.Config{
        URL: fmt.Sprintf("%s/%s", redisUrl, redisDb),
    })
}
