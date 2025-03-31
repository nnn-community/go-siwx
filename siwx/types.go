package siwx

import "github.com/gofiber/fiber/v2"

type Redis struct {
    // Provide Url string to the Redis Server, omit DB from the string as it will be added from the DB option.
    //
    // Optional, default: os.Getenv("REDIS_URL")
    Url string `json:"url"`

    // DB number where the session is stored.
    //
    // Optional, default: os.Getenv("REDIS_DB")
    DB string `json:"db"`
}

type Config struct {
    Fiber fiber.Config
    Redis Redis
}
