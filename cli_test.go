package main

import (
	"fmt"
	"strings"
	"testing"

	assert2 "github.com/stretchr/testify/assert"
)

func TestCLI_validateCLIBestRouteRequest(t *testing.T) {
	assert := assert2.New(t)

	t.Run("not contains quotation marks", func(t *testing.T) {
		route, resErr := validateCLIBestRouteRequest("GRUCDG")

		assert.Nil(route)

		if assert.NotNil(resErr) {
			assert.Contains(resErr.ToString(), "places not separated by")
		}
	})

	t.Run("max size name", func(t *testing.T) {
		route, resErr := validateCLIBestRouteRequest(fmt.Sprintf("%v-GRU", strings.Repeat("A", MaxSizePlaceName+1)))

		assert.Nil(route)

		if assert.NotNil(resErr) {
			assert.Contains(resErr.ToString(), "place name is too big")
		}
	})

	t.Run("empty place name", func(t *testing.T) {
		route, resErr := validateCLIBestRouteRequest("GRU-")

		assert.Nil(route)

		if assert.NotNil(resErr) {
			assert.Contains(resErr.ToString(), "place name is empty")
		}
	})

	t.Run("success", func(t *testing.T) {
		route, resErr := validateCLIBestRouteRequest("GRU-CDG")

		assert.Nil(resErr)

		if assert.NotNil(route) {
			assert.Equal("GRU", route.Start)
			assert.Equal("CDG", route.Target)
		}
	})
}
