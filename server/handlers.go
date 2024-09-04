package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

func CreateRoomHandler(c *fiber.Ctx, roomManager *RoomManager) error {
	roomName := c.FormValue("room_name")
	if roomName == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Room name is required")
	}

	room := roomManager.CreateRoom(roomName)
	return c.JSON(fiber.Map{
		"room_id": room.ID,
	})
}

func GetRoomHandler(c *fiber.Ctx, roomManager *RoomManager) error {
	roomName := c.Query("room_name")
	if roomName == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Room name is required")
	}

	room, found := roomManager.GetRoomByName(roomName)
	if !found {
		return c.Status(fiber.StatusNotFound).SendString("Room not found")
	}

	return c.JSON(fiber.Map{
		"room_id": room.ID,
	})
}

func WebSocketHandler(c *websocket.Conn, roomManager *RoomManager) {
	roomID := c.Query("room_id")
	clientID := c.Query("client_id")

	if roomID == "" {
		log.Println("Room ID is missing")
		c.Close()
		return
	}

	roomManager.Mutex.Lock()
	room, exists := roomManager.Rooms[roomID]
	roomManager.Mutex.Unlock()
	if !exists {
		log.Println("Room does not exist")
		c.Close()
		return
	}

	if clientID == "" {
		clientID = uuid.New().String()
	}

	client := &Client{
		ID:       clientID,
		Username: "Anonymous",
		Conn:     c,
		Room:     room,
	}

	room.Register <- client
	handleClient(roomManager, client)
}
