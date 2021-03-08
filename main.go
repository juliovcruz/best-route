package main

import (
	"best-route/database"
	"best-route/database/csv"
	"best-route/route_calculator"
	"best-route/route_calculator/djk"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

func main() {
	if len(os.Args) > 1 {
		var (
			g errgroup.Group
		)

		csvClient, err := csv.NewClient("./" + os.Args[1])
		if err != nil {
			log.Fatal(err)
		}

		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		}
		port := os.Getenv("API_PORT")

		g.Go(func() error {
			if err := RunAPI(
				&database.Database{Client: csvClient},
				&route_calculator.Router{Client: djk.NewClient()},
				":"+port,
			); err != nil {
				log.Fatal(err)
			}
			return nil
		})

		g.Go(func() error {
			if err := RunCLI(
				&database.Database{Client: csvClient},
				&route_calculator.Router{Client: djk.NewClient()},
			); err != nil {
				log.Fatal(err)
			}
			return nil
		})

		if err := g.Wait(); err != nil {
			return
		}
	}

	fmt.Println("input initial csv file is necessary\nExample: ./main input-routes.csv")
}
