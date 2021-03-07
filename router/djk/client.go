package djk

import (
	"best-route/models"
	"best-route/router"
	"fmt"
	"strings"

	djk_hub "github.com/juliovcruz/dijkstra"
)

type DjkClient struct{}

func NewDjkClient() *DjkClient {
	return &DjkClient{}
}

func (c *DjkClient) BestRoute(start string, target string, routes []*models.Route) (*router.BestRouteResponse, *models.ResponseError) {
	if len(routes) == 0 {
		return nil, &models.ResponseError{
			Message: "no have routes",
		}
	}

	if start == target {
		return &router.BestRouteResponse{
			Route: []string{start},
			Cost:  0,
		}, nil
	}

	graph := ConvertManyRoutesToGraph(routes)

	bestRoute, cost, err := graph.Path(start, target)
	if err != nil {
		return nil, convertErrorToResponseError(err)
	}

	if len(bestRoute) < 2 {
		return nil, &models.ResponseError{
			Message: fmt.Sprintf("no have route of %v to %v", start, target),
		}
	}

	return &router.BestRouteResponse{
		Route: bestRoute,
		Cost:  cost,
	}, nil
}

func ConvertManyRoutesToGraph(routes []*models.Route) djk_hub.Graph {
	graph := djk_hub.Graph{}

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

func convertErrorToResponseError(err error) *models.ResponseError {
	if strings.Contains(err.Error(), "start") {
		return &models.ResponseError{
			Field:   "start",
			Message: err.Error(),
		}
	}

	if strings.Contains(err.Error(), "target") {
		return &models.ResponseError{
			Field:   "target",
			Message: err.Error(),
		}
	}

	return nil
}
