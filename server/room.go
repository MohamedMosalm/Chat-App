package main

import (
	"log"
	"sync"
)

type Room struct {
	ID         string
	Name       string
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
	Mutex      sync.Mutex
}

func NewRoom(id, name string) *Room {
	return &Room{
		ID:         id,
		Name:       name,
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.Register:
			r.Mutex.Lock()
			r.Clients[client.ID] = client
			r.Mutex.Unlock()
			log.Printf("Client connected: %s (%s) to room %s with ID: %s", client.Username, client.ID, r.Name, r.ID)
		case client := <-r.Unregister:
			r.Mutex.Lock()
			if _, ok := r.Clients[client.ID]; ok {
				delete(r.Clients, client.ID)
				log.Printf("Client disconnected: %s (%s) from room %s", client.Username, client.ID, r.Name)
			}
			r.Mutex.Unlock()
		case message := <-r.Broadcast:
			r.Mutex.Lock()
			for _, client := range r.Clients {
				if client.ID == message.SenderID {
					continue
				}
				log.Printf("Sending message to client: %s", client.ID)
				err := client.Conn.WriteJSON(message)
				if err != nil {
					log.Printf("Error sending message to %s: %v", client.Username, err)
					client.Conn.Close()
					delete(r.Clients, client.ID)
				}
			}
			r.Mutex.Unlock()
		}
	}
}
