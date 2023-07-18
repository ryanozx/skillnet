package database

import (
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
)

type LikeAPIHandler interface {
	CreateLike(*models.Like) (*models.Like, error)
	DeleteLike(string, uint) error
	GetValue(uint) (uint64, error)
	GetLikeByID(string) (*models.Like, error)
}

type DBValueGetter interface {
	GetValue(uint) (uint64, error)
}

type LikeDB struct {
	DB *gorm.DB
}

func (db *LikeDB) CreateLike(like *models.Like) (*models.Like, error) {
	result := db.DB.Create(like)
	if result.Error != nil {
		return like, result.Error
	}
	return db.GetLikeByID(like.ID)
}

func (db *LikeDB) DeleteLike(userID string, postID uint) error {
	err := db.DB.Unscoped().Delete(&models.Like{}, "id = ?", helpers.GenerateLikeID(userID, postID)).Error
	return err
}

func (db *LikeDB) GetValue(postID uint) (uint64, error) {
	var count int64
	result := db.DB.Model(&models.Like{}).Where("post_id = ?", postID).Count(&count)
	return uint64(count), result.Error
}

func (db *LikeDB) GetLikeByID(likeID string) (*models.Like, error) {
	like := models.Like{}
	err := db.DB.Joins("Post").Joins("User").First(&like, "likes.id = ?", likeID).Error
	return &like, err
}
