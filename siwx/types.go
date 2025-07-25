package siwx

import (
    "github.com/gofiber/fiber/v2"
    "github.com/nnn-community/go-siwx/siwx/redis"
    "github.com/nnn-community/go-siwx/siwx/user"
    "time"
)

type Config struct {
    // Custom Fiber configuration to apply.
    //
    // Optional
    Fiber fiber.Config

    // Redis connection settings to store session. Read redis.Config for more info.
    //
    // Optional
    Redis redis.Config

    // DatabaseUrl connection to your database to manage users.
    //
    // Optional, default: os.Getenv("DATABASE_URL")
    DatabaseUrl string

    // CookieDuration how long will cookies last.
    //
    // Optional, default: 0 (0 = session only)
    CookieDuration time.Duration

    // Set CookieDomain if you are using auth in cross-domains (omitted on localhost).
    //
    // Optional, default: undefined
    CookieDomain string

    // Define GetUserData to get `user-data` for client (ie. profile info).
    //
    // Optional, default: nil
    GetUserData *func(user user.User) map[string]interface{}
}

var AppConfig Config
