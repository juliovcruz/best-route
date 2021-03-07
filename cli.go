package main

import (
	"best-route/database"
	"best-route/router"
	"fmt"
)

func RunCLI(db *database.Database, router *router.Router) error {
	fmt.Printf("| Example input: GRU-CDG -- To exit press CTRL+C |\n")

	for {
		var str string

		fmt.Println("please enter the route:")
		fmt.Scanf("%s", &str)

		route, resErr := validateCLIBestRouteRequest(str)
		if resErr != nil {
			fmt.Printf("precondition failed: %v\n", resErr.ToString())
			continue
		}

		routes, err := db.Client.GetAllRoutes()
		if err != nil {
			fmt.Printf("internal error: %v\n", err.Error())
			continue
		}

		res, resErr := router.Client.BestRoute(route.Start, route.Target, routes)
		if resErr != nil {
			fmt.Printf("precondition failed: %v\n", resErr.ToString())
			continue
		}

		fmt.Printf("best route: %v > $%v\n", res.Route, res.Cost)
	}
}
