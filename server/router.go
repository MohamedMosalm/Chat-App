package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupRoutes(app *fiber.App, roomManager *RoomManager) {
	app.Post("/create-room", func(c *fiber.Ctx) error {
		return CreateRoomHandler(c, roomManager)
	})

	app.Get("/get-room", func(c *fiber.Ctx) error {
		return GetRoomHandler(c, roomManager)
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		WebSocketHandler(c, roomManager)
	}))
}
