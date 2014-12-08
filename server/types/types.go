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

type PublishedMessageView struct {
	Id          string    `json:"id"`
	Url         string    `json:"url"`
	Author      string    `json:"author"`
	Content     string    `json:"content"`
	Created     time.Time `json:"created"`
	Circleid    string    `json:"circleid"`
	PublishedAt time.Time `json:"published_at"`
}

type MessageView struct {
	Id      string    `json:"id"`
	Url     string    `json:"url"`
	Author  string    `json:"author"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
}

type PublishedMessageResponseView struct {
	Objects []PublishedMessageView `json:"objects"`
}

type MessageResponseView struct {
	Objects []MessageView `json:"objects"`
	Count   int           `json:"count"`
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

type UserView struct {
	Handle    string           `json:"handle"`
	FirstName string           `json:"firstname"`
	LastName  string           `json:"lastname"`
	Gender    string           `json:"gender"`
	Birthday  time.Time        `json:"birthday"`
	Bio       string           `json:"bio"`
	Interests string           `json:"interests"`
	Languages string           `json:"languages"`
	Location  string           `json:"location"`
	Circles   []CircleResponse `json:"circles"`
	Blocked   []UserView       `json:"blocked"`
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

//
// Test types
//

type ReasonCatcher struct {
	Reason string `json:"reason"`
}
