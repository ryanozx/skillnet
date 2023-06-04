package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ryanozx/skillnet/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database *gorm.DB

func ConnectProdDatabase() *gorm.DB {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	dataSourceName := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%v",
		host,
		user,
		password,
		name,
		port,
	)
	db := connectDatabase(dataSourceName)
	return db
}

func ConnectTestDatabase() *gorm.DB {
	host := os.Getenv("DB_TEST_HOST")
	port := os.Getenv("DB_TEST_PORT")
	name := os.Getenv("DB_TEST_NAME")
	user := os.Getenv("DB_TEST_USER")
	password := os.Getenv("DB_TEST_PASSWORD")
	dataSourceName := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%v sslmode=disable",
		host,
		user,
		password,
		name,
		port,
	)
	db := connectDatabase(dataSourceName)
	return db
}

func connectDatabase(dataSourceName string) *gorm.DB {
	db, dbErr := establishConnection(dataSourceName)
	if dbErr != nil {
		log.Fatal("Failed to connect to database. \n", dbErr)
		os.Exit(2)
	}

	log.Println("Connected")
	db.Logger = logger.Default.LogMode(logger.Info)
	autoMigrate(db)
	Database = db
	return db
}

func establishConnection(dataSourceName string) (*gorm.DB, error) {
	gormOptions := initialiseGormConfigurations("")
	return gorm.Open(postgres.Open(dataSourceName), gormOptions)
}

// update this function with any new schemas
func autoMigrate(database *gorm.DB) {
	log.Println("Running migrations")
	database.AutoMigrate(&models.PostSchema{})
	database.AutoMigrate(&models.UserWithSettings{})
	//database.AutoMigrate(&models.CommentSchema{})
}

/*
Pass in an empty string for UTC.
Timestamps are returned with microsecond precision - Go's timestamps
have nanosecond precision while PostGreSQL has microsecond precision;
hence we truncate the timestamps created in Go to microsecond precision
/to ensure equality.
*/
func initialiseGormConfigurations(timezoneLocation string) *gorm.Config {
	options := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			location, _ := time.LoadLocation(timezoneLocation)
			return time.Now().Truncate(time.Microsecond).In(location)
		},
	}
	return options
}
