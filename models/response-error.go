package models

import (
	"encoding/json"
)

type ResponseError struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewResponseErrorToJSON(field string, msg string) string {
	res := &ResponseError{
		Field:   field,
		Message: msg,
	}

	js, _ := json.Marshal(res)

	return string(js)
}

func (e *ResponseError) ToJSON() []byte {
	js, _ := json.Marshal(e)
	return js
}

func (e *ResponseError) ToString() string {
	return e.Message
}
