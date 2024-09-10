package main

import (
	"log"

	"github.com/MohamedMosalm/Chat-App/db"
	"github.com/MohamedMosalm/Chat-App/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

	roomUUID, err := uuid.Parse(roomID)
	if err != nil {
		log.Println("Invalid room ID")
		c.Close()
		return
	}

	roomManager.Mutex.Lock()
	room, exists := roomManager.Rooms[roomUUID]
	roomManager.Mutex.Unlock()

	if !exists {
		roomRecord := &models.Room{}
		result := db.DB.Where("id = ?", roomUUID).First(roomRecord)
		if result.Error == gorm.ErrRecordNotFound {
			log.Println("Room does not exist")
			c.Close()
			return
		}

		room = NewRoom(roomRecord.ID, roomRecord.Name)
		roomManager.Rooms[roomUUID] = room
		go room.Run()
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
