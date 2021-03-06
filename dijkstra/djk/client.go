package djk

import (
	"best-route/models"
	"best-route/dijkstra"
	djk_hub "github.com/juliovcruz/dijkstra"
)

type DjkClient struct {}

func NewDjkClient() *DjkClient {
	return &DjkClient{}
}

func (c *DjkClient) BestRoute(start string, target string, routes []*models.Route) (*dijkstra.BestRouteResponse, error) {
	graph := ConvertManyRoutesToGraph(routes)

	bestRoute, cost, err := graph.Path(start, target)
	if err != nil {
		return nil, err
	}

	return &dijkstra.BestRouteResponse{
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