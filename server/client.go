package main

import (
	"log"
	"time"

	"github.com/gofiber/websocket/v2"
)

type Client struct {
	ID       string
	Username string
	Conn     *websocket.Conn
	Room     *Room
}

func handleClient(roomManager *RoomManager, client *Client) {
    initialMessage := Message{
        SenderID:  client.ID, 
        Content:   "Connected to the server",
        Timestamp: time.Now().Format("15:04:05"),
    }

    if err := client.Conn.WriteJSON(initialMessage); err != nil {
        log.Printf("Error sending initial message to %s: %v", client.ID, err)
        client.Room.Unregister <- client
        client.Conn.Close()
        return
    }

    defer func() {
        client.Room.Unregister <- client
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

        msg.SenderID = client.ID
        log.Printf("Received message from client %s: %s", client.ID, msg.Content)

        client.Room.Broadcast <- msg
    }
}