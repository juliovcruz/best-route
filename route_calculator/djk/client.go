package djk

import (
	"best-route/models"
	"best-route/route_calculator"
	"fmt"

	djkhub "github.com/juliovcruz/dijkstra"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) BestRoute(start string, target string, routes []*models.Route) (*route_calculator.BestRouteResponse, *models.ResponseError) {
	if len(routes) == 0 {
		return nil, &models.ResponseError{
			Message: "no have routes",
		}
	}

	graph := ConvertManyRoutesToGraph(routes)

	if resErr := ensurePlacesArePart(start, target, routes); resErr != nil {
		return nil, resErr
	}

	if start == target {
		return &route_calculator.BestRouteResponse{
			Route: []string{start},
			Cost:  0,
		}, nil
	}

	bestRoute, cost, _ := graph.Path(start, target)

	if len(bestRoute) < 2 {
		return nil, &models.ResponseError{
			Message: fmt.Sprintf("no have route of %v to %v", start, target),
		}
	}

	return &route_calculator.BestRouteResponse{
		Route: bestRoute,
		Cost:  cost,
	}, nil
}

func ConvertManyRoutesToGraph(routes []*models.Route) djkhub.Graph {
	graph := djkhub.Graph{}

	for _, start := range routes {
		targets := GetMapTargetsByStart(routes, start.Start)

		graph[start.Start] = targets
	}

	return graph
}

func GetMapTargetsByStart(routes []*models.Route, start string) map[string]int {
	routesRes := make(map[string]int)

	for _, r := range routes {
		if r.Start == start {
			if value, found := routesRes[r.Target]; found {
				if value < r.Cost {
					continue
				}
			}

			routesRes[r.Target] = r.Cost
		}
	}

	return routesRes
}

func ensurePlacesArePart(start string, target string, routes []*models.Route) *models.ResponseError {
	if !ensurePlaceArePart(start, routes) {
		return &models.ResponseError{
			Field:   "start",
			Message: fmt.Sprintf("cannot find %v in registered routes", start),
		}
	}

	if !ensurePlaceArePart(target, routes) {
		return &models.ResponseError{
			Field:   "target",
			Message: fmt.Sprintf("cannot find %v in registered routes", target),
		}
	}

	return nil
}

func ensurePlaceArePart(place string, routes []*models.Route) bool {
	for _, route := range routes {
		if route.Start == place || route.Target == place {
			return true
		}
	}

	return false
}
