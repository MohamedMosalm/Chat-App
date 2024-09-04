package main

import (
	"log"
	"sync"

	"github.com/google/uuid"
)

type RoomManager struct {
	Rooms map[string]*Room
	Mutex sync.Mutex
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		Rooms: make(map[string]*Room),
	}
}

func (rm *RoomManager) GetRoomByName(name string) (*Room, bool) {
	rm.Mutex.Lock()
	defer rm.Mutex.Unlock()

	for _, room := range rm.Rooms {
		if room.Name == name {
			return room, true
		}
	}
	return nil, false
}

func (rm *RoomManager) CreateRoom(name string) *Room {
	rm.Mutex.Lock()
	defer rm.Mutex.Unlock()

	roomID := uuid.New().String()
	newRoom := NewRoom(roomID, name)
	rm.Rooms[roomID] = newRoom
	go newRoom.Run()
	log.Printf("Created new room: %s (ID: %s)", name, roomID)
	return newRoom
}

func (rm *RoomManager) GetOrCreateRoom(name string) *Room {
	room, exists := rm.GetRoomByName(name)
	if exists {
		return room
	}
	return rm.CreateRoom(name)
}