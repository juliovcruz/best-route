package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Route struct {
	Start  string `json:"start,omitempty"`
	Target string `json:"target,omitempty"`
	Cost   int    `json:"cost,omitempty"`
}

func (r *Route) ToJSON() []byte {
	js, _ := json.Marshal(r)
	return js
}

func (r *Route) ToString() []string {
	return []string{strings.TrimSpace(r.Start), strings.TrimSpace(r.Target), fmt.Sprintf("%v", r.Cost)}
}

func (r *Route) Trim() *Route {
	return &Route{
		Start:  strings.TrimSpace(r.Start),
		Target: strings.TrimSpace(r.Target),
		Cost:   r.Cost,
	}
}

func (r *Route) Equal(route *Route) bool {
	if route.Start == r.Start &&
		route.Target == r.Target &&
		route.Cost == r.Cost {
		return true
	}

	return false
}
