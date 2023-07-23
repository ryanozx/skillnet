package database

import (
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const communitiesToReturn = 10

type CommunityDBHandler interface {
	CreateCommunity(*models.Community) (*models.Community, error)
	GetCommunityByID(uint) (*models.Community, error)
	GetCommunityByName(name string) (*models.Community, error)
	GetCommunities(*helpers.NullableUint) ([]models.Community, error)
	UpdateCommunity(*models.Community, string, string) (*models.Community, error)
}

type CommunityDB struct {
	DB *gorm.DB
}

func (db *CommunityDB) CreateCommunity(community *models.Community) (*models.Community, error) {
	result := db.DB.Create(community)
	if result.Error != nil {
		return community, result.Error
	}
	return db.GetCommunityByID(community.ID)
}

func (db *CommunityDB) GetCommunities(cutoff *helpers.NullableUint) ([]models.Community, error) {
	var communities []models.Community

	query := db.DB

	if !cutoff.IsNull() {
		cutoffVal, _ := cutoff.GetValue()
		query = query.Where("communities.id < ?", cutoffVal)
	}

	query = query.Joins("User").Order("communities.id desc").Limit(communitiesToReturn).Find(&communities)
	return communities, query.Error
}

func (db *CommunityDB) GetCommunityByID(communityID uint) (*models.Community, error) {
	community := models.Community{}
	err := db.DB.Joins("User").First(&community, "communities.id = ?", communityID).Error
	return &community, err
}

func (db *CommunityDB) GetCommunityByName(communityName string) (*models.Community, error) {
	community := models.Community{}
	err := db.DB.Joins("User").First(&community, "communities.name = ?", communityName).Error
	return &community, err
}

func (db *CommunityDB) UpdateCommunity(community *models.Community, communityName string, userID string) (*models.Community, error) {
	communityGet, err := db.GetCommunityByName(communityName)
	if err != nil {
		return community, err
	}
	if err := helpers.CheckUserIsOwner(communityGet, userID); err != nil {
		return community, err
	}
	resCommunity := &models.Community{}
	result := db.DB.Model(resCommunity).Clauses(clause.Returning{}).Where("id = ?", communityGet.ID).Updates(map[string]interface{}{
		"About": community.About,
	})
	err = result.Error
	resCommunity.User = communityGet.User
	return resCommunity, err
}
