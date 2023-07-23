/*
Contains controllers for Community API.
*/
package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
)

// Errors
var (
	ErrCannotCreateCommunity = errors.New("cannot create community")
	ErrCannotUpdateCommunity = errors.New("cannot update community")
	ErrCommunityNotFound     = errors.New("community not found")
)

func (a *APIEnv) InitialiseCommunityHandler() {
	a.CommunityDBHandler = &database.CommunityDB{
		DB: a.DB,
	}
}

func (a *APIEnv) CreateCommunity(ctx *gin.Context) {
	var newCommunity models.Community

	// If unable to bind JSON in request to the Community object, return status
	// code 400 Bad Request
	if err := helpers.BindInput(ctx, &newCommunity); err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrBadBinding)
		return
	}

	// Add userID into the corresponding field in the newCommunity object so that
	// the client does not have to pass in any userID, and overwrites any userID
	// that a malicious client might have passed in.
	userID := helpers.GetUserIDFromContext(ctx)
	newCommunity.OwnerID = userID

	community, err := a.CommunityDBHandler.CreateCommunity(&newCommunity)

	// If community cannot be created, return status code 500 Internal Service Error
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrCannotCreateCommunity)
		return
	}

	helpers.OutputData(ctx, community.CommunityView(userID))
}

func (a *APIEnv) GetCommunities(ctx *gin.Context) {
	userID := helpers.GetUserIDFromContext(ctx)
	// Ensure that cutoff is an unsigned integer or empty
	cutoff, err := helpers.GetCutoffFromQuery(ctx)
	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrBadBinding)
		return
	}

	communities, err := a.CommunityDBHandler.GetCommunities(cutoff)
	// If unable to retrieve communities, return status code 404 Not Found
	if err != nil {
		helpers.OutputError(ctx, http.StatusNotFound, ErrCommunityNotFound)
		return
	}

	var smallestID uint = 0
	var communityViews []models.CommunityView
	// Set next cutoff value
	for _, community := range communities {
		smallestID = community.ID
		communityView := community.CommunityView(userID)
		communityViews = append(communityViews, *communityView)
	}

	communitiesArray := models.CommunityArray{
		Communities: communityViews,
		NextPageURL: helpers.GenerateCommunitiesNextPageURL(models.BackendAddress, smallestID),
	}
	helpers.OutputData(ctx, communitiesArray)
}

func (a *APIEnv) GetCommunityByName(ctx *gin.Context) {
	userID := helpers.GetUserIDFromContext(ctx)

	communityName := helpers.GetCommunityNameFromContext(ctx)

	community, err := a.CommunityDBHandler.GetCommunityByName(communityName)
	// If unable to retrieve community, return status code 404 Not Found
	if err != nil {
		helpers.OutputError(ctx, http.StatusNotFound, ErrCommunityNotFound)
		return
	}

	communityView := community.CommunityView(userID)
	helpers.OutputData(ctx, communityView)
}

func (a *APIEnv) UpdateCommunity(ctx *gin.Context) {
	userID := helpers.GetUserIDFromContext(ctx)

	communityName := helpers.GetCommunityNameFromContext(ctx)

	var inputUpdate models.Community

	// If unable to bind JSON in request to the Community object, return status
	// code 400 Bad Request
	if err := helpers.BindInput(ctx, &inputUpdate); err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, ErrBadBinding)
		return
	}

	community, err := a.CommunityDBHandler.UpdateCommunity(&inputUpdate, communityName, userID)

	// If community cannot be found in the database, return status code 404 Status Not Found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		helpers.OutputError(ctx, http.StatusNotFound, ErrCommunityNotFound)
		return
	}
	// If user is not the owner of the community, return status code 403 Forbidden
	if errors.Is(err, helpers.ErrNotOwner) {
		helpers.OutputError(ctx, http.StatusForbidden, helpers.ErrNotOwner)
		return
	}
	// If community cannot be updated for any other reason, return status code 500 Internal Server Error
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, ErrCannotUpdateCommunity)
		return
	}
	helpers.OutputData(ctx, community.CommunityView(userID))
}
