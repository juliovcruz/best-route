package main

import (
	"best-route/models"
	"fmt"
	"strings"
)

func validateRequest(route *models.Route) *models.ResponseError {
	if route == nil {
		return &models.ResponseError{
			Message: "route is empty",
		}
	}

	if len(route.Start) > MaxSizePlaceName {
		return &models.ResponseError{
			Message: fmt.Sprintf("start place name is too big, max is %v characteres", MaxSizePlaceName),
			Field:   "start",
		}
	}

	if len(route.Target) > MaxSizePlaceName {
		return &models.ResponseError{
			Message: fmt.Sprintf("target place name is too big, max is %v characteres", MaxSizePlaceName),
			Field:   "target",
		}
	}

	if len(route.Start) < 1 {
		return &models.ResponseError{
			Message: "start place name is empty",
			Field:   "start",
		}
	}

	if len(route.Target) < 1 {
		return &models.ResponseError{
			Message: "target place name is empty",
			Field:   "target",
		}
	}

	return nil
}

func validateAPIInsertRequest(route *models.Route) *models.ResponseError {
	if err := validateRequest(route); err != nil {
		return err
	}

	if route.Cost < 0 {
		return &models.ResponseError{
			Message: "cost cannot be negative",
			Field:   "cost",
		}
	}

	return nil
}

func validateAPIBestRouteRequest(route *models.Route) *models.ResponseError {
	if err := validateRequest(route); err != nil {
		return err
	}

	return nil
}

func validateCLIBestRouteRequest(str string) (*models.Route, *models.ResponseError) {
	if !strings.Contains(str, "-") {
		return nil, &models.ResponseError{
			Message: "places not separated by \"-\"",
		}
	}

	split := strings.Split(str, "-")

	route := &models.Route{
		Start:  split[0],
		Target: split[1],
	}

	if err := validateRequest(route); err != nil {
		return nil, err
	}

	return route, nil
}
