package cors

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Register() fiber.Handler {
	return cors.New(cors.Config{
		AllowOriginsFunc: func(_ string) bool {
			return true
		},
		AllowMethods:     "GET, HEAD, POST, PUT, DELETE, CONNECT, OPTIONS, TRACE, PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept, Cache-Control, Pragma, Referer, User-Agent, Authorization, Upgrade, Connection, Sec-Ch-Ua, Sec-Ch-Ua-Mobile, Sec-Ch-Ua-Platform, Sec-WebSocket-Key, Sec-WebSocket-Version, Sec-WebSocket-Extensions, Guardian-Api-Key, X-Guardian-Api-Key, Guardian-Webhook-Key, X-Guardian-Webhook-Key, X-Key",
		AllowCredentials: true,
	})
}
