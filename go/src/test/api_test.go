package api_test

import (
	. "gopkg.in/check.v1"
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

func (s *TestSuite) TestSignupEmptyHandle(c *C) {
	response, err := postSignup("", "testing123", "testing123", "testing123")
	if err != nil {
		c.Error(err)
	}

	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestSignupEmptyEmail(c *C) {
	response, err := postSignup("testing123", "", "testing123", "testing123")
	if err != nil {
		c.Error(err)
	}

	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestSignupPasswordMismatch(c *C) {
	response, err := postSignup("testing123", "testing123", "testing123", "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestSignupPasswordTooShort(c *C) {
	entry := "testing"

	for i := len(entry); i >= 0; i-- {
		response, err := postSignup("testing123", "testing123", entry[:len(entry)-i], entry[:len(entry)-i])
		if err != nil {
			c.Error(err)
		}

		c.Assert(response.StatusCode, Equals, 400, Commentf("Password length = %d.", len(entry)-i))
	}
}

func (s *TestSuite) TestSignupHandleTaken(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	
	response, err := postSignup("testing123", "testing321", "testing123", "testing123")
	if err != nil {
		c.Error(err)
	}

	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestSignupEmailTaken(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	
	response, err := postSignup("testing321", "testing123", "testing123", "testing123")
	if err != nil {
		c.Error(err)
	}

	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestSignupCreated(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	
	response, err := postSignup("testing321", "testing321", "testing123", "testing123")
	if err != nil {
		c.Error(err)
	}

	c.Assert(response.StatusCode, Equals, 201)
}

func (s *TestSuite) TestLoginInvalidUsername(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	response, err := postLogin("testing321", "testing123")
	if err != nil {
		c.Error(err)
	}

	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestLoginInvalidPassword(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	response, err := postLogin("testing123", "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestLoginOK(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	response, err := postLogin("testing123", "testing123")
	if err != nil {
		c.Error(err)
	}

	c.Assert(response.StatusCode, Equals, 200)
}
