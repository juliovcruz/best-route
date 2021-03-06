package main

import (
	"best-route/database"
	"best-route/models"
	"fmt"
	"strings"
)

func RunCLI(db *database.Database) error {
	for {
		var route string

		fmt.Println("please enter the route:")
		fmt.Scanf("%s", &route)

		routes, err := db.Client.GetAllRoutes()
		if err != nil {
			return err
		}

		graph := models.ConvertManyRoutesToGraph(routes)

		bestRoute, cost, err := graph.Path(getBestRouteRequest(route))
		if err != nil {
			return err
			// TODO: not found in graph
		}

		fmt.Printf("best route: %v > $%v", bestRoute, cost)
	}
}

func getBestRouteRequest(route string) (string, string) {
	return strings.Split(route, "-")[0], strings.Split(route, "-")[1]
}
