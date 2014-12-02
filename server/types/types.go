package types

import (
	"time"
)

type Json map[string]interface{}

type JsonArray []Json

// The json annotations that accompany these structs allow json.Marshall
// to to produce proper json instead of an escaped json string.

type SearchCirclesResponse struct {
	Results  []CircleResponse `json:"results"`
	Response string           `json:"response"`
	Count    int              `json:"count"`
}

type CircleResponse struct {
	Name        string    `json:"name"`
	Url         string    `json:"url"`
	Description string    `json:"description"`
	Owner       string    `json:"owner"`
	Visibility  string    `json:"visibility"`
	Members     string    `json:"members"`
	Created     time.Time `json:"created"`
}

type MessageView struct {
	Id      string    `json:"id"`
	Url     string    `json:"url"`
	Author  string    `json:"author"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
}

type MessageResponseView struct {
	Objects []MessageView `json:"objects"`
	Count   int           `json:"count"`
}
