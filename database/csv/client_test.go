package csv

import (
	"best-route/models"
	"encoding/csv"
	"os"
	"strings"
	"testing"

	assert2 "github.com/stretchr/testify/assert"
)

const PathTest = "test-input-not-edit.csv"

var (
	client     *Client
	seedRoutes map[string]*models.Route
)

func TestClient_NewClient(t *testing.T) {
	assert := assert2.New(t)

	t.Run("already exists", func(t *testing.T) {
		res, err := NewClient("../../" + PathTest)

		assert.NoError(err)

		if assert.NotNil(res) {
			assert.Equal(7, len(res.Routes))
			assert.Equal("../../"+PathTest, res.Path)
		}
	})
}

func TestClient_InsertOne(t *testing.T) {
	assert := assert2.New(t)

	t.Run("already exists", func(t *testing.T) {
		route := seedRoutes["GRU-BRC"]

		r, err := client.InsertOneRoute(route)

		assert.Nil(r)

		if assert.Error(err) {
			msg := err.Error()
			assert.Contains(msg, "already exists")
		}
	})

	t.Run("success", func(t *testing.T) {
		route := &models.Route{
			Start:  "GO",
			Target: "TO",
			Cost:   100,
		}

		r, err := client.InsertOneRoute(route)

		assert.NoError(err)

		if assert.NotNil(r) {
			assert.Equal(route.Start, r.Start)
			assert.Equal(route.Target, r.Target)
			assert.Equal(route.Cost, r.Cost)
		}
	})
}

func TestClient_GetAllRoutes(t *testing.T) {
	assert := assert2.New(t)

	t.Run("success", func(t *testing.T) {
		res, err := client.GetAllRoutes()

		assert.NoError(err)

		if assert.NotNil(res) {
			assert.GreaterOrEqual(len(res), 7)
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

		res, err := csvReaderToRoutes(reader)

		assert.NoError(err)

		if assert.NotNil(res) {
			assert.True(res[0].Equal(seedRoutes["GRU-BRC"]))
		}
	})

	t.Run("success with errors", func(t *testing.T) {
		str := `GRU,BRC,ACA
BRC,SCL,5
GRU,CDG,75
,SCL,20
GRU,ORL,56
ORL,CDG,5
SCL,ORL,20
`
		reader := csv.NewReader(strings.NewReader(str))

		res, err := csvReaderToRoutes(reader)

		assert.NoError(err)

		if assert.NotNil(res) {
			assert.True(res[0].Equal(seedRoutes["BRC-SCL"]))
		}
	})
}

func TestMain(m *testing.M) {
	os.Exit(func() int {
		routes := SeedRoutesToTest()

		c, err := NewMockClient(&Client{
			Path:   "../../" + PathTest,
			Routes: routes,
		})
		if err != nil {
			return 0
		}

		client = c
		return m.Run()
	}())
}

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
