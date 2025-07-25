package siwx

import (
    "encoding/json"
    "errors"
    "github.com/gofiber/fiber/v2"
    "github.com/nnn-community/go-siwx/siwx/cors"
    "github.com/nnn-community/go-siwx/siwx/db"
    "github.com/nnn-community/go-siwx/siwx/redis"
    "github.com/nnn-community/go-siwx/siwx/session"
    "github.com/nnn-community/go-siwx/siwx/user"
    "os"
)

// Creates a new Fiber instance with GORM database with all required plugins for authentication process.
func NewFiber(config ...Config) *fiber.App {
    if len(config) > 0 {
        AppConfig = config[0]
    }

    if AppConfig.DatabaseUrl == "" {
        AppConfig.DatabaseUrl = os.Getenv("DATABASE_URL")
    }

    app := fiber.New(AppConfig.Fiber)

    db.Connect(AppConfig.DatabaseUrl)
    redis.Register(AppConfig.Redis)

    app.Use(cors.Register())
    app.Use(session.Register())

    setRoutes(app)

    return app
}

// Refresh the session of the current user.
func Refresh(c *fiber.Ctx) error {
    sess := session.Get(c)

    if sess.User == nil {
        return errors.New("not authenticated")
    }

    updated, err := user.GetByID(sess.User.ID, AppConfig.GetUserData)

    if err != nil {
        return errors.New("refresh failed")
    }

    sessionData, _ := json.Marshal(session.Data{
        Id:   sess.Id,
        Siwx: sess.Siwx,
        User: &updated,
    })

    redis.Store.Set(sess.Id, sessionData, AppConfig.CookieDuration)

    c.Locals("siwx.user", updated)

    return nil
}
