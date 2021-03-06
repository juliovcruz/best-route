package dijkstra

import (
	"best-route/models"
)

type Dijkstra struct {
	Client DijkstraClient
}

type BestRouteResponse struct {
	Route []string `json:"route,omitempty"`
	Cost  int      `json:"cost,omitempty"`
}

type DijkstraClient interface {
	BestRoute(start string, target string, routes []*models.Route) (*BestRouteResponse, error)
}
