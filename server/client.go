package main

import (
	"log"
	"time"

	"github.com/MohamedMosalm/Chat-App/db"
	"github.com/MohamedMosalm/Chat-App/models"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
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

    var messages []models.Message
    if err := db.DB.Where("room_id = ?", client.Room.ID).Order("timestamp asc").Find(&messages).Error; err != nil {
        log.Printf("Error loading messages for room %s: %v", client.Room.ID, err)
    }

    for _, msg := range messages {
        wsMessage := Message{
            SenderID:  msg.SenderID.String(),
            Content:   msg.Content,
            Timestamp: msg.Timestamp.Format("15:04:05"),
        }
        if err := client.Conn.WriteJSON(wsMessage); err != nil {
            log.Printf("Error sending previous message to %s: %v", client.ID, err)
            return
        }
    }

    for {
        var msg Message
        err := client.Conn.ReadJSON(&msg)
        if err != nil {
            log.Printf("Error reading message from %s: %v", client.ID, err)
            break
        }

        var senderUUID uuid.UUID
        if client.ID != "Server" {
            senderUUID, err = uuid.Parse(client.ID)  
            if err != nil {
                log.Printf("Error parsing sender ID: %v", err)
                break
            }
        }

        newMessage := models.Message{
            RoomID:   client.Room.ID,
            SenderID: senderUUID,
            Content:  msg.Content,
            Timestamp: time.Now(),
        }
        
        if err := db.DB.Create(&newMessage).Error; err != nil {
            log.Printf("Error saving message to the database: %v", err)
        }
        

        msg.SenderID = client.ID
        client.Room.Broadcast <- msg
    }
}