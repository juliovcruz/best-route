package models

import (
	"fmt"

	"github.com/juliovcruz/dijkstra"
)

type Route struct {
	Start  string
	Target string
	Cost   int
}

func (r *Route) ToString() []string {
	return []string{r.Start, r.Target, fmt.Sprintf("%v", r.Cost)}
}

func (r *Route) Equal(route *Route) bool {
	if route.Start == r.Start &&
		route.Target == r.Target &&
		route.Cost == r.Cost {
		return true
	}

	return false
}

func ConvertManyRoutesToGraph(routes []*Route) dijkstra.Graph {
	graph := dijkstra.Graph{}

	for _, start := range routes {
		targets := GetMapTargetsByStart(routes, start.Start)

		graph[start.Start] = targets
	}

	return graph
}

func GetMapTargetsByStart(routes []*Route, start string) map[string]int {
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
