package db

import (
	"fmt"
	"log"
	"os"

	"github.com/MohamedMosalm/Chat-App/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load() 
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	MigrateDB()
}

func MigrateDB() {
	err := DB.AutoMigrate(&models.User{}, &models.Room{}, &models.Message{}, &models.RoomParticipant{})
	if err != nil {
		log.Fatal("Failed to migrate models:", err)
	}
}
