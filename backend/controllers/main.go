/*
Contains common functions, structs, and constants used by handler functions.

The response body will contain one of three keys:
- data: Returns the data requested by the client (uses outputData)
- error: Returns an error message (uses outputError)
- message: Returns a message accompanied with status code 200 OK (uses outputMessage)
*/
package controllers

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"cloud.google.com/go/storage"
	"github.com/redis/go-redis/v9"
	"github.com/ryanozx/skillnet/database"
	"gorm.io/gorm"
)

// APIEnv is a wrapper for the shared database instance
type APIEnv struct {
	DB                   *gorm.DB
	NotifRedis           *redis.Client
	PostDBHandler        database.PostDBHandler
	UserDBHandler        database.UserDBHandler
	AuthDBHandler        database.AuthDBHandler
	LikeDBHandler        database.LikeAPIHandler
	CommentDBHandler     database.CommentsDBHandler
	CommunityDBHandler   database.CommunityDBHandler
	ProjectDBHandler     database.ProjectDBHandler
	GoogleCloud          *storage.Client
	LikesCacheHandler    CacheHandler
	CommentsCacheHandler CacheHandler
	NotificationPoster   NotificationPoster
}

// General
var (
	ErrBadBinding         = errors.New("invalid request")
	ErrCookieSaveFail     = errors.New("cookie failed to save")
	ErrDBValueFailed      = errors.New("unable to retrieve value")
	ErrNoValidSession     = errors.New("no valid session")
	ErrSessionClearFailed = errors.New("failed to clear session")
	// ErrTest is a test error that can be used to simulate an unexpected error returned by the database helper functions
	ErrTest                   = errors.New("test error")
	ErrUpdateCacheValueFailed = errors.New("failed to update cache value")
)

type CacheHandler interface {
	GetCacheVal(context.Context, uint) (uint64, error)
	SetCacheVal(context.Context, uint) (uint64, error)
}

type Cache struct {
	redisDB   *redis.Client
	DBHandler database.DBValueGetter
}

func (c *Cache) GetCacheVal(ctx context.Context, key uint) (uint64, error) {
	val, err := c.redisDB.Get(ctx, fmt.Sprintf("%v", key)).Result()
	if err == redis.Nil {
		return c.SetCacheVal(ctx, key)
	}
	return strconv.ParseUint(val, 10, 32)
}

func (c *Cache) SetCacheVal(ctx context.Context, id uint) (uint64, error) {
	newVal, err := c.DBHandler.GetValue(id)
	if err != nil {
		return newVal, ErrDBValueFailed
	}
	err = c.redisDB.Set(ctx, fmt.Sprintf("%v", id), newVal, 0).Err()
	if err != nil {
		return newVal, ErrUpdateCacheValueFailed
	}
	return newVal, nil
}
