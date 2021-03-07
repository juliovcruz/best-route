package models

import (
	"fmt"
	"strings"
)

func ValidateRequest(route *Route) *ResponseError {
	if route == nil {
		return &ResponseError{
			Message: "route is empty",
		}
	}

	if len(route.Start) > MaxSizePlaceName {
		return &ResponseError{
			Message: fmt.Sprintf("start place name is too big, max is %v characteres", MaxSizePlaceName),
			Field:   "start",
		}
	}

	if len(route.Target) > MaxSizePlaceName {
		return &ResponseError{
			Message: fmt.Sprintf("target place name is too big, max is %v characteres", MaxSizePlaceName),
			Field:   "target",
		}
	}

	if len(route.Start) < 1 {
		return &ResponseError{
			Message: "start place name is empty",
			Field:   "start",
		}
	}

	if len(route.Target) < 1 {
		return &ResponseError{
			Message: "target place name is empty",
			Field:   "target",
		}
	}

	return nil
}

func ValidateInsertRequest(route *Route) *ResponseError {
	if err := ValidateRequest(route); err != nil {
		return err
	}

	if route.Cost < 0 {
		return &ResponseError{
			Message: "cost cannot be negative",
			Field:   "cost",
		}
	}

	return nil
}

func ValidateCLIBestRouteRequest(str string) (*Route, *ResponseError) {
	if !strings.Contains(str, "-") {
		return nil, &ResponseError{
			Message: "places not separated by \"-\"",
		}
	}

	split := strings.Split(str, "-")

	route := &Route{
		Start:  split[0],
		Target: split[1],
	}

	if err := ValidateRequest(route); err != nil {
		return nil, err
	}

	return route, nil
}
