package database

import (
	"log"
	"os"
	"time"

	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// Returns the production database
func ConnectProdDatabase() *gorm.DB {
	env := helpers.RetrieveProdDBEnv()
	db := connectDatabase(env)
	return db
}

// Returns the test database
func ConnectTestDatabase() *gorm.DB {
	env := helpers.RetrieveTestDBEnv()
	db := connectDatabase(env)
	return db
}

type DataSourceNamer interface {
	DataSourceName() string
}

// Returns a database connection
func connectDatabase(env DataSourceNamer) *gorm.DB {
	db, err := establishConnection(env)
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}
	log.Println("Connected")

	db.Logger = logger.Default.LogMode(logger.Info)
	autoMigrate(db)
	return db
}

// Connects to a database based on the data source name supplied by DataSourceNamer
func establishConnection(env DataSourceNamer) (*gorm.DB, error) {
	dataSourceName := env.DataSourceName()
	gormOptions := initialiseGormConfigurations("")
	return gorm.Open(postgres.Open(dataSourceName), gormOptions)
}

// Performs migration automatically based on schemas specified in method body
func autoMigrate(database *gorm.DB) {
	log.Println("Running migrations")
	database.AutoMigrate(&models.Post{}, &models.User{}, &models.Like{}, &models.Comment{})
	// Add more schemas above as necessary
}

// Pass in an empty string for UTC.
func initialiseGormConfigurations(timezoneLocation string) *gorm.Config {
	/*
		Timestamps are returned with microsecond precision - Go's timestamps
		have nanosecond precision while PostGreSQL has microsecond precision;
		hence we truncate the timestamps created in Go to microsecond precision
		/to ensure equality.
	*/
	nowFunc := func() time.Time {
		location, _ := time.LoadLocation(timezoneLocation)
		return time.Now().Truncate(time.Microsecond).In(location)
	}
	options := &gorm.Config{
		SkipDefaultTransaction:                   false,
		NamingStrategy:                           nil,
		FullSaveAssociations:                     false,
		Logger:                                   logger.Default.LogMode(logger.Info),
		NowFunc:                                  nowFunc,
		DryRun:                                   false,
		PrepareStmt:                              false,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		IgnoreRelationshipsWhenMigrating:         false,
		DisableNestedTransaction:                 false,
		AllowGlobalUpdate:                        false,
		QueryFields:                              false,
		CreateBatchSize:                          0,
		TranslateError:                           false,
		ClauseBuilders:                           map[string]clause.ClauseBuilder{},
		ConnPool:                                 nil,
		Dialector:                                nil,
		Plugins:                                  map[string]gorm.Plugin{},
	}
	return options
}
