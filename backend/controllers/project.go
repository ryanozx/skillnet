/*
Contains controllers for Community API.
*/
package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
)

// Messages
const (
	ProjectDeletedMsg = "Project successfully deleted"
)

// Errors
var (
	ErrCannotCreateProject = errors.New("cannot create project")
	ErrCannotDeleteProject = errors.New("cannot delete project")
	ErrCannotUpdateProject = errors.New("cannot update project")
	ErrProjectNotFound     = errors.New("project not found")
)

func (a *APIEnv) InitialiseProjectHandler() {
	a.ProjectDBHandler = &database.ProjectDB{
		DB: a.DB,
	}
}

func (a *APIEnv) CreateProject(ctx *gin.Context) {
	var newProject models.Project

	// If unable to bind JSON in request to the Project object, return status
	// code 400 Bad Request
	if err := helpers.BindInput(ctx, &newProject); err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrBadBinding)
		return
	}

	// Add userID into the corresponding field in the newProject object so that
	// the client does not have to pass in any userID, and overwrites any userID
	// that a malicious client might have passed in.
	userID := helpers.GetUserIDFromContext(ctx)
	newProject.OwnerID = userID

	project, err := a.ProjectDBHandler.CreateProject(&newProject)

	// If project cannot be created, return status code 500 Internal Service Error
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrCannotCreateProject)
		return
	}

	helpers.OutputData(ctx, project.ProjectView(userID))
}

func (a *APIEnv) DeleteProject(ctx *gin.Context) {
	userID := helpers.GetUserIDFromContext(ctx)

	// Ensure that projectID is an unsigned integer
	projectID, err := helpers.GetProjectIDFromContext(ctx)
	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrProjectNotFound)
		return
	}

	err = a.ProjectDBHandler.DeleteProject(projectID, userID)
	// If project cannot be found in the database return status code 404 Status Not Found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		helpers.OutputError(ctx, http.StatusNotFound, ErrProjectNotFound)
		return
	}
	// If user is not the owner of the project, return status code 403 Forbidden
	if errors.Is(err, helpers.ErrNotOwner) {
		helpers.OutputError(ctx, http.StatusForbidden, helpers.ErrNotOwner)
		return
	}
	// If project cannot be deleted for any other reason, return status code 500 Internal Server Error
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrCannotDeleteProject)
		return
	}
	helpers.OutputMessage(ctx, ProjectDeletedMsg)
}

func (a *APIEnv) GetProjects(ctx *gin.Context) {
	// Ensure that cutoff is an unsigned integer or empty
	cutoff, err := helpers.GetCutoffFromQuery(ctx)
	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrBadBinding)
		return
	}

	communityID, err := helpers.GetCommunityIDFromQuery(ctx)
	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrCommunityNotFound)
		return
	}

	fmt.Print(communityID)

	username := helpers.GetUsernameFromQuery(ctx)

	projects, err := a.ProjectDBHandler.GetProjects(cutoff, communityID, username)
	// If unable to retrieve projects, return status code 404 Not Found
	if err != nil {
		helpers.OutputError(ctx, http.StatusNotFound, ErrProjectNotFound)
		return
	}

	var smallestID uint = 0
	var projectMinimals []models.ProjectMinimal
	// Set next cutoff value
	for _, project := range projects {
		smallestID = project.ID
		projectMinimal := project.GetProjectMinimal()
		projectMinimals = append(projectMinimals, *projectMinimal)
	}

	projectsArray := models.ProjectsArray{
		Projects:    projectMinimals,
		NextPageURL: helpers.GenerateProjectsNextPageURL(models.BackendAddress, smallestID, communityID, username),
	}
	helpers.OutputData(ctx, projectsArray)
}

func (a *APIEnv) GetProjectByID(ctx *gin.Context) {
	userID := helpers.GetUserIDFromContext(ctx)

	// Ensure that projectID is an unsigned integer
	projectID, err := helpers.GetProjectIDFromContext(ctx)
	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrProjectNotFound)
		return
	}

	project, err := a.ProjectDBHandler.GetProjectByID(projectID)
	// If unable to retrieve project, return status code 404 Not Found
	if err != nil {
		helpers.OutputError(ctx, http.StatusNotFound, ErrProjectNotFound)
		return
	}

	projectView := project.ProjectView(userID)
	helpers.OutputData(ctx, projectView)
}

func (a *APIEnv) UpdateProject(ctx *gin.Context) {
	userID := helpers.GetUserIDFromContext(ctx)

	// Ensure that projectID is an unsigned integer
	projectID, err := helpers.GetProjectIDFromContext(ctx)

	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrProjectNotFound)
		return
	}

	var inputUpdate models.Project

	// If unable to bind JSON in request to the Project object, return status
	// code 400 Bad Request
	if err := helpers.BindInput(ctx, &inputUpdate); err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrBadBinding)
		return
	}

	project, err := a.ProjectDBHandler.UpdateProject(&inputUpdate, projectID, userID)

	// If project cannot be found in the database, return status code 404 Status Not Found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		helpers.OutputError(ctx, http.StatusNotFound, ErrProjectNotFound)
		return
	}
	// If user is not the owner of the project, return status code 403 Forbidden
	if errors.Is(err, helpers.ErrNotOwner) {
		helpers.OutputError(ctx, http.StatusForbidden, helpers.ErrNotOwner)
		return
	}
	// If project cannot be updated for any other reason, return status code 500 Internal Server Error
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrCannotUpdateProject)
		return
	}
	helpers.OutputData(ctx, project.ProjectView(userID))
}
