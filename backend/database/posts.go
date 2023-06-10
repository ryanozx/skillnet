package database

import (
	"errors"

	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
)

var ErrNotOwner = errors.New("unauthorised action")

func GetPosts(db *gorm.DB) ([]models.PostView, error) {
	var posts []models.Post
	var postViews []models.PostView

	// Retrieve all posts from database
	query := db.Preload("User").Find(&posts)
	if err := query.Find(&posts).Error; err != nil {
		return nil, err
	}

	// Fill in user details for each post using userID of post creator
	for _, post := range posts {
		post := bindUserToPost(&post)
		postViews = append(postViews, *post)
	}
	return postViews, nil
}

func CreatePost(db *gorm.DB, post *models.Post) (*models.PostView, error) {
	result := db.Create(post)
	if err := result.Error; err != nil {
		return nil, err
	}
	newPost := bindUserToPost(post)
	return newPost, nil
}

func GetPostByID(db *gorm.DB, id string) (*models.PostView, error) {
	post := models.Post{}
	err := db.First(&post, id).Error
	if err != nil {
		return nil, err
	}
	postView := bindUserToPost(&post)
	return postView, nil
}

func UpdatePost(db *gorm.DB, post *models.Post, id string) (*models.PostView, error) {
	originalPostView, err := GetPostByID(db, id)
	if err != nil {
		return nil, err
	}
	originalPost := originalPostView.Post
	userID := originalPost.UserID.String()
	if err := checkUserIsOwner(originalPostView, userID); err != nil {
		return nil, err
	}
	err = db.Model(originalPost).Updates(post).Error
	if err != nil {
		return nil, err
	}
	postView := bindUserToPost(&originalPost)
	return postView, nil
}

func DeletePost(db *gorm.DB, id string, userID string) error {
	postView, err := GetPostByID(db, id)
	if err != nil {
		return err
	}
	if err := checkUserIsOwner(postView, userID); err != nil {
		return err
	}
	post := postView.Post
	err = db.Delete(&post).Error
	return err
}

func bindUserToPost(post *models.Post) *models.PostView {
	postView := post.PostView()
	return postView
}

func checkUserIsOwner(postView *models.PostView, userID string) error {
	if postView.Post.UserID.String() != userID {
		return ErrNotOwner
	}
	return nil
}
