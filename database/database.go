package database

import (
	"best-route/models"
)

type Database struct {
	Client DatabaseClient
}

type DatabaseClient interface {
	GetAllRoutes() ([]*models.Route, error)
	InsertOneRoute(route *models.Route) (*models.Route, error)
}
