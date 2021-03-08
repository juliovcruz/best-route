package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

const MaxSizePlaceName = 30

type Route struct {
	Start  string `json:"start"`
	Target string `json:"target"`
	Cost   int    `json:"cost"`
}

func (r *Route) ToJSON() []byte {
	js, _ := json.Marshal(r)
	return js
}

func (r *Route) ToArray() []string {
	return []string{strings.ToUpper(strings.TrimSpace(r.Start)), strings.ToUpper(strings.TrimSpace(r.Target)), fmt.Sprintf("%v", r.Cost)}
}

func (r *Route) Trim() *Route {
	return &Route{
		Start:  strings.ToUpper(strings.TrimSpace(r.Start)),
		Target: strings.ToUpper(strings.TrimSpace(r.Target)),
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
