package database

import (
	"log"
	"os"
	"time"

	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectProdDatabase() *gorm.DB {
	env := helpers.RetrieveProdDBEnv()
	db := connectDatabase(env)
	return db
}

func ConnectTestDatabase() *gorm.DB {
	env := helpers.RetrieveTestDBEnv()
	db := connectDatabase(env)
	return db
}

type DataSourceNamer interface {
	DataSourceName() string
}

func connectDatabase(env DataSourceNamer) *gorm.DB {
	db, dbErr := establishConnection(env)
	if dbErr != nil {
		log.Fatal("Failed to connect to database. \n", dbErr)
		os.Exit(2)
	}

	log.Println("Connected")
	db.Logger = logger.Default.LogMode(logger.Info)
	autoMigrate(db)
	return db
}

func establishConnection(env DataSourceNamer) (*gorm.DB, error) {
	dataSourceName := env.DataSourceName()
	gormOptions := initialiseGormConfigurations("")
	return gorm.Open(postgres.Open(dataSourceName), gormOptions)
}

// update this function with any new schemas
func autoMigrate(database *gorm.DB) {
	log.Println("Running migrations")
	database.AutoMigrate(&models.Post{})
	database.AutoMigrate(&models.User{})
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
