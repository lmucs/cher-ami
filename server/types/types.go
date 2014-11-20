package types

import (
	"time"
)

type Json map[string]interface{}

type JsonArray []Json

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

type SignupProposal struct {
	Handle          string `validate:"handle"`
	Email           string
	Password        string
	ConfirmPassword string
}
