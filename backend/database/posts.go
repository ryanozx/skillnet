package database

import (
	"errors"
	"fmt"

	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrNotOwner = errors.New("unauthorised action")

const postsToReturn = 10

type PostDBHandler interface {
	CreatePost(*models.Post) (*models.Post, error)
	DeletePost(string, string) error
	GetPosts(string, string) ([]models.Post, error)
	GetPostByID(string, string) (*models.Post, error)
	UpdatePost(*models.Post, string, string) (*models.Post, error)
}

// PostDB implements PostDBHandler
type PostDB struct {
	DB *gorm.DB
}

func (db *PostDB) CreatePost(post *models.Post) (*models.Post, error) {
	result := db.DB.Create(post)
	if result.Error != nil {
		return post, result.Error
	}
	return db.GetPostByID(fmt.Sprintf("%v", post.ID), post.UserID)
}

func (db *PostDB) DeletePost(id string, userID string) error {
	post, err := db.GetPostByID(id, userID)
	if err != nil {
		return err
	}
	if err := checkUserIsOwner(post, userID); err != nil {
		return err
	}
	err = db.DB.Delete(&post).Error
	return err
}

func (db *PostDB) GetPosts(cutoff string, userID string) ([]models.Post, error) {
	var posts []models.Post

	if cutoff == "" {
		// Retrieve all posts from database
		query := db.DB.Limit(postsToReturn).Joins("User").Preload("Likes").Joins("LEFT JOIN likes ON (posts.ID = likes.post_id AND likes.user_id = ?)", userID).Order("posts.id desc").Find(&posts)
		if err := query.Find(&posts).Error; err != nil {
			return posts, err
		}
	} else {
		query := db.DB.Where("posts.id < ?", cutoff).Joins("User").Limit(postsToReturn).Preload("Likes").Joins("LEFT JOIN likes ON (posts.ID = likes.post_id AND likes.user_id = ?)", userID).Order("posts.id desc").Find(&posts)
		if err := query.Find(&posts).Error; err != nil {
			return posts, err
		}
	}
	return posts, nil
}

func (db *PostDB) GetPostByID(id string, userID string) (*models.Post, error) {
	post := models.Post{}
	err := db.DB.Joins("User").First(&post, id).Preload("Likes").Joins("LEFT JOIN likes ON (posts.ID = likes.post_id AND likes.user_id = ?)", userID).Error
	return &post, err
}

func (db *PostDB) UpdatePost(post *models.Post, postid string, userID string) (*models.Post, error) {
	postGet, err := db.GetPostByID(postid, userID)
	if err != nil {
		return post, err
	}
	if err := checkUserIsOwner(postGet, userID); err != nil {
		return post, err
	}
	resPost := &models.Post{}
	result := db.DB.Model(resPost).Clauses(clause.Returning{}).Where("id = ?", postid).Updates(post).Preload("Likes").Joins("LEFT JOIN likes ON (posts.ID = likes.post_id AND likes.user_id = ?)", userID)
	err = result.Error
	resPost.User = postGet.User
	return resPost, err
}

func checkUserIsOwner(post *models.Post, userID string) error {
	if post.UserID != userID {
		return ErrNotOwner
	}
	return nil
}

type PostViewer interface {
	GetPost() *models.Post
}
