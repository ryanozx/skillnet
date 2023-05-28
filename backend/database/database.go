package database

import (
	"fmt"
	"log"
	"os"

	"github.com/ryanozx/skillnet/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database *gorm.DB

func ConnectDatabase() {

	db, dbErr := establishConnection()
	if dbErr != nil {
		log.Fatal("Failed to connect to database. \n", dbErr)
		os.Exit(2)
	}

	log.Println("Connected")
	db.Logger = logger.Default.LogMode(logger.Info)

	autoMigrate(db)

	Database = db
}

func establishConnection() (*gorm.DB, error) {
	dataSourceName := fmt.Sprintf("host=db user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Singapore",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	gormOptions := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	return gorm.Open(postgres.Open(dataSourceName), gormOptions)
}

// update this function with any new schemas
func autoMigrate(database *gorm.DB) {
	log.Println("Running migrations")
	database.AutoMigrate(&models.PostSchema{})
	database.AutoMigrate(&models.UserWithSettings{})
	database.AutoMigrate(&models.UserCredentials{})
	database.AutoMigrate(&models.CommentSchema{})
}
