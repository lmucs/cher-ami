package api_test

import (
	api "../api"
	routes "../routes"
	"./helper"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jadengore/goconfig"
	. "gopkg.in/check.v1"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Flag for local testing.
var local = flag.Bool("local", false, "For local testing")

var (
	server *httptest.Server
	a      *api.Api

	signupURL      string
	sessionsURL    string
	userURL        string
	usersURL       string
	messagesURL    string
	publishURL     string
	joindefaultURL string
	joinURL        string
	blockURL       string
	circlesURL     string

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

	handler, err := routes.MakeHandler(*a)
	if err != nil {
		log.Fatal(err)
	}

	server = httptest.NewServer(&handler)

	signupURL = fmt.Sprintf("%s/signup", server.URL)
	sessionsURL = fmt.Sprintf("%s/sessions", server.URL)
	userURL = fmt.Sprintf("%s/users/user", server.URL)
	usersURL = fmt.Sprintf("%s/users", server.URL)
	messagesURL = fmt.Sprintf("%s/messages", server.URL)
	publishURL = fmt.Sprintf("%s/publish", server.URL)
	joindefaultURL = fmt.Sprintf("%s/joindefault", server.URL)
	joinURL = fmt.Sprintf("%s/join", server.URL)
	blockURL = fmt.Sprintf("%s/block", server.URL)
	circlesURL = fmt.Sprintf("%s/circles", server.URL)
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

//
// Send/Receive Calls to API:
//

func postSignup(handle string, email string, password string, confirmPassword string) (*http.Response, error) {
	proposal := map[string]interface{}{
		"handle":          handle,
		"email":           email,
		"password":        password,
		"confirmpassword": confirmPassword,
	}

	return helper.Execute("POST", signupURL, proposal)
}

func postSessions(handle string, password string) (*http.Response, error) {
	payload := map[string]interface{}{
		"handle":   handle,
		"password": password,
	}

	return helper.Execute("POST", sessionsURL, payload)
}

func deleteSessions(handle string) (*http.Response, error) {
	payload := map[string]interface{}{
		"handle": handle,
	}

	return helper.Execute("DELETE", sessionsURL, payload)
}

func getUser(handle string) (*http.Response, error) {
	payload := map[string]interface{}{
		"handle": handle,
	}

	return helper.Execute("GET", userURL, payload)
}

func getUsers() (*http.Response, error) {
	payload := map[string]interface{}{}

	return helper.Execute("GET", usersURL, payload)
}

func deleteUser(handle string, password string, sessionid string) (*http.Response, error) {
	payload := map[string]interface{}{
		"handle":    handle,
		"password":  password,
		"sessionid": sessionid,
	}

	return helper.Execute("DELETE", userURL, payload)
}

func postCircles(handle string, sessionid string, circleName string, public bool) (*http.Response, error) {
	payload := map[string]interface{}{
		"handle":     handle,
		"sessionid":  sessionid,
		"circlename": circleName,
		"public":     public,
	}

	return helper.Execute("POST", circlesURL, payload)
}

func postBlock(handle string, sessionid string, target string) (*http.Response, error) {
	payload := map[string]interface{}{
		"handle":    handle,
		"sessionid": sessionid,
		"target":    target,
	}

	return helper.Execute("POST", blockURL, payload)
}

func postJoinDefault(handle string, sessionid string, target string) (*http.Response, error) {
	payload := map[string]interface{}{
		"handle":    handle,
		"sessionid": sessionid,
		"target":    target,
	}

	return helper.Execute("POST", joindefaultURL, payload)
}

func postJoin(handle string, sessionid string, target string, circle string) (*http.Response, error) {
	payload := map[string]interface{}{
		"handle":    handle,
		"sessionid": sessionid,
		"target":    target,
		"circle":    circle,
	}

	return helper.Execute("POST", joinURL, payload)
}

//
// Read Body of Response:
//

func getJsonResponseMessage(response *http.Response) string {
	type Json struct {
		Response string
	}

	var message Json

	if body, err := ioutil.ReadAll(response.Body); err != nil {
		log.Fatal(err)
	} else if err := json.Unmarshal(body, &message); err != nil {
		log.Fatal(err)
	}

	return message.Response
}

func getJsonUserData(response *http.Response) string {
	type Json struct {
		Handle string
		Name   string
	}

	var user Json

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}

	return user.Handle
}

func getJsonUsersData(response *http.Response) []string {
	type Json struct {
		Handle string
		Joined string
	}

	var users []Json
	data := []map[string]string{}
	var handles []string

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	for i := range data {
		user := Json{data[i]["Handle"], data[i]["Joined"]}
		users = append(users, user)
		handles = append(handles, data[i]["Handle"])
	}

	return handles
}

//
// Read info from headers:
//

func getSessionFromResponse(response *http.Response) string {
	authentication := struct {
		Response  string
		Sessionid string
	}{}
	var (
		body []byte
		err  error
	)
	if body, err = ioutil.ReadAll(response.Body); err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(body, &authentication); err != nil {
		log.Fatal(err)
	}

	return authentication.Sessionid
}

//
// Signup Tests:
//

func (s *TestSuite) TestSignupEmptyHandle(c *C) {
	response, err := postSignup("", "test@test.io", "password1", "password1")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Handle is a required field for signup")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestSignupEmptyEmail(c *C) {
	response, err := postSignup("handleA", "", "password1", "password1")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Email is a required field for signup")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestSignupPasswordMismatch(c *C) {
	response, err := postSignup("handleA", "handleA@test.io", "testing777", "testing888")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Passwords do not match")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestSignupPasswordTooShort(c *C) {
	entry := "testing"

	for i := len(entry); i >= 0; i-- {
		pass := entry[:len(entry)-i]
		response, err := postSignup("handleA", "test@test.io", pass, pass)
		if err != nil {
			c.Error(err)
		}

		c.Check(getJsonResponseMessage(response), Equals, "Passwords must be at least 8 characters long")
		c.Assert(response.StatusCode, Equals, 400, Commentf("Password length = %d.", len(entry)-i))
	}
}

func (s *TestSuite) TestSignupHandleTaken(c *C) {
	postSignup("handleA", "test@test.io", "password1", "password1")

	response, err := postSignup("handleA", "b@test.io", "password2", "password2")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Sorry, handle or email is already taken")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestSignupEmailTaken(c *C) {
	postSignup("handleA", "test@test.io", "password1", "password1")

	response, err := postSignup("handleB", "test@test.io", "password2", "password2")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Sorry, handle or email is already taken")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestSignupCreated(c *C) {
	response, err := postSignup("handleA", "test@test.io", "password1", "password1")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Signed up a new user!")
	c.Assert(response.StatusCode, Equals, 201)
}

//
// Login Tests:
//

func (s *TestSuite) TestLoginUserNoExist(c *C) {
	response, err := postSessions("handleA", "password1")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Invalid username or password, please try again.")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestLoginInvalidUsername(c *C) {
	postSignup("handleA", "test@test.io", "password1", "password1")

	response, err := postSessions("wrong_username", "password1")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Invalid username or password, please try again.")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestLoginInvalidPassword(c *C) {
	postSignup("handleA", "test@test.io", "password1", "password1")

	response, err := postSessions("handleA", "wrong_password")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Invalid username or password, please try again.")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestLoginOK(c *C) {
	postSignup("handleA", "test@test.io", "password1", "password1")

	response, err := postSessions("handleA", "password1")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Logged in handleA. Note your session id.")
	c.Assert(response.StatusCode, Equals, 200)
}

//
// Logout Tests:
//

func (s *TestSuite) TestLogoutUserNoExist(c *C) {
	response, err := deleteSessions("testing123")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "That user doesn't exist")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestLogoutOK(c *C) {
	postSignup("handleA", "handleA@test.io", "password1", "password1")

	postSessions("handleA", "password1")

	response, err := deleteSessions("handleA")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Goodbye handleA, have a nice day")
	c.Assert(response.StatusCode, Equals, 200)
}

//
// Get User Tests:
//

func (s *TestSuite) TestGetUserNotFound(c *C) {
	response, err := getUser("testing123")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "No results found")
	c.Assert(response.StatusCode, Equals, 404)
}

func (s *TestSuite) TestGetUserOK(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	response, err := getUser("testing123")
	if err != nil {
		c.Error(err)
	}

	handle := getJsonUserData(response)

	c.Check(handle, Equals, "testing123")
	c.Assert(response.StatusCode, Equals, 200)
}

//
// Get Users Tests:
//

func (s *TestSuite) TestGetUsersNotFound(c *C) {
	response, err := getUsers()
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "No results found")
	c.Assert(response.StatusCode, Equals, 404)
}

func (s *TestSuite) TestGetUsersOK(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	postSignup("testing132", "testing132", "testing132", "testing132")
	postSignup("testing213", "testing213", "testing213", "testing213")
	postSignup("testing231", "testing231", "testing231", "testing231")
	postSignup("testing312", "testing312", "testing312", "testing312")
	postSignup("testing321", "testing321", "testing321", "testing321")

	response, err := getUsers()
	if err != nil {
		c.Error(err)
	}

	handles := getJsonUsersData(response)

	c.Check(handles[0], Equals, "testing123")
	c.Check(handles[1], Equals, "testing132")
	c.Check(handles[2], Equals, "testing213")
	c.Check(handles[3], Equals, "testing231")
	c.Check(handles[4], Equals, "testing312")
	c.Check(handles[5], Equals, "testing321")
	c.Assert(response.StatusCode, Equals, 200)
}

//
// Delete User Tests:
//

func (s *TestSuite) TestDeleteUserInvalidUsername(c *C) {
	postSignup("handleA", "handleA@test.io", "password1", "password1")

	response, _ := postSessions("handleA", "password1")
	sessionid := getSessionFromResponse(response)

	response, err := deleteUser("notHandleA", "password1", sessionid)
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Invalid username or password, please try again.")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestDeleteUserInvalidPassword(c *C) {
	postSignup("handleA", "handleA@test.io", "password1", "password1")

	response, _ := postSessions("handleA", "password1")
	sessionid := getSessionFromResponse(response)

	response, err := deleteUser("handleA", "notPassword1", sessionid)
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Invalid username or password, please try again.")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestDeleteUserOK(c *C) {
	postSignup("handleA", "handleA@test.io", "password1", "password1")

	response, _ := postSessions("handleA", "password1")
	sessionid := getSessionFromResponse(response)

	deleteUserResponse, err := deleteUser("handleA", "password1", sessionid)
	if err != nil {
		c.Error(err)
	}
	// TODO check if user really deleted
	// getUserResponse, err := getUser("handleA")
	// if err != nil {
	// 	c.Error(err)
	// }

	c.Check(getJsonResponseMessage(deleteUserResponse), Equals, "Deleted handleA")
	c.Check(deleteUserResponse.StatusCode, Equals, 200)
	// c.Check(getJsonResponseMessage(getUserResponse), Equals, "No results found")
	// c.Assert(getUserResponse.StatusCode, Equals, 404)
}

//
// Post Circles Tests:
//

func (s *TestSuite) TestPostCirclesUserNoExist(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	response, _ := postSessions("testing123", "testing123")
	sessionid := getSessionFromResponse(response)

	deleteUser("testing123", "testing123", sessionid)

	response, err := postCircles("testing123", sessionid, "testing123", true)
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostCirclesUserNoSession(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	response, _ := postSessions("testing123", "testing123")
	sessionid := getSessionFromResponse(response)

	deleteSessions("testing123")

	response, err := postCircles("testing123", sessionid, "testing123", true)
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostCirclesNameReservedGold(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	response, _ := postSessions("testing123", "testing123")
	sessionid := getSessionFromResponse(response)

	response, err := postCircles("testing123", sessionid, "Gold", false)
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Gold is a reserved circle name")
	c.Assert(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestPostCirclesNameReservedBroadcast(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	response, _ := postSessions("testing123", "testing123")
	sessionid := getSessionFromResponse(response)

	response, err := postCircles("testing123", sessionid, "Broadcast", true)
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Broadcast is a reserved circle name")
	c.Assert(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestPostCirclesPublicCircleCreated(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	response, _ := postSessions("testing123", "testing123")
	sessionid := getSessionFromResponse(response)

	response, err := postCircles("testing123", sessionid, "testing123", true)
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Created new circle testing123 for testing123")
	c.Assert(response.StatusCode, Equals, 201)
}

func (s *TestSuite) TestPostCirclesPrivateCircleCreated(c *C) {
	postSignup("handleA", "test@test.io", "password1", "password1")

	response, _ := postSessions("handleA", "password1")
	sessionid := getSessionFromResponse(response)

	response, err := postCircles("handleA", sessionid, "PrivateCircleForA", true)
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Created new circle PrivateCircleForA for handleA")
	c.Assert(response.StatusCode, Equals, 201)
}

//
// Post Block Tests:
//

func (s *TestSuite) TestPostBlockUserNoExist(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	postSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := postSessions("testing123", "testing123")
	sessionid := getSessionFromResponse(response)

	deleteUser("testing123", "testing123", sessionid)

	response, err := postBlock("testing123", sessionid, "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostBlockTargetNoExist(c *C) {
	postSignup("handleA", "test@test.io", "password1", "password1")

	response, _ := postSessions("handleA", "password1")
	sessionid := getSessionFromResponse(response)

	response, err := postBlock("handleA", sessionid, "handleB")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Bad request, user handleB wasn't found")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostBlockUserNoSession(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	postSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := postSessions("testing123", "testing123")
	sessionid := getSessionFromResponse(response)

	deleteSessions("testing123")

	response, err := postBlock("testing123", sessionid, "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostBlockOK(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	postSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := postSessions("testing123", "testing123")
	sessionid := getSessionFromResponse(response)

	response, err := postBlock("testing123", sessionid, "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "User testing321 has been blocked")
	c.Assert(response.StatusCode, Equals, 200)
}

//
// Post Join Default Tests:
//

func (s *TestSuite) TestPostJoinDefaultUserNoSession(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	postSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := postSessions("testing123", "testing123")
	sessionid := getSessionFromResponse(response)

	deleteSessions("testing123")

	response, err := postJoinDefault("testing123", sessionid, "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostJoinDefaultTargetNoExist(c *C) {
	postSignup("handleA", "test@test.io", "password1", "password1")

	response, _ := postSessions("handleA", "password1")
	sessionid := getSessionFromResponse(response)

	response, err := postJoinDefault("handleA", sessionid, "handleB")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Bad request, user handleB wasn't found")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostJoinDefaultUserBlocked(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	postSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := postSessions("testing321", "testing321")
	sessionid := getSessionFromResponse(response)

	postBlock("testing321", sessionid, "testing123")

	response, _ = postSessions("testing123", "testing123")
	sessionid = getSessionFromResponse(response)

	response, err := postJoinDefault("testing123", sessionid, "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Server refusal to comply with join request")
	c.Assert(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestPostJoinDefaultCreated(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	postSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := postSessions("testing123", "testing123")
	sessionid := getSessionFromResponse(response)

	response, err := postJoinDefault("testing123", sessionid, "testing123")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "JoinDefault request successful!")
	c.Assert(response.StatusCode, Equals, 201)
}

//
// Post Join Tests:
//

func (s *TestSuite) TestPostJoinUserNoSession(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	postSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := postSessions("testing321", "testing321")
	sessionid := getSessionFromResponse(response)

	postCircles("testing321", sessionid, "testing321", true)

	response, _ = postSessions("testing123", "testing123")
	sessionid = getSessionFromResponse(response)

	deleteSessions("testing123")

	response, err := postJoin("testing123", sessionid, "testing321", "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostJoinTargetNoExist(c *C) {
	postSignup("handleA", "handleA@test.io", "password1", "password1")
	postSignup("handleB", "handleB@test.io", "password2", "password2")

	response, _ := postSessions("handleB", "password2")
	sessionid := getSessionFromResponse(response)

	postCircles("handleB", sessionid, "CircleOfB", true)

	response, _ = postSessions("handleA", "password1")
	sessionid = getSessionFromResponse(response)

	response, err := postJoin("handleA", sessionid, "handleC", "CircleOfB")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Bad request, user handleC wasn't found")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostJoinUserBlocked(c *C) {
	postSignup("handleA", "handleA@test.io", "password1", "password1")
	postSignup("handleB", "handleB@test.io", "password2", "password2")

	response, _ := postSessions("handleB", "password2")
	sessionid := getSessionFromResponse(response)

	postCircles("handleB", sessionid, "CircleOfHandleB", true)

	postBlock("handleB", sessionid, "handleA")

	response, _ = postSessions("handleA", "password1")
	sessionid = getSessionFromResponse(response)

	response, err := postJoin("handleA", sessionid, "handleB", "CircleOfHandleB")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Server refusal to comply with join request")
	c.Assert(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestPostJoinCircleNoExist(c *C) {
	postSignup("handleA", "testA@test.io", "password1", "password1")
	postSignup("handleB", "testB@test.io", "password2", "password2")

	response, _ := postSessions("handleA", "password1")
	sessionid := getSessionFromResponse(response)

	response, err := postJoin("handleA", sessionid, "handleB", "NonExistentCircle")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Could not find target circle, join failed")
	c.Assert(response.StatusCode, Equals, 404)
}

func (s *TestSuite) TestPostJoinCreated(c *C) {
	postSignup("handleA", "testA@test.io", "password1", "password1")
	postSignup("handleB", "testB@test.io", "password2", "password2")

	response_B, _ := postSessions("handleB", "password2")
	sessionid_B := getSessionFromResponse(response_B)

	postCircles("handleB", sessionid_B, "MyCircle", true)

	response_A, _ := postSessions("handleA", "password1")
	sessionid_A := getSessionFromResponse(response_A)

	response, err := postJoin("handleA", sessionid_A, "handleB", "MyCircle")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Join request successful!")
	c.Assert(response.StatusCode, Equals, 201)
}
