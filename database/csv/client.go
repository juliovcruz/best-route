package csv

import (
	"best-route/models"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
)

type CsvClient struct {
	Path  string
	Routes []*models.Route
}

func NewCsvClient(path string) (*CsvClient, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)

	routes, err := CsvReaderToRoutes(reader)
	if err != nil {
		return nil, err
	}

	return &CsvClient{
		Path:  path,
		Routes: routes,
	}, nil
}

func (c *CsvClient) InsertOne(route *models.Route) (*models.Route, error) {
	var lines [][]string

	for _, r := range c.Routes {
		// check if already exists
		if r.Equal(route) {
			return nil, errors.New("already exists")
		}

		lines = append(lines, r.ToString())
	}
	lines = append(lines, route.ToString())

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

func CsvReaderToRoutes(reader *csv.Reader) ([]*models.Route, error) {
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
			return nil, err
		}

		routes = append(routes, &models.Route{
			Start:  record[0],
			Target: record[1],
			Cost:   cost,
		})
	}

	return routes, nil
}

func NewMockCsvClient(opts *CsvClient) (*CsvClient, error) {
	var lines [][]string

	file, err := os.Create(opts.Path)
	if err != nil {
		return nil, err
	}

	for _, r := range opts.Routes {
		lines = append(lines, r.ToString())
	}

	w := csv.NewWriter(file)
	if err = w.WriteAll(lines); err != nil {
		file.Close()
		return nil, err
	}

	return &CsvClient{
		Path:   opts.Path,
		Routes: opts.Routes,
	}, nil
}