package database

import (
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
)

type LikeAPIHandler interface {
	CreateLike(*models.Like) (*models.Like, error)
	DeleteLike(string, string) error
}

type LikeDB struct {
	DB *gorm.DB
}

func (db *LikeDB) CreateLike(like *models.Like) (*models.Like, error) {
	result := db.DB.Create(like)
	return like, result.Error
}

func (db *LikeDB) DeleteLike(userID string, postID string) error {
	err := db.DB.Unscoped().Delete(&models.Like{}, "id = ?", helpers.GenerateLikeID(userID, postID)).Error
	return err
}
