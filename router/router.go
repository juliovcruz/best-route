package router

import (
	"best-route/models"
)

type Router struct {
	Client RouterClient
}

type BestRouteResponse struct {
	Route []string `json:"route,omitempty"`
	Cost  int      `json:"cost,omitempty"`
}

type RouterClient interface {
	BestRoute(start string, target string, routes []*models.Route) (*BestRouteResponse, *models.ResponseError)
}
