package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

func main() {
	app := fiber.New()

	room := NewRoom()
	go room.Run()

	app.Static("/", "./../client")

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		clientID := uuid.New().String()
		username := "Anonymous"

		client := &Client{
			ID:       clientID,
			Username: username,
			Conn:     c,
		}

		room.Register <- client
		handleClient(room, client)
	}))

	log.Println("Server is running on http://localhost:4000")
	log.Fatal(app.Listen(":4000"))
}

func handleClient(room *Room, client *Client) {
	initialMessage := Message{
		SenderID:  "Server",
		Content:   client.ID,
		Timestamp: time.Now().Format("15:04:05"),
	}

	if err := client.Conn.WriteJSON(initialMessage); err != nil {
		log.Printf("Error sending client ID to %s: %v", client.ID, err)
		room.Unregister <- client
		client.Conn.Close()
		return
	}

	defer func() {
		room.Unregister <- client
		if client.Conn != nil {
			log.Println("Closing WebSocket connection for client:", client.ID)
			client.Conn.Close()
		}
	}()

	for {
		var msg Message
		err := client.Conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message from %s: %v", client.ID, err)
			break
		}

		room.Broadcast <- msg
	}
}
