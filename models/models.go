package models

import (
	"fmt"
)

type Route struct {
	Start  string `json:"start,omitempty"`
	Target string `json:"target,omitempty"`
	Cost   int    `json:"cost,omitempty"`
}

func (r *Route) ToString() []string {
	return []string{r.Start, r.Target, fmt.Sprintf("%v", r.Cost)}
}

func (r *Route) Equal(route *Route) bool {
	if route.Start == r.Start &&
		route.Target == r.Target &&
		route.Cost == r.Cost {
		return true
	}

	return false
}

