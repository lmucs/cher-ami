package types

import (
	"time"
)

type Json map[string]interface{}

type JsonArray []Json

//
// Account Types
//

type SignupProposal struct {
	Handle          string `json:"handle" validate:"handle"`
	Email           string `json:"email" validate:"email"`
	Password        string `json:"password" validate:"password"`
	ConfirmPassword string `json:"confirmpassword" validate:"password"`
}

type LoginCredentials struct {
	Handle   string `json:"handle" validate:"handle"`
	Password string `json:"password" validate:"password"`
}

//
// Message Types
//

type Message struct {
	Id      string    `json:"id"`
	Url     string    `json:"url"`
	Author  string    `json:"author"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
}

type NewMessage struct {
	Content string
	Circles []string
}

type MessagePatch struct {
	Op       string `json:"op" validate:"messageop"`
	Resource string `json:"resource" validate:"messageresource"`
	Value    string `json:"value" validate:"messagevalue"`
}

//
// User Types
//

type UserAttributes struct {
	FirstName string
	LastName  string
	Gender    string
	Birthday  time.Time
	Bio       string
	Interests string
	Languages string
	Location  string
}

type UserPatch struct {
	Resource string `json:"resource" validate:"userresource"`
	Value    string `json:"value" validate:"uservalue"`
}

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
