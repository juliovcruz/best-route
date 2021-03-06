package main

import (
	"best-route/database"
	"best-route/database/csv"
	"log"
	"testing"

	assert2 "github.com/stretchr/testify/assert"
)

func TestCLI(t *testing.T) {
	assert := assert2.New(t)

	t.Run("success", func(t *testing.T) {
		csvClient, err := csv.NewCsvClient("./input-routes.csv")
		if err != nil {
			log.Fatal(err)
		}

		err = RunCLI(&database.Database{
			Client: csvClient,
		})

		assert.NoError(err)
	})
}
