package main

import (
	"best-route/database"
	"best-route/models"
	"best-route/route_calculator"
	"fmt"
	"time"
)

func RunCLI(db *database.Database, router *route_calculator.Router) error {
	time.Sleep(time.Millisecond * 5)

	fmt.Printf("example input in CLI: GRU-CDG -- to exit press CTRL+C\n")

	for {
		var str string

		fmt.Println("please enter the route:")
		fmt.Scanf("%s", &str)

		route, resErr := models.ValidateCLIBestRouteRequest(str)
		if resErr != nil {
			fmt.Printf("precondition failed: %v\n", resErr.ToString())
			continue
		}
		route = route.Trim()

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
