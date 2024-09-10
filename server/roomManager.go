package main

import (
	"log"
	"sync"

	"github.com/MohamedMosalm/Chat-App/db"
	"github.com/MohamedMosalm/Chat-App/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomManager struct {
	Rooms map[uuid.UUID]*Room
	Mutex sync.Mutex
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		Rooms: make(map[uuid.UUID]*Room),
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

	var dbRoom models.Room
	result := db.DB.Where("name = ?", name).First(&dbRoom)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, false
	}

	room := NewRoom(dbRoom.ID, dbRoom.Name)
	rm.Rooms[dbRoom.ID] = room
	go room.Run()
	return room, true
}

func (rm *RoomManager) CreateRoom(name string) *Room {
	rm.Mutex.Lock()
	defer rm.Mutex.Unlock()

	roomID := uuid.New()
	newRoom := NewRoom(roomID, name)
	rm.Rooms[roomID] = newRoom
	go newRoom.Run()

	newRoomRecord := models.Room{
		ID:        roomID,
		Name:      name,
	}
	if err := db.DB.Create(&newRoomRecord).Error; err != nil {
		log.Printf("Error creating room in the database: %v", err)
	}

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