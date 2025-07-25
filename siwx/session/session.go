package session

import (
    "encoding/json"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/encryptcookie"
    "github.com/nnn-community/go-siwx/siwx/redis"
    "github.com/nnn-community/go-siwx/siwx/user"
)

type Data struct {
    Id   string     `json:"id"`
    Siwx *Siwx      `json:"siwx"`
    User *user.User `json:"user"`
}

type Siwx struct {
    Address string `json:"address"`
    ChainId string `json:"chain-id"`
}

var Unauthorized = Data{
    Id: "",
    Siwx: &Siwx{
        Address: "",
        ChainId: "",
    },
    User: &user.Unauthorized,
}

func Get(c *fiber.Ctx) Data {
    sessId := GetId(c)
    s, err := redis.Store.Get(sessId)

    if err != nil {
        return Unauthorized
    }

    if string(s) == "" {
        return Unauthorized
    }

    var data Data

    if err := json.Unmarshal(s, &data); err != nil {
        return Unauthorized
    }

    if data.User == nil {
        return Unauthorized
    }

    data.Id = sessId

    return data
}

func GetId(c *fiber.Ctx) string {
    sid := c.Cookies("session-id")

    if sid == "" {
        return encryptcookie.GenerateKey()
    }

    return sid
}

func Delete(c *fiber.Ctx) {
    sessId := GetId(c)
    redis.Store.Delete(sessId)
    c.Locals("siwx.user", nil)
}

func Register() fiber.Handler {
    return func(c *fiber.Ctx) error {
        s := Get(c)
        c.Locals("siwx.user", s.User)

        return c.Next()
    }
}
