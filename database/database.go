package database

import (
	"best-route/models"
)

const AlreadyExistsMessageError = "this route already exists"

type Database struct {
	Client Client
}

type Client interface {
	GetAllRoutes() ([]*models.Route, error)
	InsertOneRoute(route *models.Route) (*models.Route, error)
}
