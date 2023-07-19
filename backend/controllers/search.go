package controllers

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
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

	projectResults := []models.SearchResult{}
	communityResults := []models.SearchResult{}
	results := append(userResults, append(projectResults, communityResults...)...)

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
