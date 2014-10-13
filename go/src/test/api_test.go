package api_test

import (
	api "../api"
	routes "../routes"
	"encoding/json"
	"fmt"
	. "gopkg.in/check.v1"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	//"time"
)

var (
	server *httptest.Server
	a      *api.Api

	signupURL      string
	loginURL       string
	logoutURL      string
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
	uri := "http://localhost:7474/db/data"

	a = api.NewApi(uri)

	handler, err := routes.MakeHandler(*a)
	if err != nil {
		log.Fatal(err)
	}

	server = httptest.NewServer(&handler)

	signupURL = fmt.Sprintf("%s/signup", server.URL)
	loginURL = fmt.Sprintf("%s/login", server.URL)
	logoutURL = fmt.Sprintf("%s/logout", server.URL)
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
	proposal := "{\"Handle\": \"" + handle + "\", \"Email\": \"" + email + "\", \"Password\": \"" + password + "\", \"ConfirmPassword\": \"" + confirmPassword + "\"}"

	reader = strings.NewReader(proposal)

	request, err := http.NewRequest("POST", signupURL, reader)
	if err != nil {
		log.Fatal(err)
	}

	return http.DefaultClient.Do(request)
}

func postLogin(handle string, password string) (*http.Response, error) {
	credentials := "{\"Handle\": \"" + handle + "\", \"Password\": \"" + password + "\"}"

	reader = strings.NewReader(credentials)

	request, err := http.NewRequest("POST", loginURL, reader)
	if err != nil {
		log.Fatal(err)
	}

	return http.DefaultClient.Do(request)
}

func postLogout(handle string) (*http.Response, error) {
	user := "{\"Handle\": \"" + handle + "\"}"

	reader = strings.NewReader(user)

	request, err := http.NewRequest("POST", logoutURL, reader)
	if err != nil {
		log.Fatal(err)
	}

	return http.DefaultClient.Do(request)
}

func getUser(handle string) (*http.Response, error) {
	user := "{\"Handle\": \"" + handle + "\"}"

	reader = strings.NewReader(user)

	request, err := http.NewRequest("GET", userURL, reader)
	if err != nil {
		log.Fatal(err)
	}

	return http.DefaultClient.Do(request)
}

func getUsers() (*http.Response, error) {
	request, err := http.NewRequest("GET", usersURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	return http.DefaultClient.Do(request)
}

func deleteUser(handle string, password string) (*http.Response, error) {
	credentials := "{\"Handle\": \"" + handle + "\", \"Password\": \"" + password + "\"}"

	reader = strings.NewReader(credentials)

	request, err := http.NewRequest("DELETE", userURL, reader)
	if err != nil {
		log.Fatal(err)
	}

	return http.DefaultClient.Do(request)
}

func postCircles(handle string, sessionId string, circleName string, public bool) (*http.Response, error) {
	payload := "{\"Handle\": \"" + handle + "\", \"SessionId\": \"" + sessionId + "\", \"CircleName\": \"" + circleName + "\", \"Public\": \"" + strconv.FormatBool(public) + "\"}"

	reader = strings.NewReader(payload)

	request, err := http.NewRequest("POST", circlesURL, reader)
	if err != nil {
		log.Fatal(err)
	}

	return http.DefaultClient.Do(request)
}

func postBlock(handle string, sessionId string, target string) (*http.Response, error) {
	payload := "{\"Handle\": \"" + handle + "\", \"SessionId\": \"" + sessionId + "\", \"Target\": \"" + target + "\"}"

	reader = strings.NewReader(payload)

	request, err := http.NewRequest("POST", blockURL, reader)
	if err != nil {
		log.Fatal(err)
	}

	return http.DefaultClient.Do(request)
}

func postJoinDefault(handle string, sessionId string, target string) (*http.Response, error) {
	payload := "{\"Handle\": \"" + handle + "\", \"SessionId\": \"" + sessionId + "\", \"Target\": \"" + target + "\"}"

	reader = strings.NewReader(payload)

	request, err := http.NewRequest("POST", joindefaultURL, reader)
	if err != nil {
		log.Fatal(err)
	}

	return http.DefaultClient.Do(request)
}

func postJoin(handle string, sessionId string, target string, circle string) (*http.Response, error) {
	payload := "{\"Handle\": \"" + handle + "\", \"SessionId\": \"" + sessionId + "\", \"Target\": \"" + target + "\", \"Circle\": \"" + circle + "\" }"

	reader = strings.NewReader(payload)

	request, err := http.NewRequest("POST", joinURL, reader)
	if err != nil {
		log.Fatal(err)
	}

	return http.DefaultClient.Do(request)
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

func getJsonAuthenticationData(response *http.Response) (string, string) {
	type Json struct {
		Response  string
		SessionId string
	}

	var authentication Json

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &authentication)
	if err != nil {
		log.Fatal(err)
	}

	return authentication.Response, authentication.SessionId
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
	response, err := postLogin("handleA", "password1")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Invalid username or password, please try again.")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestLoginInvalidUsername(c *C) {
	postSignup("handleA", "test@test.io", "password1", "password1")

	response, err := postLogin("wrong_username", "password1")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Invalid username or password, please try again.")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestLoginInvalidPassword(c *C) {
	postSignup("handleA", "test@test.io", "password1", "password1")

	response, err := postLogin("handleA", "wrong_password")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Invalid username or password, please try again.")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestLoginOK(c *C) {
	postSignup("handleA", "test@test.io", "password1", "password1")

	response, err := postLogin("handleA", "password1")
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
	response, err := postLogout("testing123")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "That user doesn't exist")
	c.Assert(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestLogoutOK(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	postLogin("testing123", "testing123")

	response, err := postLogout("testing123")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Goodbye testing123, have a nice day")
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

func (s *TestSuite) TestDeleteUserNoExist(c *C) {
	response, err := deleteUser("testing123", "testing123")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Could not delete user with supplied credentials")
	c.Assert(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestDeleteUserInvalidUsername(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	response, err := deleteUser("testing321", "testing123")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Could not delete user with supplied credentials")
	c.Assert(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestDeleteUserInvalidPassword(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	response, err := deleteUser("testing123", "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Could not delete user with supplied credentials")
	c.Assert(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestDeleteUserOK(c *C) {
	postSignup("handleA", "handleA@test.io", "password1", "password1")

	deleteUserResponse, err := deleteUser("handleA", "password1")
	if err != nil {
		c.Error(err)
	}

	getUserResponse, err := getUser("handleA")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(deleteUserResponse), Equals, "Deleted handleA")
	c.Check(deleteUserResponse.StatusCode, Equals, 200)
	c.Check(getJsonResponseMessage(getUserResponse), Equals, "No results found")
	c.Assert(getUserResponse.StatusCode, Equals, 404)
}

//
// Post Circles Tests:
//

func (s *TestSuite) TestPostCirclesUserNoExist(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	response, _ := postLogin("testing123", "testing123")
	_, sessionId := getJsonAuthenticationData(response)

	deleteUser("testing123", "testing123")

	response, err := postCircles("testing123", sessionId, "testing123", true)
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostCirclesUserNoSession(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	response, _ := postLogin("testing123", "testing123")
	_, sessionId := getJsonAuthenticationData(response)

	postLogout("testing123")

	response, err := postCircles("testing123", sessionId, "testing123", true)
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostCirclesNameReservedGold(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	response, _ := postLogin("testing123", "testing123")
	_, sessionId := getJsonAuthenticationData(response)

	response, err := postCircles("testing123", sessionId, "Gold", false)
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Gold is a reserved circle name")
	c.Assert(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestPostCirclesNameReservedBroadcast(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	response, _ := postLogin("testing123", "testing123")
	_, sessionId := getJsonAuthenticationData(response)

	response, err := postCircles("testing123", sessionId, "Broadcast", true)
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Broadcast is a reserved circle name")
	c.Assert(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestPostCirclesPublicCircleCreated(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	response, _ := postLogin("testing123", "testing123")
	_, sessionId := getJsonAuthenticationData(response)

	response, err := postCircles("testing123", sessionId, "testing123", true)
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Created new circle testing123 for testing123")
	c.Assert(response.StatusCode, Equals, 201)
}

func (s *TestSuite) TestPostCirclesPrivateCircleCreated(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	response, _ := postLogin("testing123", "testing123")
	_, sessionId := getJsonAuthenticationData(response)

	response, err := postCircles("testing123", sessionId, "testing123", false)
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Created new circle testing123 for testing123")
	c.Assert(response.StatusCode, Equals, 201)
}

//
// Post Block Tests:
//

func (s *TestSuite) TestPostBlockUserNoExist(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	postSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := postLogin("testing123", "testing123")
	_, sessionId := getJsonAuthenticationData(response)

	deleteUser("testing123", "testing123")

	response, err := postBlock("testing123", sessionId, "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostBlockTargetNoExist(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	postSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := postLogin("testing123", "testing123")
	_, sessionId := getJsonAuthenticationData(response)

	deleteUser("testing321", "testing321")

	response, err := postBlock("testing123", sessionId, "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Bad request, user testing321 wasn't found")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostBlockUserNoSession(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	postSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := postLogin("testing123", "testing123")
	_, sessionId := getJsonAuthenticationData(response)

	postLogout("testing123")

	response, err := postBlock("testing123", sessionId, "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostBlockOK(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	postSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := postLogin("testing123", "testing123")
	_, sessionid := getJsonAuthenticationData(response)

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

func (s *TestSuite) TestPostJoinDefaultUserNoExist(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	postSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := postLogin("testing123", "testing123")
	_, sessionId := getJsonAuthenticationData(response)

	deleteUser("testing123", "testing123")

	response, err := postJoinDefault("testing123", sessionId, "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostJoinDefaultUserNoSession(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	postSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := postLogin("testing123", "testing123")
	_, sessionId := getJsonAuthenticationData(response)

	postLogout("testing123")

	response, err := postJoinDefault("testing123", sessionId, "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostJoinDefaultTargetNoExist(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	postSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := postLogin("testing123", "testing123")
	_, sessionId := getJsonAuthenticationData(response)

	deleteUser("testing321", "testing321")

	response, err := postJoinDefault("testing123", sessionId, "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Bad request, user testing321 wasn't found")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostJoinDefaultUserBlocked(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	postSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := postLogin("testing321", "testing321")
	_, sessionId := getJsonAuthenticationData(response)

	postBlock("testing321", sessionId, "testing123")

	response, _ = postLogin("testing123", "testing123")
	_, sessionId = getJsonAuthenticationData(response)

	response, err := postJoinDefault("testing123", sessionId, "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Server refusal to comply with join request")
	c.Assert(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestPostJoinDefaultCreated(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	postSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := postLogin("testing123", "testing123")
	_, sessionId := getJsonAuthenticationData(response)

	response, err := postJoinDefault("testing123", sessionId, "testing123")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "JoinDefault request successful!")
	c.Assert(response.StatusCode, Equals, 201)
}

//
// Post Join Tests:
//

func (s *TestSuite) TestPostJoinUserNoExist(c *C) {
	postSignup("handleA", "testA@test.io", "password1", "password1")
	postSignup("handleB", "testB@test.io", "password2", "password2")

	response, _ := postLogin("handleB", "password2")
	_, sessionId := getJsonAuthenticationData(response)

	postCircles("handleB", sessionId, "MyPublicCircle", true)

	response, _ = postLogin("handleA", "password1")
	_, sessionId = getJsonAuthenticationData(response)

	deleteUser("handleA", "password1")

	response, err := postJoin("handleA", sessionId, "handleB", "MyPublicCircle")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostJoinUserNoSession(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	postSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := postLogin("testing321", "testing321")
	_, sessionId := getJsonAuthenticationData(response)

	postCircles("testing321", sessionId, "testing321", true)

	response, _ = postLogin("testing123", "testing123")
	_, sessionId = getJsonAuthenticationData(response)

	postLogout("testing123")

	response, err := postJoin("testing123", sessionId, "testing321", "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostJoinTargetNoExist(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	postSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := postLogin("testing321", "testing321")
	_, sessionId := getJsonAuthenticationData(response)

	postCircles("testing321", sessionId, "testing321", true)

	response, _ = postLogin("testing123", "testing123")
	_, sessionId = getJsonAuthenticationData(response)

	deleteUser("testing321", "testing321")

	response, err := postJoin("testing123", sessionId, "testing321", "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Bad request, user testing321 wasn't found")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostJoinUserBlocked(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")
	postSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := postLogin("testing321", "testing321")
	_, sessionId := getJsonAuthenticationData(response)

	postCircles("testing321", sessionId, "testing321", true)

	postBlock("testing321", sessionId, "testing123")

	response, _ = postLogin("testing123", "testing123")
	_, sessionId = getJsonAuthenticationData(response)

	response, err := postJoin("testing123", sessionId, "testing321", "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Server refusal to comply with join request")
	c.Assert(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestPostJoinCircleNoExist(c *C) {
	postSignup("handleA", "testA@test.io", "password1", "password1")
	postSignup("handleB", "testB@test.io", "password2", "password2")

	response, _ := postLogin("handleA", "password1")
	_, sessionId := getJsonAuthenticationData(response)

	response, err := postJoin("handleA", sessionId, "handleB", "NonExistentCircle")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Could not find target circle, join failed")
	c.Assert(response.StatusCode, Equals, 404)
}

func (s *TestSuite) TestPostJoinCreated(c *C) {
	postSignup("handleA", "testA@test.io", "password1", "password1")
	postSignup("handleB", "testB@test.io", "password2", "password2")

	response_B, _ := postLogin("handleB", "password2")
	_, sessionid_B := getJsonAuthenticationData(response_B)

	postCircles("handleB", sessionid_B, "MyCircle", true)

	response_A, _ := postLogin("handleA", "password1")
	_, sessionid_A := getJsonAuthenticationData(response_A)

	response, err := postJoin("handleA", sessionid_A, "handleB", "MyCircle")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Join request successful!")
	c.Assert(response.StatusCode, Equals, 201)
}
