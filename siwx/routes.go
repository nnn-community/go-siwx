package siwx

import (
    "encoding/json"
    "github.com/gofiber/fiber/v2"
    "github.com/nnn-community/go-siwx/siwx/redis"
    "github.com/nnn-community/go-siwx/siwx/session"
    "github.com/nnn-community/go-siwx/siwx/user"
    "github.com/nnn-community/go-siwx/siwx/verification"
    "github.com/nnn-community/go-utils/env"
    "github.com/nnn-community/go-utils/strings"
    "time"
)

type siwxData struct {
    Address string `json:"address"`
    ChainId int    `json:"chainId"`
}

func setRoutes(app *fiber.App, cfg Config) {
    //
    // Get current logged user.
    //
    app.Get("/siwx/session", func(c *fiber.Ctx) error {
        sess := session.Get(c)

        if !sess.User.IsLoggedIn {
            return c.Status(fiber.StatusOK).JSON(nil)
        }

        return c.Status(fiber.StatusOK).JSON(siwxData{
            Address: sess.Siwx.Address,
            ChainId: strings.ToInt(sess.Siwx.ChainId, 0),
        })
    })

    //
    // Verify signature and save the user to as redis session.
    //
    app.Post("/siwx/verify", func(c *fiber.Ctx) error {
        payload := struct {
            Message   string `json:"message"`
            Signature string `json:"signature"`
        }{}

        if err := c.BodyParser(&payload); err != nil {
            return c.SendStatus(fiber.StatusBadRequest)
        }

        v, err := verification.VerifyUser(payload.Message, payload.Signature)

        if err != nil || !v.Success {
            return c.SendStatus(fiber.StatusUnauthorized)
        }

        sid := session.GetId(c)
        u, err := user.GetByAddress(v.Result.Address, cfg.GetUserData)

        if err != nil {
            return c.SendStatus(fiber.StatusUnauthorized)
        }

        var expiration time.Time

        if cfg.CookieDuration > 0 {
            expiration = time.Now().Add(cfg.CookieDuration)
        }

        cookieCfg := fiber.Cookie{
            Name:     "session-id",
            Value:    sid,
            Expires:  expiration,
            Path:     "/",
            Secure:   !env.IsLocal(),
            HTTPOnly: true,
        }

        if !env.IsLocal() {
            cookieCfg.SameSite = "none"
        }

        if !env.IsLocal() && cfg.CookieDomain != "" {
            cookieCfg.Domain = cfg.CookieDomain
        }

        sessionData, _ := json.Marshal(session.Data{
            Id:   sid,
            Siwx: v.Result,
            User: &u,
        })

        c.Cookie(&cookieCfg)
        redis.Store.Set(sid, sessionData, cfg.CookieDuration)

        return c.Status(fiber.StatusOK).JSON(true)
    })

    //
    // Logout current user.
    //
    app.Delete("/siwx/session", func(c *fiber.Ctx) error {
        session.Delete(c)
        c.Locals("siwx", nil)

        return c.Status(fiber.StatusOK).JSON(true)
    })

    //
    // Get logged user data/profile.
    //
    app.Get("/siwx/profile", Authenticated, func(c *fiber.Ctx) error {
        return c.Status(fiber.StatusOK).JSON(user.Get(c))
    })
}
