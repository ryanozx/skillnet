package database

import (
	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
)

func GetPosts(db *gorm.DB) ([]models.Post, error) {
	var posts []models.PostSchema
	query := db.Find(&posts)
	output := models.ConvertToPosts(posts)
	if err := query.Find(&posts).Error; err != nil {
		return output, err
	}
	return output, nil
}

func CreatePost(db *gorm.DB, postInput *models.PostInput) (*models.Post, error) {
	newPostSchema := postInput.ConvertToPostSchema()
	result := db.Create(newPostSchema)
	newPost := newPostSchema.ConvertToPost()
	if err := result.Error; err != nil {
		return newPost, err
	}
	return newPost, nil
}

func GetPostSchemaByID(db *gorm.DB, id string) (*models.PostSchema, error) {
	postSchema := models.PostSchema{}
	err := db.First(&postSchema, id).Error
	if err != nil {
		return &postSchema, err
	}
	return &postSchema, nil
}

func UpdatePost(db *gorm.DB, updateInput *models.PostInput, id string) (*models.PostSchema, error) {
	originalPost, err := GetPostSchemaByID(db, id)
	if err != nil {
		return originalPost, err
	}
	err = db.Model(originalPost).Updates(updateInput).Error
	return originalPost, err
}

func DeletePost(db *gorm.DB, id string) error {
	postSchema, err := GetPostSchemaByID(db, id)
	if err != nil {
		return err
	}
	err = db.Delete(&postSchema).Error
	return err
}
