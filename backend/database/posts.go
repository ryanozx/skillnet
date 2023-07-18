package database

import (
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const postsToReturn = 10

type PostDBHandler interface {
	CreatePost(*models.Post) (*models.Post, error)
	DeletePost(uint, string) error
	GetPosts(cutoff *helpers.NullableUint, communityID *helpers.NullableUint, projectID *helpers.NullableUint, userID string) ([]models.Post, error)
	GetPostByID(uint, string) (*models.Post, error)
	UpdatePost(*models.Post, uint, string) (*models.Post, error)
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
	return db.GetPostByID(post.ID, post.UserID)
}

func (db *PostDB) DeletePost(postID uint, userID string) error {
	post, err := db.GetPostByID(postID, userID)
	if err != nil {
		return err
	}
	if err := helpers.CheckUserIsOwner(post, userID); err != nil {
		return err
	}
	err = db.DB.Delete(&post).Error
	return err
}

func (db *PostDB) GetPosts(cutoff *helpers.NullableUint, communityID *helpers.NullableUint,
	projectID *helpers.NullableUint, userID string) ([]models.Post, error) {
	var posts []models.Post

	var tx *gorm.DB

	if !projectID.IsNull() {
		// Check if we should filter for project (e.g. project feed)
		projectIDVal, _ := projectID.GetValue()
		tx = db.DB.Where("posts.project_id = ?", projectIDVal)
	} else if !communityID.IsNull() {
		// Check if we should filter for community (e.g. community feed)
		communityIDVal, _ := communityID.GetValue()
		tx = db.DB.Where("posts.community_id = ?", communityIDVal)
	} else {
		// Global feed
		tx = db.DB
	}

	if !cutoff.IsNull() {
		cutoffVal, _ := cutoff.GetValue()
		tx = tx.Where("posts.id < ?", cutoffVal)
	}

	query := tx.Joins("User").Limit(postsToReturn).Preload("Likes").
		Joins("LEFT JOIN likes ON (posts.ID = likes.post_id AND likes.user_id = ?)", userID).
		Order("posts.id desc").Find(&posts)

	if err := query.Find(&posts).Error; err != nil {
		return posts, err
	}
	return posts, nil
}

func (db *PostDB) GetPostByID(postID uint, userID string) (*models.Post, error) {
	post := models.Post{}
	var err error
	if userID == "" {
		err = db.DB.Joins("User").First(&post, postID).Error
	} else {
		err = db.DB.Joins("User").First(&post, postID).Preload("Likes").Joins("LEFT JOIN likes ON (posts.ID = likes.post_id AND likes.user_id = ?)", userID).Error
	}
	return &post, err
}

func (db *PostDB) UpdatePost(post *models.Post, postID uint, userID string) (*models.Post, error) {
	postGet, err := db.GetPostByID(postID, userID)
	if err != nil {
		return post, err
	}
	if err := helpers.CheckUserIsOwner(postGet, userID); err != nil {
		return post, err
	}
	resPost := &models.Post{}
	result := db.DB.Model(resPost).Clauses(clause.Returning{}).Where("id = ?", postID).Updates(post).Preload("Likes").Joins("LEFT JOIN likes ON (posts.ID = likes.post_id AND likes.user_id = ?)", userID)
	err = result.Error
	resPost.User = postGet.User
	return resPost, err
}

type PostViewer interface {
	GetPost() *models.Post
}
