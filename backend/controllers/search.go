package controllers

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/helpers"
)

func (a *APIEnv) GetSearchResults(ctx *gin.Context) {
	searchTerm := ctx.Query("q")
	limit := ctx.DefaultQuery("limit", "10")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, err)
		return
	}

	userResults, err := a.UserDBHandler.QueryUser(searchTerm, limitInt)
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, err)
		return
	}

	results := userResults

	projectResults, err := a.ProjectDBHandler.QueryProject(searchTerm, limitInt)
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, err)
		return
	}

	results = append(results, projectResults...)
	communityResults, err := a.CommunityDBHandler.QueryCommunity(searchTerm, limitInt)
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, err)
		return
	}
	results = append(results, communityResults...)

	// Sort the results by score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// Return the top 10 results
	if len(results) > 10 {
		results = results[:10]
	}

	helpers.OutputData(ctx, results)
}
