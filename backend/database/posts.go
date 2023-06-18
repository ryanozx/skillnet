package database

import (
	"errors"

	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
)

var ErrNotOwner = errors.New("unauthorised action")

const postsToReturn = 10

type PostDBHandler interface {
	CreatePost(*models.Post) (*models.PostView, error)
	DeletePost(string, string) error
	GetPosts(string, string) ([]models.PostView, uint, error)
	GetPostByID(string, string) (*models.PostView, error)
	UpdatePost(*models.Post, string, string) (*models.PostView, error)
}

// PostDB implements PostDBHandler
type PostDB struct {
	DB *gorm.DB
}

func (db *PostDB) CreatePost(post *models.Post) (*models.PostView, error) {
	result := db.DB.Create(post)
	newPostView := post.PostView(post.UserID)
	return newPostView, result.Error
}

func (db *PostDB) DeletePost(id string, userID string) error {
	postView, err := db.GetPostByID(id, userID)
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

func (db *PostDB) GetPosts(cutoff string, userID string) ([]models.PostView, uint, error) {
	var posts []models.Post
	var postViews []models.PostView

	if cutoff == "" {
		// Retrieve all posts from database
		query := db.DB.Limit(postsToReturn).Joins("User").Order("posts.id desc").Find(&posts)
		if err := query.Find(&posts).Error; err != nil {
			return postViews, 0, err
		}
	} else {
		query := db.DB.Where("posts.id < ?", cutoff).Joins("User").Limit(postsToReturn).Order("posts.id desc").Find(&posts)
		if err := query.Find(&posts).Error; err != nil {
			return postViews, 0, err
		}
	}

	var smallestID uint = 0
	// Fill in user details for each post using userID of post creator
	for _, post := range posts {
		smallestID = post.ID
		post := post.PostView(userID)
		postViews = append(postViews, *post)
	}
	return postViews, smallestID, nil
}

func (db *PostDB) GetPostByID(id string, userID string) (*models.PostView, error) {
	post := models.Post{}
	err := db.DB.First(&post, id).Error
	postView := post.PostView(userID)
	return postView, err
}

func (db *PostDB) UpdatePost(post *models.Post, postid string, userID string) (*models.PostView, error) {
	originalPostView, err := db.GetPostByID(postid, userID)
	if err != nil {
		return originalPostView, err
	}
	originalPost := originalPostView.GetPost()
	if err := checkUserIsOwner(originalPostView, userID); err != nil {
		return originalPostView, err
	}
	err = db.DB.Model(originalPost).Updates(post).Error
	postView := post.PostView(userID)
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
