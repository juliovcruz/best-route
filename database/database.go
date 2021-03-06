package database

import (
	"best-route/models"
)

type DatabaseClient interface {
	InsertOne(route *models.Route) (*models.Route, error)
}
