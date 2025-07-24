package goupload

import (
    "fmt"
    "github.com/gofiber/fiber/v2"
    "github.com/nnn-community/go-siwx/siwx"
    "github.com/nnn-community/go-siwx/siwx/redis"
    "github.com/nnn-community/go-siwx/siwx/user"
    "github.com/nnn-community/go-utils/env"
    "os"
)

func main() {
    env.Load()

    app := siwx.NewFiber(siwx.Config{
        Fiber: fiber.Config{
            BodyLimit: 4 * 1024 * 1024,
        },
        Redis: redis.Config{
            Url: os.Getenv("REDIS_URL"),
            DB:  os.Getenv("REDIS_DB"),
        },
    })

    // Use middleware to validate user (sends error 401 if unauthenticated)
    app.Post("/save-profile", siwx.Authenticated, func(c *fiber.Ctx) error {
        u := user.Get(c)

        fmt.Println("logged user:", u)
        fmt.Println("is admin:", u.Can([]string{"is-admin"}))

        return c.SendStatus(fiber.StatusOK)
    })

    app.Listen(":" + os.Getenv("PORT"))
}
