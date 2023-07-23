package database

import (
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const commentsToReturn = 10

type CommentsDBHandler interface {
	CreateComment(*models.Comment) (*models.Comment, error)
	DeleteComment(uint, string) (uint, error)
	GetComments(uint, *helpers.NullableUint) ([]models.Comment, error)
	GetCommentByID(uint) (*models.Comment, error)
	UpdateComment(*models.Comment, uint, string) (*models.Comment, error)
	GetValue(uint) (uint64, error)
}

// CommentDB implements CommentDBHandler
type CommentDB struct {
	DB *gorm.DB
}

func (db *CommentDB) CreateComment(comment *models.Comment) (*models.Comment, error) {
	result := db.DB.Create(comment)
	if result.Error != nil {
		return comment, result.Error
	}
	return db.GetCommentByID(comment.ID)
}

func (db *CommentDB) DeleteComment(commentID uint, userID string) (uint, error) {
	comment, err := db.GetCommentByID(commentID)
	if err != nil {
		return comment.PostID, err
	}
	if err := helpers.CheckUserIsOwner(comment, userID); err != nil {
		return comment.PostID, err
	}
	err = db.DB.Delete(&comment).Error
	return comment.PostID, err
}

func (db *CommentDB) GetComments(postID uint, cutoff *helpers.NullableUint) ([]models.Comment, error) {
	var comments []models.Comment

	query := db.DB.Where("comments.post_id = ?", postID)

	if !cutoff.IsNull() {
		cutoffVal, _ := cutoff.GetValue()
		query = query.Where("comments.id < ?", cutoffVal)
	}

	query = query.Joins("User").Order("comments.id desc").Limit(commentsToReturn).Find(&comments)
	return comments, query.Error
}

func (db *CommentDB) GetCommentByID(commentID uint) (*models.Comment, error) {
	comment := models.Comment{}
	err := db.DB.Joins("Post").Joins("User").First(&comment, "comments.id = ?", commentID).Error
	return &comment, err
}

func (db *CommentDB) UpdateComment(comment *models.Comment, commentID uint, userID string) (*models.Comment, error) {
	commentGet, err := db.GetCommentByID(commentID)
	if err != nil {
		return comment, err
	}
	if err := helpers.CheckUserIsOwner(commentGet, userID); err != nil {
		return comment, err
	}
	resComment := &models.Comment{}
	result := db.DB.Model(resComment).Clauses(clause.Returning{}).Where("id = ?", commentID).Updates(comment)
	err = result.Error
	resComment.User = commentGet.User
	return resComment, err
}

func (db *CommentDB) GetValue(postID uint) (uint64, error) {
	var count int64
	result := db.DB.Model(&models.Comment{}).Where("post_id = ?", postID).Count(&count)
	return uint64(count), result.Error
}
