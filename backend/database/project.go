package database

import (
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const projectsToReturn = 10

type ProjectDBHandler interface {
	CreateProject(*models.Project) (*models.Project, error)
	DeleteProject(uint, string) error
	GetProjectByID(uint) (*models.Project, error)
	GetProjects(cutoff *helpers.NullableUint, communityID *helpers.NullableUint, username string) ([]models.Project, error)
	UpdateProject(*models.Project, uint, string) (*models.Project, error)
}

type ProjectDB struct {
	DB *gorm.DB
}

func (db *ProjectDB) CreateProject(project *models.Project) (*models.Project, error) {
	result := db.DB.Create(project)
	if result.Error != nil {
		return project, result.Error
	}
	return db.GetProjectByID(project.ID)
}

func (db *ProjectDB) DeleteProject(projectID uint, userID string) error {
	project, err := db.GetProjectByID(projectID)
	if err != nil {
		return err
	}
	if err := helpers.CheckUserIsOwner(project, userID); err != nil {
		return err
	}
	err = db.DB.Delete(&project).Error
	return err
}

func (db *ProjectDB) GetProjects(cutoff *helpers.NullableUint, communityID *helpers.NullableUint, username string) ([]models.Project, error) {
	var projects []models.Project

	query := db.DB

	if !cutoff.IsNull() {
		cutoffVal, _ := cutoff.GetValue()
		query = query.Where("projects.id < ?", cutoffVal)
	}

	query = query.Joins("User")

	if !communityID.IsNull() {
		communityIDVal, _ := communityID.GetValue()
		query = query.Where("projects.community_id = ?", communityIDVal).Joins("Community")
	} else if username != "" {
		query = query.Where("users.username = ?", username)
	}

	query = query.Order("projects.id desc").Limit(projectsToReturn).Find(&projects)
	return projects, query.Error
}

func (db *ProjectDB) GetProjectByID(projectID uint) (*models.Project, error) {
	project := models.Project{}
	err := db.DB.Joins("User").First(&project, "projects.id = ?", projectID).Error
	return &project, err
}

func (db *ProjectDB) UpdateProject(project *models.Project, projectID uint, userID string) (*models.Project, error) {
	projectGet, err := db.GetProjectByID(projectID)
	if err != nil {
		return project, err
	}
	if err := helpers.CheckUserIsOwner(projectGet, userID); err != nil {
		return project, err
	}
	resProject := &models.Project{}
	result := db.DB.Model(resProject).Clauses(clause.Returning{}).Where("id = ?", projectID).Updates(project)
	err = result.Error
	resProject.User = projectGet.User
	return resProject, err
}
