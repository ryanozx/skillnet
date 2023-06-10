package database

import (
	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, userCreds *models.UserCredentials) (*models.User, error) {
	newUser := userCreds.NewUser()
	result := db.Create(newUser)
	if err := result.Error; err != nil {
		return newUser, err
	}
	return newUser, nil
}

func GetUserByID(db *gorm.DB, id string) (*models.User, error) {
	user := models.User{}
	err := db.First(&user, "id = ?", id).Error
	if err != nil {
		return &user, err
	}
	return &user, nil
}

/*
Since ID will not change after creation, GetUserByID is preferred over GetUserByUsername for
operations involving updating/deleting.
*/
func GetUserByUsername(db *gorm.DB, username string) (*models.User, error) {
	user := models.User{}
	if err := db.First(&user, "username = ?", username).Error; err != nil {
		return &user, err
	}
	return &user, nil
}

func UpdateUser(db *gorm.DB, user *models.User, id string) (*models.User, error) {
	originalUser, err := GetUserByID(db, id)
	if err != nil {
		return originalUser, err
	}
	err = db.Model(originalUser).Updates(user).Error
	return originalUser, err
}

func DeleteUser(db *gorm.DB, id string) error {
	user, err := GetUserByID(db, id)
	if err != nil {
		return err
	}
	err = db.Delete(&user).Error
	return err
}
