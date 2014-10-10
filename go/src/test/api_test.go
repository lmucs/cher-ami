package api_test

import (
	. "gopkg.in/check.v1"
	api "../api"
	routes "../routes"
	"encoding/json"
	"fmt"
	"github.com/jmcvetta/neoism"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	server *httptest.Server
	a *api.Api

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

/* 
	Hook up gocheck into the "go test" runner.
*/
func Test(t *testing.T) {
	TestingT(t)
}

/*
	Suite-based grouping of tests.
*/
type TestSuite struct {
}

/*
	Suite registers the given value as a test suite to be run. 
	Any methods starting with the Test prefix in the given value will be considered as a test method.
*/
var _ = Suite(&TestSuite{})

/*
	Run once when the suite starts running.
*/
func (s *TestSuite) SetUpSuite(c *C) {
	uri := "http://192.241.226.228:7474/db/data"
	neo4jdb, err := neoism.Connect(uri)
	if err != nil {
		log.Fatal(err)
	}

	a = &api.Api{neo4jdb}
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

/*
	Run before each test or benchmark starts running.
*/
func (s *TestSuite) SetUpTest(c *C) {
	a.DatabaseInit()
}

/*
	Run after each test or benchmark runs.
*/
func (s *TestSuite) TearDownTest(c *C) {
	a.Db.Cypher(&neoism.CypherQuery {
		Statement: `
            MATCH (n)
            OPTIONAL MATCH (n)-[r]-()
            DELETE n, r
        `,
	})
}

/*
	Run once after all tests or benchmarks have finished running.
*/
func (s *TestSuite) TearDownSuite(c *C) {
	server.Close()
}

func postSignup(handle string, email string, password string, confirmPassword string) (*http.Response, error) {
	proposal := "{\"Handle\": \"" + handle + "\", \"Email\": \"" + email + "\", \"Password\": \"" + password + "\", \"ConfirmPassword\": \"" + confirmPassword + "\"}"

	reader = strings.NewReader(proposal)

	request, err := http.NewRequest("POST", signupURL, reader)
	if err != nil {
		log.Fatal(err)
	}

	response, err := http.DefaultClient.Do(request)

	return response, err
}

func postLogin(handle string, password string) (*http.Response, error) {
	credentials := "{\"Handle\": \"" + handle + "\", \"Password\": \"" + password + "\"}"

	reader = strings.NewReader(credentials)

	request, err := http.NewRequest("POST", loginURL, reader)
	if err != nil {
		log.Fatal(err)
	}

	response, err := http.DefaultClient.Do(request)

	return response, err
}

func postLogout(handle string) (*http.Response, error) {
	user := "{\"Handle\": \"" + handle + "\"}"

	reader = strings.NewReader(user)

	request, err := http.NewRequest("POST", logoutURL, reader)
	if err != nil {
		log.Fatal(err)
	}

	response, err := http.DefaultClient.Do(request)

	return response, err
}

func getJsonResponseMessage(response *http.Response) (string) {
	type Json struct {
		Response string
	}

	var message Json

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}  

	err = json.Unmarshal(body, &message)
	if err != nil {
		log.Fatal(err)
	}

	return message.Response
}

func getJsonErrorMessage(response *http.Response) (string) {
	type Json struct {
		Error string
	}

	var message Json

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &message)
	if err != nil {
		log.Fatal(err)
	}
 
	return message.Error
}

func (s *TestSuite) TestSignupEmptyHandle(c *C) {
	response, err := postSignup("", "testing123", "testing123", "testing123")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Handle is a required field for signup")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestSignupEmptyEmail(c *C) {
	response, err := postSignup("testing123", "", "testing123", "testing123")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Email is a required field for signup")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestSignupPasswordMismatch(c *C) {
	response, err := postSignup("testing123", "testing123", "testing123", "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Passwords do not match")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestSignupPasswordTooShort(c *C) {
	entry := "testing"

	for i := len(entry); i >= 0; i-- {
		response, err := postSignup("testing123", "testing123", entry[:len(entry)-i], entry[:len(entry)-i])
		if err != nil {
			c.Error(err)
		}

		c.Check(getJsonResponseMessage(response), Equals, "Passwords must be at least 8 characters long")
		c.Assert(response.StatusCode, Equals, 400, Commentf("Password length = %d.", len(entry)-i))
	}
}

func (s *TestSuite) TestSignupHandleTaken(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	
	response, err := postSignup("testing123", "testing321", "testing123", "testing123")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Sorry, testing123 is already taken")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestSignupEmailTaken(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	
	response, err := postSignup("testing321", "testing123", "testing123", "testing123")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Sorry, testing123 is already taken")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestSignupCreated(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	
	response, err := postSignup("testing321", "testing321", "testing123", "testing123")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Signed up a new user!")
	c.Assert(response.StatusCode, Equals, 201)
}

func (s *TestSuite) TestLoginInvalidUsername(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	response, err := postLogin("testing321", "testing123")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonErrorMessage(response), Equals, "Invalid username or password, please try again.")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestLoginInvalidPassword(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	response, err := postLogin("testing123", "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonErrorMessage(response), Equals, "Invalid username or password, please try again.")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestLoginOK(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	response, err := postLogin("testing123", "testing123")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Logged in testing123. Note your session id.")
	c.Assert(response.StatusCode, Equals, 200)
}

func (s *TestSuite) TestLogoutUserNoExist(c *C) {
	response, err := postLogout("testing123")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "No user was logged out")
	c.Assert(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestLogoutUserNoLogin(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	postLogout("testing123")

	response, err := postLogout("testing123")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "No user was logged out")
	c.Assert(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestLogoutOK(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	postLogin("testing123", "testing123")

	response, err := postLogout("testing123")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Logged out testing123, have a nice day")
	c.Assert(response.StatusCode, Equals, 200)
}
