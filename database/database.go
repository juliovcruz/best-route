package database

import (
	"best-route/models"
)

type Client interface {
	InsertOne(route *models.Route) (*models.Route, error)
}
