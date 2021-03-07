package main

import (
	"best-route/models"
	"fmt"
	"strings"
	"testing"

	assert2 "github.com/stretchr/testify/assert"
)

func TestCLI_validateCLIBestRouteRequest(t *testing.T) {
	assert := assert2.New(t)

	t.Run("not contains quotation marks", func(t *testing.T) {
		route, resErr := models.ValidateCLIBestRouteRequest("GRUCDG")

		assert.Nil(route)

		if assert.NotNil(resErr) {
			assert.Contains(resErr.ToString(), "places not separated by")
		}
	})

	t.Run("max size name", func(t *testing.T) {
		route, resErr := models.ValidateCLIBestRouteRequest(fmt.Sprintf("%v-GRU", strings.Repeat("A", models.MaxSizePlaceName+1)))

		assert.Nil(route)

		if assert.NotNil(resErr) {
			assert.Contains(resErr.ToString(), "place name is too big")
		}
	})

	t.Run("empty place name", func(t *testing.T) {
		route, resErr := models.ValidateCLIBestRouteRequest("GRU-")

		assert.Nil(route)

		if assert.NotNil(resErr) {
			assert.Contains(resErr.ToString(), "place name is empty")
		}
	})

	t.Run("success", func(t *testing.T) {
		route, resErr := models.ValidateCLIBestRouteRequest("GRU-CDG")

		assert.Nil(resErr)

		if assert.NotNil(route) {
			assert.Equal("GRU", route.Start)
			assert.Equal("CDG", route.Target)
		}
	})
}
