package djk

import (
	"best-route/models"
	"os"
	"testing"

	assert2 "github.com/stretchr/testify/assert"
)

var (
	client *DjkClient
	routes []*models.Route
)

func TestClient_BestRoute(t *testing.T) {
	assert := assert2.New(t)

	t.Run("target not found", func(t *testing.T) {
		res, err := client.BestRoute("GRU", "AAA", routes)

		assert.Nil(res)

		if assert.NotNil(err) {
			assert.Contains(err.Message, "cannot find target")
		}
	})

	t.Run("start not found", func(t *testing.T) {
		res, err := client.BestRoute("AAA", "CDG", routes)

		assert.Nil(res)

		if assert.NotNil(err) {
			assert.Contains(err.Message, "cannot find start")
		}
	})

	t.Run("no have route", func(t *testing.T) {
		res, err := client.BestRoute("ORL", "SCL", routes)

		assert.Nil(res)

		if assert.NotNil(err) {
			assert.Contains(err.Message, "no have route")
		}
	})

	t.Run("no have routes", func(t *testing.T) {
		res, err := client.BestRoute("ORL", "SCL", []*models.Route{})

		assert.Nil(res)

		if assert.NotNil(err) {
			assert.Contains(err.Message, "no have routes")
		}
	})

	t.Run("success same place", func(t *testing.T) {
		res, err := client.BestRoute("GRU", "GRU", routes)

		assert.Nil(err)

		if assert.NotNil(res) {
			assert.Equal(1, len(res.Route))
			assert.Equal(0, res.Cost)
		}
	})

	t.Run("success less", func(t *testing.T) {
		route := routes[1]

		res, err := client.BestRoute(route.Start, route.Target, routes)

		assert.Nil(err)

		if assert.NotNil(res) {
			assert.Equal(2, len(res.Route))
			assert.Equal(route.Cost, res.Cost)
		}
	})

	t.Run("success", func(t *testing.T) {
		res, err := client.BestRoute("GRU", "CDG", routes)

		assert.Nil(err)

		if assert.NotNil(res) {
			assert.Equal(5, len(res.Route))
			assert.Equal(40, res.Cost)
		}
	})
}

func TestMain(m *testing.M) {
	os.Exit(func() int {
		routes = []*models.Route{
			{
				Start:  "GRU",
				Target: "BRC",
				Cost:   10,
			},
			{
				Start:  "BRC",
				Target: "SCL",
				Cost:   5,
			},
			{
				Start:  "GRU",
				Target: "CDG",
				Cost:   75,
			},
			{
				Start:  "GRU",
				Target: "SCL",
				Cost:   20,
			},
			{
				Start:  "GRU",
				Target: "ORL",
				Cost:   56,
			},
			{
				Start:  "ORL",
				Target: "CDG",
				Cost:   5,
			},
			{
				Start:  "SCL",
				Target: "ORL",
				Cost:   20,
			},
			{
				Start:  "SCL",
				Target: "ORL",
				Cost:   30,
			},
		}

		client = NewDjkClient()
		return m.Run()
	}())
}
