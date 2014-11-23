package types

import (
	"time"
)

//
// Json Types
//

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
	Id      string
	Url     string
	Author  string
	Content string
	Created time.Time
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
	Interests []string
	Languages []string
	Location  string
}
