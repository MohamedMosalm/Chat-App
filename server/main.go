package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	roomManager := NewRoomManager()

	SetupRoutes(app, roomManager)

	app.Static("/", "./../client")

	log.Println("Server is running on http://localhost:4000")
	log.Fatal(app.Listen(":4000"))
}