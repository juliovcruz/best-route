package main

import (
	"best-route/database/csv"
	"best-route/models"
	"best-route/route_calculator"
	"best-route/route_calculator/djk"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	assert2 "github.com/stretchr/testify/assert"
)

const PathTest = "test-input-not-edit.csv"

var (
	s          *server
	seedRoutes map[string]*models.Route
)

func TestAPI_Insert(t *testing.T) {
	assert := assert2.New(t)

	t.Run("body nil", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/add", nil)
		if err != nil {
			t.Fatal(err)
		}

		res := httptest.NewRecorder()

		handler := http.HandlerFunc(s.handleInsert)
		handler.ServeHTTP(res, req)

		if assert.NotNil(res) {
			var resErr *models.ResponseError

			err := json.Unmarshal(res.Body.Bytes(), &resErr)
			assert.NoError(err)
			assert.Contains(resErr.Message, "body cannot be empty")

			assert.Equal(http.StatusBadRequest, res.Code)
		}
	})

	t.Run("wrong method", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/add", nil)
		if err != nil {
			t.Fatal(err)
		}

		res := httptest.NewRecorder()

		handler := http.HandlerFunc(s.handleInsert)
		handler.ServeHTTP(res, req)

		if assert.NotNil(res) {
			var resErr *models.ResponseError

			err := json.Unmarshal(res.Body.Bytes(), &resErr)
			assert.NoError(err)
			assert.Contains(resErr.Message, "wrong method")

			assert.Equal(http.StatusMethodNotAllowed, res.Code)
		}
	})

	t.Run("big place name", func(t *testing.T) {
		route := &models.Route{
			Start:  strings.Repeat("A", models.MaxSizePlaceName+1),
			Target: "TO",
			Cost:   10,
		}

		req, err := http.NewRequest(http.MethodPost, "/add", bytes.NewBuffer(route.ToJSON()))
		if err != nil {
			t.Fatal(err)
		}

		res := httptest.NewRecorder()

		handler := http.HandlerFunc(s.handleInsert)
		handler.ServeHTTP(res, req)

		if assert.NotNil(res) {
			var resErr *models.ResponseError

			err := json.Unmarshal(res.Body.Bytes(), &resErr)
			assert.NoError(err)
			assert.Contains(resErr.Message, "place name is too big")

			assert.Equal(http.StatusPreconditionFailed, res.Code)
		}
	})

	t.Run("empty place name", func(t *testing.T) {
		route := &models.Route{
			Start:  "",
			Target: "",
			Cost:   10,
		}

		req, err := http.NewRequest(http.MethodPost, "/add", bytes.NewBuffer(route.ToJSON()))
		if err != nil {
			t.Fatal(err)
		}

		res := httptest.NewRecorder()

		handler := http.HandlerFunc(s.handleInsert)
		handler.ServeHTTP(res, req)

		if assert.NotNil(res) {
			var resErr *models.ResponseError

			err := json.Unmarshal(res.Body.Bytes(), &resErr)
			assert.NoError(err)
			assert.Contains(resErr.Message, "place name is empty")

			assert.Equal(http.StatusPreconditionFailed, res.Code)
		}
	})

	t.Run("negative cost", func(t *testing.T) {
		route := &models.Route{
			Start:  "TO",
			Target: "SA",
			Cost:   -5,
		}

		req, err := http.NewRequest(http.MethodPost, "/add", bytes.NewBuffer(route.ToJSON()))
		if err != nil {
			t.Fatal(err)
		}

		res := httptest.NewRecorder()

		handler := http.HandlerFunc(s.handleInsert)
		handler.ServeHTTP(res, req)

		if assert.NotNil(res) {
			var resErr *models.ResponseError

			err := json.Unmarshal(res.Body.Bytes(), &resErr)
			assert.NoError(err)
			assert.Contains(resErr.Message, "cost cannot be negative")

			assert.Equal(http.StatusPreconditionFailed, res.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		route := &models.Route{
			Start:  "SP ",
			Target: "TO",
			Cost:   10,
		}

		req, err := http.NewRequest(http.MethodPost, "/add", bytes.NewBuffer(route.ToJSON()))
		if err != nil {
			t.Fatal(err)
		}

		res := httptest.NewRecorder()

		handler := http.HandlerFunc(s.handleInsert)
		handler.ServeHTTP(res, req)

		if assert.NotNil(res) {
			var routeRes *models.Route

			err := json.Unmarshal(res.Body.Bytes(), &routeRes)
			assert.NoError(err)
			assert.True(route.Trim().Equal(routeRes))

			assert.Equal(http.StatusOK, res.Code)
		}
	})
}

func TestAPI_BestRoute(t *testing.T) {
	assert := assert2.New(t)

	t.Run("body nil", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/best", nil)
		if err != nil {
			t.Fatal(err)
		}

		res := httptest.NewRecorder()

		handler := http.HandlerFunc(s.handleBestRoute)
		handler.ServeHTTP(res, req)

		if assert.NotNil(res) {
			var resErr *models.ResponseError

			err := json.Unmarshal(res.Body.Bytes(), &resErr)
			assert.NoError(err)
			assert.Contains(resErr.Message, "body cannot be empty")

			assert.Equal(http.StatusBadRequest, res.Code)
		}
	})

	t.Run("wrong method", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/best", nil)
		if err != nil {
			t.Fatal(err)
		}

		res := httptest.NewRecorder()

		handler := http.HandlerFunc(s.handleBestRoute)
		handler.ServeHTTP(res, req)

		if assert.NotNil(res) {
			var resErr *models.ResponseError

			err := json.Unmarshal(res.Body.Bytes(), &resErr)
			assert.NoError(err)
			assert.Contains(resErr.Message, "wrong method")

			assert.Equal(http.StatusMethodNotAllowed, res.Code)
		}
	})

	t.Run("big place name", func(t *testing.T) {
		route := &models.Route{
			Start:  strings.Repeat("A", models.MaxSizePlaceName+1),
			Target: "TO",
		}

		req, err := http.NewRequest(http.MethodGet, "/best", bytes.NewBuffer(route.ToJSON()))
		if err != nil {
			t.Fatal(err)
		}

		res := httptest.NewRecorder()

		handler := http.HandlerFunc(s.handleBestRoute)
		handler.ServeHTTP(res, req)

		if assert.NotNil(res) {
			var resErr *models.ResponseError

			err := json.Unmarshal(res.Body.Bytes(), &resErr)
			assert.NoError(err)
			assert.Contains(resErr.Message, "place name is too big")

			assert.Equal(http.StatusPreconditionFailed, res.Code)
		}
	})

	t.Run("empty place name", func(t *testing.T) {
		route := &models.Route{
			Start:  "",
			Target: "",
		}

		req, err := http.NewRequest(http.MethodGet, "/best", bytes.NewBuffer(route.ToJSON()))
		if err != nil {
			t.Fatal(err)
		}

		res := httptest.NewRecorder()

		handler := http.HandlerFunc(s.handleBestRoute)
		handler.ServeHTTP(res, req)

		if assert.NotNil(res) {
			var resErr *models.ResponseError

			err := json.Unmarshal(res.Body.Bytes(), &resErr)
			assert.NoError(err)
			assert.Contains(resErr.Message, "place name is empty")

			assert.Equal(http.StatusPreconditionFailed, res.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		route := &models.Route{
			Start:  "GRU",
			Target: "CDG",
		}

		req, err := http.NewRequest(http.MethodGet, "/best", bytes.NewBuffer(route.ToJSON()))
		if err != nil {
			t.Fatal(err)
		}

		res := httptest.NewRecorder()

		handler := http.HandlerFunc(s.handleBestRoute)
		handler.ServeHTTP(res, req)

		if assert.NotNil(res) {
			var routeRes *route_calculator.BestRouteResponse

			err := json.Unmarshal(res.Body.Bytes(), &routeRes)
			assert.NoError(err)
			assert.Equal(len(routeRes.Route), 5)
			assert.Equal(routeRes.Cost, 40)

			assert.Equal(http.StatusOK, res.Code)
		}
	})
}

func TestMain(m *testing.M) {
	os.Exit(func() int {
		routes := SeedRoutesToTest()

		csvClient, err := csv.NewMockClient(&csv.Client{
			Path:   "./" + PathTest,
			Routes: routes,
		})
		if err != nil {
			log.Fatal(err)
		}

		sv := &server{
			DatabaseClient: csvClient,
			RouterClient:   djk.NewClient(),
			Addr:           ":3001",
		}

		s = sv

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
