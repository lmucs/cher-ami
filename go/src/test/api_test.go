package api_test

import (
	"../api"
	"../routes"
	"./requester"
	"flag"
	"github.com/jadengore/goconfig"
	. "gopkg.in/check.v1"
	"io"
	"log"
	"net/http/httptest"
	"testing"
)

// Flag for local testing.
var local = flag.Bool("local", false, "For local testing")

var (
	server *httptest.Server
	a      *api.Api
	req    *requester.Requester
	reader io.Reader
)

//
// Hook up gocheck into the "go test" runner.
//

func Test(t *testing.T) {
	TestingT(t)
}

//
// Suite-based grouping of tests.
//

type TestSuite struct {
}

//
// Suite registers the given value as a test suite to be run.
// Any methods starting with the Test prefix in the given value will be considered as a test method.
//

var _ = Suite(&TestSuite{})

//
// Run once when the suite starts running.
//

func (s *TestSuite) SetUpSuite(c *C) {
	config, err := goconfig.ReadConfigFile("../../../config.cfg")
	var location string
	if *local {
		location = "local-test"
	} else {
		location = "api-test"
	}
	uri, err := config.GetString(location, "url")

	a = api.NewApi(uri)

	handler, err := routes.MakeHandler(*a, false)
	if err != nil {
		log.Fatal(err)
	}

	server = httptest.NewServer(&handler)

	// routes.signupURL = fmt.Sprintf("%s/signup", server.URL)
	// routes.changePassURL = fmt.Sprintf("%s/changepassword", server.URL)
	// routes.sessionsURL = fmt.Sprintf("%s/sessions", server.URL)
	// routes.userURL = fmt.Sprintf("%s/users/user", server.URL)
	// routes.usersURL = fmt.Sprintf("%s/users", server.URL)
	// routes.messagesURL = fmt.Sprintf("%s/messages", server.URL)
	// routes.publishURL = fmt.Sprintf("%s/publish", server.URL)
	// routes.joindefaultURL = fmt.Sprintf("%s/joindefault", server.URL)
	// routes.joinURL = fmt.Sprintf("%s/join", server.URL)
	// routes.blockURL = fmt.Sprintf("%s/block", server.URL)
	// routes.circlesURL = fmt.Sprintf("%s/circles", server.URL)

	req = requester.NewRequester(server.URL)

}

//
// Run before each test or benchmark starts running.
//

func (s *TestSuite) SetUpTest(c *C) {
}

//
// Run after each test or benchmark runs.
//

func (s *TestSuite) TearDownTest(c *C) {
	a.Svc.FreshInitialState()
}

//
// Run once after all tests or benchmarks have finished running.
//

func (s *TestSuite) TearDownSuite(c *C) {
	server.Close()
}
