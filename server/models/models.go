package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username    string    `gorm:"unique;not null"`
	Email       string    `gorm:"unique;not null"`
	Password    string    `gorm:"not null"`
	DisplayName string    `gorm:"not null"`
}

type Room struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name     string    `gorm:"unique;not null"`
}

type Message struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	RoomID   uuid.UUID `gorm:"type:uuid;not null"`  
	SenderID uuid.UUID `gorm:"type:uuid;not null"`  
	Content  string    `gorm:"not null"`
	Timestamp time.Time `gorm:"not null"`
}

type RoomParticipant struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	RoomID uuid.UUID `gorm:"type:uuid;not null"` 
	UserID uuid.UUID `gorm:"type:uuid;not null"` 
}
