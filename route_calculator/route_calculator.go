package route_calculator

import (
	"best-route/models"
)

type Router struct {
	Client Client
}

type BestRouteResponse struct {
	Route []string `json:"route"`
	Cost  int      `json:"cost"`
}

type Client interface {
	BestRoute(start string, target string, routes []*models.Route) (*BestRouteResponse, *models.ResponseError)
}
