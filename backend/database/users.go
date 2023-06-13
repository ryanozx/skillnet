package database

import (
	"log"

	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserDBHandler interface {
	CreateUser(NewUser) (*models.User, error)
	DeleteUser(string) error
	GetUserByID(string) (*models.User, error)
	GetUserByUsername(string) (*models.User, error)
	UpdateUser(*models.User, string) (*models.User, error)
}

// UserDB implements both UserDBHandler and AuthAPIHandler
type UserDB struct {
	DB *gorm.DB
}

// Creates a new user from a username and password
func (db *UserDB) CreateUser(userCreds NewUser) (*models.User, error) {
	newUser := userCreds.NewUser()
	result := db.DB.Model(models.User{}).Create(newUser)
	err := result.Error
	return newUser, err
}

type NewUser interface {
	NewUser() *models.User
}

// Deletes user's profile
func (db *UserDB) DeleteUser(id string) error {
	user, err := db.GetUserByID(id)
	if err != nil {
		return err
	}
	err = db.DB.Delete(&user).Error
	return err
}

// Retrieves a User object by ID
func (db *UserDB) GetUserByID(id string) (*models.User, error) {
	log.Println(id)
	user := models.User{}
	err := db.DB.First(&user, "id = ?", id).Error
	return &user, err
}

// Retrieves user by username. Since ID will not change after creation, GetUserByID is preferred
// over GetUserByUsername for update/delete operations
func (db *UserDB) GetUserByUsername(username string) (*models.User, error) {
	user := models.User{}
	err := db.DB.First(&user, "username = ?", username).Error
	return &user, err
}

// Updates user's profile.
func (db *UserDB) UpdateUser(user *models.User, id string) (*models.User, error) {
	resUser := &models.User{}
	result := db.DB.Model(resUser).Clauses(clause.Returning{}).Where("id = ?", id).Updates(user)
	err := result.Error
	return resUser, err
}
