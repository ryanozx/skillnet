package database

import (
	"errors"

	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
)

var ErrNotOwner = errors.New("unauthorised action")

type PostDBHandler interface {
	CreatePost(*models.Post) (*models.PostView, error)
	DeletePost(string, string) error
	GetPosts() ([]models.PostView, error)
	GetPostByID(string) (*models.PostView, error)
	UpdatePost(*models.Post, string, string) (*models.PostView, error)
}

// PostDB implements PostDBHandler
type PostDB struct {
	DB *gorm.DB
}

func (db *PostDB) CreatePost(post *models.Post) (*models.PostView, error) {
	result := db.DB.Create(post)
	newPostView := post.PostView()
	return newPostView, result.Error
}

func (db *PostDB) DeletePost(id string, userID string) error {
	postView, err := db.GetPostByID(id)
	if err != nil {
		return err
	}
	if err := checkUserIsOwner(postView, userID); err != nil {
		return err
	}
	post := postView.GetPost()
	err = db.DB.Delete(&post).Error
	return err
}

func (db *PostDB) GetPosts() ([]models.PostView, error) {
	var posts []models.Post
	var postViews []models.PostView

	// Retrieve all posts from database
	query := db.DB.Joins("User").Find(&posts)
	if err := query.Find(&posts).Error; err != nil {
		return postViews, err
	}

	// Fill in user details for each post using userID of post creator
	for _, post := range posts {
		post := post.PostView()
		postViews = append(postViews, *post)
	}
	return postViews, nil
}

func (db *PostDB) GetPostByID(id string) (*models.PostView, error) {
	post := models.Post{}
	err := db.DB.First(&post, id).Error
	postView := post.PostView()
	return postView, err
}

func (db *PostDB) UpdatePost(post *models.Post, postid string, userID string) (*models.PostView, error) {
	originalPostView, err := db.GetPostByID(postid)
	if err != nil {
		return originalPostView, err
	}
	originalPost := originalPostView.GetPost()
	if err := checkUserIsOwner(originalPostView, userID); err != nil {
		return originalPostView, err
	}
	err = db.DB.Model(originalPost).Updates(post).Error
	postView := post.PostView()
	return postView, err
}

func checkUserIsOwner(postView PostViewer, userID string) error {
	if postView.GetPost().UserID != userID {
		return ErrNotOwner
	}
	return nil
}

type PostViewer interface {
	GetPost() *models.Post
}
