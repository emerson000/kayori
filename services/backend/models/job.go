package models

import "encoding/json"

type Job struct {
	ID      string          `json:"id"`
	Title   string          `json:"title"`
	Service string          `json:"service"`
	Task    json.RawMessage `json:"task"`
}
