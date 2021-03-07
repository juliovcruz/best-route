package csv

import (
	"best-route/database"
	"best-route/models"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Client struct {
	Path   string
	Routes []*models.Route
}

func NewClient(path string) (*Client, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)

	routes, err := csvReaderToRoutes(reader)
	if err != nil {
		return nil, err
	}

	return &Client{
		Path:   path,
		Routes: routes,
	}, nil
}

func csvReaderToRoutes(reader *csv.Reader) ([]*models.Route, error) {
	var routes []*models.Route

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		cost, err := strconv.Atoi(record[2])
		if err != nil {
			fmt.Printf("route not used because cost isn't a valid number, %v > %v $ %v\n error: %v\n", record[0], record[1], record[2], err.Error())
			continue
		}

		route := &models.Route{
			Start:  record[0],
			Target: record[1],
			Cost:   cost,
		}

		if resErr := models.ValidateInsertRequest(route); resErr != nil {
			fmt.Printf("route not used because %v\n", resErr.Message)
			continue
		}

		routes = append(routes, route)
	}

	return routes, nil
}

func (c *Client) InsertOneRoute(route *models.Route) (*models.Route, error) {
	var lines [][]string

	for _, r := range c.Routes {
		// check if already exists
		if r.Equal(route) {
			return nil, errors.New(database.AlreadyExistsMessageError)
		}

		lines = append(lines, r.ToArray())
	}
	lines = append(lines, route.ToArray())

	file, err := os.Create(c.Path)
	if err != nil {
		return nil, err
	}

	w := csv.NewWriter(file)
	if err = w.WriteAll(lines); err != nil {
		file.Close()
		return nil, err
	}

	c.Routes = append(c.Routes, route)

	return route, err
}

func (c *Client) GetAllRoutes() ([]*models.Route, error) {
	return c.Routes, nil
}

func NewMockClient(opts *Client) (*Client, error) {
	var lines [][]string

	file, err := os.Create(opts.Path)
	if err != nil {
		return nil, err
	}

	for _, r := range opts.Routes {
		lines = append(lines, r.ToArray())
	}

	w := csv.NewWriter(file)
	if err = w.WriteAll(lines); err != nil {
		file.Close()
		return nil, err
	}

	return &Client{
		Path:   opts.Path,
		Routes: opts.Routes,
	}, nil
}
