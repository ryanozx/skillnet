package database

import (
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
	return nil, nil
}

func (db *LikeDB) DeleteLike(userID string, postID string) error {
	return nil
}
