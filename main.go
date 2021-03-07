package main

import (
	"best-route/database"
	"best-route/database/csv"
	"best-route/route_calculator"
	"best-route/route_calculator/djk"
	"fmt"
	"log"
	"os"
)

func main() {
	var number int

	if len(os.Args) > 1 {
		csvClient, err := csv.NewClient("./" + os.Args[1])
		if err != nil {
			log.Fatal(err)
		}

		for number != 1 && number != 2 {
			fmt.Println("Choice one option:\n1 -- CLI\n2 -- API")
			fmt.Scanf("%d", &number)
		}

		if number == 1 {
			if err := RunCLI(
				&database.Database{Client: csvClient},
				&route_calculator.Router{Client: djk.NewClient()},
			); err != nil {
				log.Fatal(err)
			}
			return
		}

		if err := RunAPI(
			&database.Database{Client: csvClient},
			&route_calculator.Router{Client: djk.NewClient()},
			":3000",
		); err != nil {
			log.Fatal(err)

		}
	}

	fmt.Println("input initial csv file is necessary\nExample: ./main input-routes.csv")
}
