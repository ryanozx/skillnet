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
	if result.Error != nil {
		return post.PostView(post.UserID), result.Error
	}
	return db.GetPostByID(fmt.Sprintf("%v", post.ID), post.UserID)
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
		query := db.DB.Limit(postsToReturn).Joins("User").Preload("Likes").Joins("LEFT JOIN likes ON (posts.ID = likes.post_id AND posts.user_id = likes.user_id)").Order("posts.id desc").Find(&posts)
		if err := query.Find(&posts).Error; err != nil {
			return postViews, 0, err
		}
	} else {
		query := db.DB.Where("posts.id < ?", cutoff).Joins("User").Limit(postsToReturn).Preload("Likes").Joins("LEFT JOIN likes ON (posts.ID = likes.post_id AND posts.user_id = likes.user_id)").Order("posts.id desc").Find(&posts)
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
	err := db.DB.Joins("User").First(&post, id).Preload("Likes").Joins("LEFT JOIN likes ON (posts.ID = likes.post_id AND posts.user_id = likes.user_id)").Error
	postView := post.PostView(userID)
	return postView, err
}

func (db *PostDB) UpdatePost(post *models.Post, postid string, userID string) (*models.PostView, error) {
	postView, err := db.GetPostByID(postid, userID)
	if err != nil {
		return postView, err
	}
	if err := checkUserIsOwner(postView, userID); err != nil {
		return postView, err
	}
	resPost := &models.Post{}
	result := db.DB.Model(resPost).Clauses(clause.Returning{}).Where("id = ?", postid).Updates(post).Preload("Likes").Joins("LEFT JOIN likes ON (posts.ID = likes.post_id AND posts.user_id = likes.user_id)")
	err = result.Error
	postView.Post = *resPost
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
