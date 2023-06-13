/*
Contains common functions, structs, and constants used by handler functions.

The response body will contain one of three keys:
- data: Returns the data requested by the client (uses outputData)
- error: Returns an error message (uses outputError)
- message: Returns a message accompanied with status code 200 OK (uses outputMessage)
*/
package controllers

import (
	"cloud.google.com/go/storage"
	"github.com/ryanozx/skillnet/database"
	"gorm.io/gorm"
)

// APIEnv is a wrapper for the shared database instance
type APIEnv struct {
	DB            *gorm.DB
	PostDBHandler database.PostDBHandler
	UserDBHandler database.UserDBHandler
	AuthDBHandler database.AuthDBHandler
	GoogleCloud   *storage.Client
}

// Creates an APIEnv object from a Database object
func CreateAPIEnv(db *gorm.DB, gc *storage.Client) *APIEnv {
	apiEnv := &APIEnv{
		DB:          db,
		GoogleCloud: gc,
	}
	return apiEnv
}
