package goupload

import (
    "fmt"
    "github.com/gofiber/fiber/v2"
    "github.com/nnn-community/go-siwx/middleware"
    "github.com/nnn-community/go-siwx/siwx"
    "github.com/nnn-community/go-utils/env"
    "os"
)

func main() {
    env.Load()

    app := siwx.New(siwx.Config{
        Fiber: fiber.Config{
            BodyLimit: 4 * 1024 * 1024,
        },
        Redis: siwx.Redis{
            Url: os.Getenv("REDIS_URL"),
            DB:  os.Getenv("REDIS_DB"),
        },
    })

    // Use middleware to validate user (sends error 401 if unauthenticated)
    app.Post("/save-user", middleware.SiwxMiddleware, func(c *fiber.Ctx) error {
        user := c.Locals("siwx").(siwx.User)

        fmt.Println("logged user:", user)
        fmt.Println("is admin:", user.Can([]string{"is-admin"}))

        return c.SendStatus(fiber.StatusOK)
    })

    app.Listen(":" + os.Getenv("PORT"))
}
