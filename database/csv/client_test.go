package csv

import (
	"best-route/models"
	"encoding/csv"
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

const PathTest = "../../input-routes_test.csv"

var (
	client *CsvClient
	seedRoutes map[string]*models.Route
)

func SeedRoutesToTest() []*models.Route {
	seedRoutes = map[string]*models.Route{
		"GRU-BRC": {
			Start:  "GRU",
			Target: "BRC",
			Cost:   10,
		},
		"BRC-SCL": {
			Start:  "BRC",
			Target: "SCL",
			Cost:   5,
		},
		"GRU-CDG": {
			Start:  "GRU",
			Target: "CDG",
			Cost:   75,
		},
		"GRU-SCL": {
			Start:  "GRU",
			Target: "SCL",
			Cost:   20,
		},
		"GRU-ORL": {
			Start:  "GRU",
			Target: "ORL",
			Cost:   56,
		},
		"ORL-CDG": {
			Start:  "ORL",
			Target: "CDG",
			Cost:   5,
		},
		"SCL-ORL": {
			Start:  "SCL",
			Target: "ORL",
			Cost:   20,
		},
	}

	// Convert map to slice of values.
	var routes []*models.Route
	for _, r := range seedRoutes {
		routes = append(routes, r)
	}

	return routes
}

func TestMain(m *testing.M){
	os.Exit(func() int {
		routes := SeedRoutesToTest()

		c, err := NewMockCsvClient(&CsvClient{
			Path:   PathTest,
			Routes: routes,
		})
		if err != nil {
			return 0
		}

		client = c
		return m.Run()
	}())
}

func TestClient_InsertOne(t *testing.T) {
	assert := assert2.New(t)

	t.Run("already exists", func(t *testing.T) {
		route := seedRoutes["GRU-BRC"]

		r, err := client.InsertOne(route)

		assert.Nil(r)

		if assert.Error(err) {
			msg := err.Error()
			assert.Equal(msg, "already exists")
		}
	})

	t.Run("success", func(t *testing.T) {
		route := &models.Route{
			Start:  "GO",
			Target: "TO",
			Cost:   100,
		}

		r, err := client.InsertOne(route)

		assert.NoError(err)

		if assert.NotNil(r) {
			assert.Equal(r.Start, route.Start)
			assert.Equal(r.Target, route.Target)
			assert.Equal(r.Cost, route.Cost)
		}
	})
}

func TestClient_CsvReaderToRoutes(t *testing.T) {
	assert := assert2.New(t)

	t.Run("success", func(t *testing.T) {
		str := `GRU,BRC,10
BRC,SCL,5
GRU,CDG,75
GRU,SCL,20
GRU,ORL,56
ORL,CDG,5
SCL,ORL,20
`
		reader := csv.NewReader(strings.NewReader(str))

		res, err := CsvReaderToRoutes(reader)

		assert.NoError(err)

		if assert.NotNil(res) {
			res[0].Equal(seedRoutes["GRU-BRC"])
		}
	})
}