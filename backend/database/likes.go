package database

import (
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
)

type LikeAPIHandler interface {
	CreateLike(*models.Like) (*models.Like, error)
	DeleteLike(string, string) error
	GetLikeCount(string) (uint64, error)
}

type LikeDBCountGetter interface {
	GetLikeCount(string) (uint64, error)
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

func (db *LikeDB) GetLikeCount(postID string) (uint64, error) {
	var count int64
	result := db.DB.Model(&models.Like{}).Where("post_id = ?", postID).Count(&count)
	return uint64(count), result.Error
}
