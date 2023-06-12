package database

import "github.com/ryanozx/skillnet/models"

type AuthDBHandler interface {
	GetUserByUsername(string) (*models.User, error)
}
