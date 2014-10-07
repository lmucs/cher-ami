package api_test

import (
	api "../api"
	routes "../routes"
	"fmt"
	"github.com/jmcvetta/neoism"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	server *httptest.Server

	signupURL   string
	loginURL    string
	logoutURL   string
	usersURL    string
	messagesURL string
	publishURL  string
	followURL   string
	blockURL    string
	circlesURL  string

	reader io.Reader
)

func setup() {
	uri := "http://192.241.226.228:7474/db/data"
	neo4jdb, err := neoism.Connect(uri)
	if err != nil {
		log.Fatal(err)
	}

	a := &api.Api{neo4jdb}
	handler, err := routes.MakeHandler(*a)
	if err != nil {
		log.Fatal(err)
	}

	server = httptest.NewServer(&handler)

	signupURL = fmt.Sprintf("%s/signup", server.URL)
	loginURL = fmt.Sprintf("%s/login", server.URL)
	logoutURL = fmt.Sprintf("%s/logout", server.URL)
	usersURL = fmt.Sprintf("%s/users", server.URL)
	messagesURL = fmt.Sprintf("%s/messages", server.URL)
	publishURL = fmt.Sprintf("%s/publish", server.URL)
	followURL = fmt.Sprintf("%s/follow", server.URL)
	blockURL = fmt.Sprintf("%s/block", server.URL)
	circlesURL = fmt.Sprintf("%s/circles", server.URL)
}

func teardown() {
	server.Close()
}

func postSignup(handle string, email string, password string, confirmPassword string) (*http.Response, error) {
	proposal := "{\"Handle\": \"" + handle + "\", \"Email\": \"" + email + "\", \"Password\": \"" + password + "\", \"ConfirmPassword\": \"" + confirmPassword + "\"}"

	reader = strings.NewReader(proposal)

	request, err := http.NewRequest("POST", signupURL, reader)

	response, err := http.DefaultClient.Do(request)

	return response, err
}

func postLogin(handle string, password string) (*http.Response, error) {
	credentials := "{\"Handle\": \"" + handle + "\", \"Password\": \"" + password + "\"}"

	reader = strings.NewReader(credentials)

	request, err := http.NewRequest("POST", loginURL, reader)

	response, err := http.DefaultClient.Do(request)

	return response, err
}

func TestSignupEmptyHandle(t *testing.T) {
	setup()
	defer teardown()

	response, err := postSignup("", "testing123", "testing123", "testing123")

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 400 {
		t.Errorf("HTTP Status Code: %d", response.StatusCode)
	}
}

func TestSignupEmptyEmail(t *testing.T) {
	setup()
	defer teardown()

	response, err := postSignup("testing123", "", "testing123", "testing123")

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 400 {
		t.Errorf("HTTP Status Code: %d", response.StatusCode)
	}
}

func TestSignupPasswordMismatch(t *testing.T) {
	setup()
	defer teardown()

	response, err := postSignup("testing123", "testing123", "testing123", "testing321")

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 400 {
		t.Errorf("HTTP Status Code: %d", response.StatusCode)
	}
}

func TestSignupPasswordTooShort(t *testing.T) {
	setup()
	defer teardown()

	entry := "testing"

	for i := len(entry); i >= 0; i-- {

		response, err := postSignup("testing123", "testing123", entry[:len(entry)-i], entry[:len(entry)-i])

		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != 400 {
			t.Errorf("HTTP Status Code: %d, Password Length: %d", response.StatusCode, len(entry)-i)
		}
	}
}

func TestSignupHandleTaken(t *testing.T) {
	setup()
	defer teardown()

	postSignup("testing123", "testing123", "testing123", "testing123")
	response, err := postSignup("testing123", "testing321", "testing123", "testing123")

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 400 {
		t.Errorf("HTTP Status Code: %d", response.StatusCode)
	}
}

func TestSignupEmailTaken(t *testing.T) {
	setup()
	defer teardown()

	postSignup("testing123", "testing123", "testing123", "testing123")
	response, err := postSignup("testing321", "testing123", "testing123", "testing123")

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 400 {
		t.Errorf("HTTP Status Code: %d", response.StatusCode)
	}
}

func TestSignupCreated(t *testing.T) {
	setup()
	defer teardown()

	postSignup("testing123", "testing123", "testing123", "testing123")
	response, err := postSignup("testing321", "testing321", "testing123", "testing123")

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 201 {
		t.Errorf("HTTP Status Code: %d", response.StatusCode)
	}
}

func TestLoginInvalidUsername(t *testing.T) {
	setup()
	defer teardown()

	postSignup("testing123", "testing123", "testing123", "testing123")

	response, err := postLogin("testing321", "testing123")

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 400 {
		t.Errorf("HTTP Status Code: %d", response.StatusCode)
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	setup()
	defer teardown()

	postSignup("testing123", "testing123", "testing123", "testing123")

	response, err := postLogin("testing123", "testing321")

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 400 {
		t.Errorf("HTTP Status Code: %d", response.StatusCode)
	}
}

func TestLoginOK(t *testing.T) {
	setup()
	defer teardown()

	postSignup("testing123", "testing123", "testing123", "testing123")

	response, err := postLogin("testing123", "testing123")

	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 200 {
		t.Errorf("HTTP Status Code: %d", response.StatusCode)
	}
}
