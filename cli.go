package main

import (
	"best-route/database"
	"best-route/dijkstra"
	"fmt"
	"strings"
)

func RunCLI(db *database.Database, djk *dijkstra.Dijkstra) error {
	for {
		var route string

		fmt.Println("please enter the route:")
		fmt.Scanf("%s", &route)

		routes, err := db.Client.GetAllRoutes()
		if err != nil {
			return err
		}

		start, target := getBestRouteRequest(route)

		res, err := djk.Client.BestRoute(start, target, routes)
		if err != nil {
			return err
			// TODO: not found in graph
		}

		fmt.Printf("best route: %v > $%v\n", res.Route, res.Cost)
	}
}

func getBestRouteRequest(route string) (string, string) {
	return strings.Split(route, "-")[0], strings.Split(route, "-")[1]
}
