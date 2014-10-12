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

	signupURL   string
	loginURL    string
	logoutURL   string
	userURL     string
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
	a.Svc.FreshInitialState()
}

/*
	Run once after all tests or benchmarks have finished running.
*/
func (s *TestSuite) TearDownSuite(c *C) {
	server.Close()
}

/*
	Send/Receive Calls to API:
*/

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

func getUser(handle string) (*http.Response, error) {
	user := "{\"Handle\": \"" + handle + "\"}"

	reader = strings.NewReader(user)

	request, err := http.NewRequest("GET", userURL, reader)
	if err != nil {
		log.Fatal(err)
	}

	response, err := http.DefaultClient.Do(request)

	return response, err
}

func getUsers() (*http.Response, error) {
	request, err := http.NewRequest("GET", usersURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	response, err := http.DefaultClient.Do(request)

	return response, err
}

func deleteUser(handle string, password string) (*http.Response, error) {
	credentials := "{\"Handle\": \"" + handle + "\", \"Password\": \"" + password + "\"}"

	reader = strings.NewReader(credentials)

	request, err := http.NewRequest("DELETE", userURL, reader)
	if err != nil {
		log.Fatal(err)
	}

	response, err := http.DefaultClient.Do(request)

	return response, err
}

func postCircles(handle string, sessionId string, circleName string, public bool) (*http.Response, error) {
	payload := "{\"Handle\": \"" + handle + "\", \"SessionId\": \"" + sessionId + "\", \"CircleName\": \"" + circleName + "\", \"Public\": \"" + strconv.FormatBool(public) + "\"}"

	reader = strings.NewReader(payload)

	request, err := http.NewRequest("POST", circlesURL, reader)
	if err != nil {
		log.Fatal(err)
	}

	response, err := http.DefaultClient.Do(request)

	return response, err
}

/*
	Read Body of Response:
*/
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

func getJsonErrorMessage(response *http.Response) string {
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

func getJsonAuthenticationData(response *http.Response) (string, string) {
	type Json struct {
		Response string
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

func getJsonUserData(response *http.Response) (string, string, string) {
	type Json struct {
		Handle   string
		Email    string
		Password string
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

	return user.Handle, user.Email, user.Password
}

func getJsonUsersData(response *http.Response) ([]string) {
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

/*
	Signup Tests:
*/

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

	c.Check(getJsonResponseMessage(response), Equals, "Sorry, handle or email is already taken")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestSignupEmailTaken(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	response, err := postSignup("testing321", "testing123", "testing123", "testing123")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(response), Equals, "Sorry, handle or email is already taken")
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

/*
	Login Tests:
*/

func (s *TestSuite) TestLoginUserNoExist(c *C) {
	response, err := postLogin("testing123", "testing123")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonErrorMessage(response), Equals, "Invalid username or password, please try again.")
	c.Assert(response.StatusCode, Equals, 400)
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

/*
	Logout Tests:
*/

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

/*
	Get User Tests:
*/

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

	handle, email, password := getJsonUserData(response)

	c.Check(handle, Equals, "testing123")
	c.Check(email, Equals, "testing123")
	c.Check(password, Equals, "testing123")
	c.Assert(response.StatusCode, Equals, 200)
}

/*
	Get Users Tests:
*/

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

/*
	Delete User Tests:
*/

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
	postSignup("testing123", "testing123", "testing123", "testing123")

	deleteUserResponse, err := deleteUser("testing123", "testing123")
	if err != nil {
		c.Error(err)
	}

	getUserResponse, err := getUser("testing123")
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonResponseMessage(deleteUserResponse), Equals, "Deleted testing123")
	c.Check(deleteUserResponse.StatusCode, Equals, 200)
	c.Check(getJsonResponseMessage(getUserResponse), Equals, "No results found")
	c.Assert(getUserResponse.StatusCode, Equals, 404)
}

/*
	Post Circles Tests:
*/

func (s *TestSuite) TestPostCirclesUserNoExist(c *C) {
	postSignup("testing123", "testing123", "testing123", "testing123")

	response, _ := postLogin("testing123", "testing123")
	_, sessionId := getJsonAuthenticationData(response)

	deleteUser("testing123", "testing123")

	response, err := postCircles("testing123", sessionId, "testing123", true)
	if err != nil {
		c.Error(err)
	}

	c.Check(getJsonErrorMessage(response), Equals, "Could not authenticate user testing123")
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

	c.Check(getJsonErrorMessage(response), Equals, "Could not authenticate user testing123")
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
	// create user, login user, create circle (name = "Broadcast") => error
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