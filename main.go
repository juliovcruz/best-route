package main

import (
	"best-route/database"
	"best-route/database/csv"
	"best-route/route_calculator"
	"best-route/route_calculator/djk"
	"log"
	"os"
)

const (
	defaultPath      = "input-routes.csv"
	MaxSizePlaceName = 30
)

func main() {
	if len(os.Args) > 1 {
		csvClient, err := csv.NewClient("./" + os.Args[1])
		if err != nil {
			log.Fatal(err)
		}

		if err := RunCLI(
			&database.Database{Client: csvClient},
			&route_calculator.Router{Client: djk.NewClient()},
		); err != nil {
			log.Fatal(err)
		}
		return
	}

	csvClient, err := csv.NewClient("./" + defaultPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := RunAPI(
		&database.Database{Client: csvClient},
		&route_calculator.Router{Client: djk.NewClient()},
		":3000",
	); err != nil {
		log.Fatal(err)
	}
}
