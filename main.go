package main

import (
	"best-route/database"
	"best-route/database/csv"
	"best-route/router"
	"best-route/router/djk"
	"log"
	"os"
)

const (
	path             = "input-routes.csv"
	MaxSizePlaceName = 30
)

func main() {
	if len(os.Args) > 1 {
		csvClient, err := csv.NewCsvClient("./" + os.Args[1])
		if err != nil {
			log.Fatal(err)
		}

		if err := RunCLI(
			&database.Database{Client: csvClient},
			&router.Router{Client: djk.NewDjkClient()},
		); err != nil {
			log.Fatal(err)
		}
		return
	}

	csvClient, err := csv.NewCsvClient("./" + path)
	if err != nil {
		log.Fatal(err)
	}

	if err := RunAPI(
		&database.Database{Client: csvClient},
		&router.Router{Client: djk.NewDjkClient()},
		":3000",
	); err != nil {
		log.Fatal(err)
	}
}
