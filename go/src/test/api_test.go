package api_test

import (
	api "../api"
	routes "../routes"
	"./helper"
	requester "./requester"
	"encoding/json"
	"flag"
	"github.com/jadengore/goconfig"
	. "gopkg.in/check.v1"
	"io"
	"log"
	"net/http/httptest"
	"testing"
	"time"
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

//
// Signup Tests:
//

func (s *TestSuite) TestSignupEmptyHandle(c *C) {
	response, err := req.PostSignup("", "test@test.io", "password1", "password1")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Handle is a required field for signup")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestSignupEmptyEmail(c *C) {
	response, err := req.PostSignup("handleA", "", "password1", "password1")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Email is a required field for signup")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestSignupPasswordMismatch(c *C) {
	response, err := req.PostSignup("handleA", "handleA@test.io", "testing777", "testing888")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Passwords do not match")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestSignupPasswordTooShort(c *C) {
	entry := "testing"

	for i := len(entry); i >= 0; i-- {
		pass := entry[:len(entry)-i]
		response, err := req.PostSignup("handleA", "test@test.io", pass, pass)
		if err != nil {
			c.Error(err)
		}

		c.Check(helper.GetJsonResponseMessage(response), Equals, "Passwords must be at least 8 characters long")
		c.Assert(response.StatusCode, Equals, 400, Commentf("Password length = %d.", len(entry)-i))
	}
}

func (s *TestSuite) TestSignupHandleTaken(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")

	response, err := req.PostSignup("handleA", "b@test.io", "password2", "password2")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Sorry, handle or email is already taken")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestSignupEmailTaken(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")

	response, err := req.PostSignup("handleB", "test@test.io", "password2", "password2")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Sorry, handle or email is already taken")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestSignupCreated(c *C) {
	response, err := req.PostSignup("handleA", "test@test.io", "password1", "password1")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Signed up a new user!")
	c.Assert(response.StatusCode, Equals, 201)
}

//
// Login Tests:
//

func (s *TestSuite) TestLoginUserNoExist(c *C) {
	response, err := req.PostSessions("handleA", "password1")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Invalid username or password, please try again.")
	c.Assert(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestLoginInvalidUsername(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")

	response, err := req.PostSessions("wrong_username", "password1")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Invalid username or password, please try again.")
	c.Assert(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestLoginInvalidPassword(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")

	response, err := req.PostSessions("handleA", "wrong_password")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Invalid username or password, please try again.")
	c.Assert(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestLoginOK(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")

	response, err := req.PostSessions("handleA", "password1")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Logged in handleA. Note your session id.")
	c.Assert(response.StatusCode, Equals, 200)
}

//
// Logout Tests:
//

func (s *TestSuite) TestLogoutUserNoExist(c *C) {
	response, err := req.DeleteSessions("testing123")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Cannot invalidate token because it is missing")
	c.Assert(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestLogoutOK(c *C) {
	req.PostSignup("handleA", "handleA@test.io", "password1", "password1")

	req.PostSessions("handleA", "password1")

	response, err := req.DeleteSessions("handleA")
	if err != nil {
		c.Error(err)
	}

	c.Assert(response.StatusCode, Equals, 204)
}

//
// ChangePassword Tests:
//

func (s *TestSuite) TestChangePasswordUserNoExist(c *C) {
	response, err := req.PostChangePassword("handleA", "1g2fg3j4", "password1", "password123", "password123")
	if err != nil {
		c.Error(err)
	}
	//fmt.Println(helper.GetJsonResponseMessage(response))
	c.Check(helper.GetJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Assert(response.StatusCode, Equals, 401)
}

func (s *TestSuite) TestChangePasswordSamePassword(c *C) {
	req.PostSignup("handleA", "handleA@test.io", "password1", "password1")

	response, _ := req.PostSessions("handleA", "password1")
	sessionid := helper.GetSessionFromResponse(response)
	response, err := req.PostChangePassword("handleA", sessionid, "password1", "password1", "password1")
	if err != nil {
		c.Error(err)
	}
	c.Check(helper.GetJsonResponseMessage(response), Equals, "Current/new password are same, please provide a new password.")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestChangePasswordOK(c *C) {
	req.PostSignup("handleA", "handleA@test.io", "password1", "password1")

	sessionResponse, _ := req.PostSessions("handleA", "password1")
	sessionid := helper.GetSessionFromResponse(sessionResponse)
	response, err := req.PostChangePassword("handleA", sessionid, "password1", "password2", "password2")
	if err != nil {
		c.Error(err)
	}

	c.Assert(response.StatusCode, Equals, 204)
}

//
// Get User Tests:
//

// func (s *TestSuite) TestGetUserNotFound(c *C) {
// 	response, err := getUser("testing123")
// 	if err != nil {
// 		c.Error(err)
// 	}

// 	c.Check(helper.GetJsonResponseMessage(response), Equals, "No results found")
// 	c.Assert(response.StatusCode, Equals, 404)
// }

// func (s *TestSuite) TestGetUserOK(c *C) {
// 	req.PostSignup("testing123", "testing123", "testing123", "testing123")

// 	response, err := getUser("testing123")
// 	if err != nil {
// 		c.Error(err)
// 	}

// 	handle := getJsonUserData(response)

// 	c.Check(handle, Equals, "testing123")
// 	c.Assert(response.StatusCode, Equals, 200)
// }

//
// Get Users Tests:
//

func (s *TestSuite) TestSearchUsersOK(c *C) {
	req.PostSignup("cat", "test1@test.io", "testing123", "testing123")
	req.PostSignup("bat", "test2@test.io", "testing132", "testing132")
	req.PostSignup("cat_woman", "test3@test.io", "testing213", "testing213")
	req.PostSignup("catsawesome", "test4@test.io", "testing231", "testing231")
	req.PostSignup("smart", "test5@test.io", "testing312", "testing312")
	req.PostSignup("battle", "test6@test.io", "testing321", "testing321")

	if response, err := req.SearchForUsers("", "cat", 0, 10, "handle"); err != nil {
		c.Error(err)
	} else {
		data := struct {
			Results  string
			Response string
			Reason   string
			Count    int
		}{}
		helper.Unmarshal(response, &data)
		type UserResult struct {
			Handle string
			Name   string
			Id     int
		}

		results := make([]UserResult, 0)
		json.Unmarshal([]byte(data.Results), &results)
		c.Check(data.Count, Equals, 3)
		c.Check(data.Response, Equals, "Search complete")
		c.Check(data.Reason, Equals, "")
		c.Check(len(results), Equals, 3)
	}

	// handles := getJsonUsersData(response)

	// c.Check(handles[1], Equals, "testing132")
	// c.Check(handles[2], Equals, "testing213")
	// c.Check(handles[3], Equals, "testing231")
	// c.Check(handles[4], Equals, "testing312")
	// c.Check(handles[5], Equals, "testing321")
	// c.Assert(response.StatusCode, Equals, 200)
}

//
// Delete User Tests:
//

func (s *TestSuite) TestDeleteUserInvalidUsername(c *C) {
	req.PostSignup("handleA", "handleA@test.io", "password1", "password1")

	response, _ := req.PostSessions("handleA", "password1")
	sessionid := helper.GetSessionFromResponse(response)

	response, err := req.DeleteUser("notHandleA", "password1", sessionid)
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Invalid username or password, please try again.")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestDeleteUserInvalidPassword(c *C) {
	req.PostSignup("handleA", "handleA@test.io", "password1", "password1")

	response, _ := req.PostSessions("handleA", "password1")
	sessionid := helper.GetSessionFromResponse(response)

	response, err := req.DeleteUser("handleA", "notpassword1", sessionid)
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Invalid username or password, please try again.")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestDeleteUserOK(c *C) {
	req.PostSignup("handleA", "handleA@test.io", "password1", "password1")

	response, _ := req.PostSessions("handleA", "password1")
	sessionid := helper.GetSessionFromResponse(response)

	deleteUserResponse, err := req.DeleteUser("handleA", "password1", sessionid)
	if err != nil {
		c.Error(err)
	}
	// TODO check if user really deleted
	// getUserResponse, err := getUser("handleA")
	// if err != nil {
	// 	c.Error(err)
	// }

	c.Check(deleteUserResponse.StatusCode, Equals, 204)
	// c.Check(helper.GetJsonResponseMessage(getUserResponse), Equals, "No results found")
	// c.Assert(getUserResponse.StatusCode, Equals, 404)
}

//
// Get Authored Messages Tests
//
func (s *TestSuite) TestGetAuthoredMessagesInvalidAuth(c *C) {
	req.PostSignup("handleA", "handleA@test.io", "password1", "password1")
	res, _ := req.GetAuthoredMessages("handleA", "")
	c.Check(res.StatusCode, Equals, 401)
}

func (s *TestSuite) TestGetAuthoredMessagesOK(c *C) {
	req.PostSignup("handleA", "handleA@test.io", "password1", "password1")

	response, _ := req.PostSessions("handleA", "password1")
	sessionid_A := helper.GetSessionFromResponse(response)

	req.PostMessages("Go is going gophers!", sessionid_A)
	req.PostMessages("Hypothesize about stuff", sessionid_A)
	req.PostMessages("The nearest exit may be behind you", sessionid_A)
	req.PostMessages("I make soap.", sessionid_A)

	res, _ := req.GetAuthoredMessages("handleA", sessionid_A)

	data := struct {
		Response string
		Objects  string
		Count    int
	}{}
	helper.Unmarshal(res, &data)
	type Message struct {
		Id      string
		Author  string
		Content string
		Created time.Time
	}

	objects := make([]Message, 0)
	json.Unmarshal([]byte(data.Objects), &objects)

	c.Check(data.Response, Equals, "Found messages for user handleA")
	c.Check(res.StatusCode, Equals, 200)
	c.Check(data.Count, Equals, 4)
	c.Check(objects[0].Author, Equals, "handleA")
	c.Check(objects[0].Content, Equals, "Go is going gophers!")
	c.Check(objects[1].Content, Equals, "Hypothesize about stuff")
	c.Check(objects[2].Content, Equals, "The nearest exit may be behind you")
	c.Check(objects[3].Content, Equals, "I make soap.")
}

//
// Post Circles Tests:
//

func (s *TestSuite) TestPostCirclesUserNoExist(c *C) {
	req.PostSignup("testing123", "testing123", "testing123", "testing123")

	response, _ := req.PostSessions("testing123", "testing123")
	sessionid := helper.GetSessionFromResponse(response)

	req.DeleteUser("testing123", "testing123", sessionid)

	response, err := req.PostCircles("testing123", sessionid, "testing123", true)
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostCirclesUserNoSession(c *C) {
	req.PostSignup("testing123", "testing123", "testing123", "testing123")

	response, _ := req.PostSessions("testing123", "testing123")
	sessionid := helper.GetSessionFromResponse(response)

	req.DeleteSessions("testing123")

	response, err := req.PostCircles("testing123", sessionid, "testing123", true)
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostCirclesNameReservedGold(c *C) {
	req.PostSignup("testing123", "testing123", "testing123", "testing123")

	response, _ := req.PostSessions("testing123", "testing123")
	sessionid := helper.GetSessionFromResponse(response)

	response, err := req.PostCircles("testing123", sessionid, "Gold", false)
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Gold is a reserved circle name")
	c.Assert(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestPostCirclesNameReservedBroadcast(c *C) {
	req.PostSignup("testing123", "testing123", "testing123", "testing123")

	response, _ := req.PostSessions("testing123", "testing123")
	sessionid := helper.GetSessionFromResponse(response)

	response, err := req.PostCircles("testing123", sessionid, "Broadcast", true)
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Broadcast is a reserved circle name")
	c.Assert(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestPostCirclesPublicCircleCreated(c *C) {
	req.PostSignup("testing123", "testing123", "testing123", "testing123")

	response, _ := req.PostSessions("testing123", "testing123")
	sessionid := helper.GetSessionFromResponse(response)

	response, err := req.PostCircles("testing123", sessionid, "testing123", true)
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Created new circle testing123 for testing123")
	c.Assert(response.StatusCode, Equals, 201)
}

func (s *TestSuite) TestPostCirclesPrivateCircleCreated(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")

	response, _ := req.PostSessions("handleA", "password1")
	sessionid := helper.GetSessionFromResponse(response)

	response, err := req.PostCircles("handleA", sessionid, "PrivateCircleForA", true)
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Created new circle PrivateCircleForA for handleA")
	c.Assert(response.StatusCode, Equals, 201)
}

//
// Post Block Tests:
//

func (s *TestSuite) TestPostBlockUserNoExist(c *C) {
	req.PostSignup("testing123", "testing123", "testing123", "testing123")
	req.PostSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := req.PostSessions("testing123", "testing123")
	sessionid := helper.GetSessionFromResponse(response)

	req.DeleteUser("testing123", "testing123", sessionid)

	response, err := req.PostBlock("testing123", sessionid, "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Assert(response.StatusCode, Equals, 401)
}

func (s *TestSuite) TestPostBlockTargetNoExist(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")

	response, _ := req.PostSessions("handleA", "password1")
	sessionid := helper.GetSessionFromResponse(response)

	response, err := req.PostBlock("handleA", sessionid, "handleB")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Bad request, user handleB wasn't found")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostBlockUserNoSession(c *C) {
	req.PostSignup("testing123", "testing123", "testing123", "testing123")
	req.PostSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := req.PostSessions("testing123", "testing123")
	sessionid := helper.GetSessionFromResponse(response)

	req.DeleteSessions("testing123")

	response, err := req.PostBlock("testing123", sessionid, "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Assert(response.StatusCode, Equals, 401)
}

func (s *TestSuite) TestPostBlockOK(c *C) {
	req.PostSignup("testing123", "testing123", "testing123", "testing123")
	req.PostSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := req.PostSessions("testing123", "testing123")
	sessionid := helper.GetSessionFromResponse(response)

	response, err := req.PostBlock("testing123", sessionid, "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "User testing321 has been blocked")
	c.Assert(response.StatusCode, Equals, 200)
}

//
// Post Join Default Tests:
//

func (s *TestSuite) TestPostJoinDefaultUserNoSession(c *C) {
	req.PostSignup("testing123", "testing123", "testing123", "testing123")
	req.PostSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := req.PostSessions("testing123", "testing123")
	sessionid := helper.GetSessionFromResponse(response)

	req.DeleteSessions("testing123")

	response, err := req.PostJoinDefault("testing123", sessionid, "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Assert(response.StatusCode, Equals, 401)
}

func (s *TestSuite) TestPostJoinDefaultTargetNoExist(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")

	response, _ := req.PostSessions("handleA", "password1")
	sessionid := helper.GetSessionFromResponse(response)

	response, err := req.PostJoinDefault("handleA", sessionid, "handleB")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Bad request, user handleB wasn't found")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostJoinDefaultUserBlocked(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	response, _ := req.PostSessions("handleB", "password2")
	sessionid := helper.GetSessionFromResponse(response)

	req.PostBlock("handleB", sessionid, "handleA")

	response, _ = req.PostSessions("handleA", "password1")
	sessionid = helper.GetSessionFromResponse(response)

	response, err := req.PostJoinDefault("handleA", sessionid, "handleB")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Server refusal to comply with join request")
	c.Assert(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestPostJoinDefaultCreated(c *C) {
	req.PostSignup("testing123", "testing123", "testing123", "testing123")
	req.PostSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := req.PostSessions("testing123", "testing123")
	sessionid := helper.GetSessionFromResponse(response)

	response, err := req.PostJoinDefault("testing123", sessionid, "testing123")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "JoinDefault request successful!")
	c.Assert(response.StatusCode, Equals, 201)
}

//
// Post Join Tests:
//

func (s *TestSuite) TestPostJoinUserNoSession(c *C) {
	req.PostSignup("testing123", "testing123", "testing123", "testing123")
	req.PostSignup("testing321", "testing321", "testing321", "testing321")

	response, _ := req.PostSessions("testing321", "testing321")
	sessionid := helper.GetSessionFromResponse(response)

	req.PostCircles("testing321", sessionid, "testing321", true)

	response, _ = req.PostSessions("testing123", "testing123")
	sessionid = helper.GetSessionFromResponse(response)

	req.DeleteSessions("testing123")

	response, err := req.PostJoin("testing123", sessionid, "testing321", "testing321")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Assert(response.StatusCode, Equals, 401)
}

func (s *TestSuite) TestPostJoinTargetNoExist(c *C) {
	req.PostSignup("handleA", "handleA@test.io", "password1", "password1")
	req.PostSignup("handleB", "handleB@test.io", "password2", "password2")

	response, _ := req.PostSessions("handleB", "password2")
	sessionid := helper.GetSessionFromResponse(response)

	req.PostCircles("handleB", sessionid, "CircleOfB", true)

	response, _ = req.PostSessions("handleA", "password1")
	sessionid = helper.GetSessionFromResponse(response)

	response, err := req.PostJoin("handleA", sessionid, "handleC", "CircleOfB")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Bad request, user handleC wasn't found")
	c.Assert(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostJoinUserBlocked(c *C) {
	req.PostSignup("handleA", "handleA@test.io", "password1", "password1")
	req.PostSignup("handleB", "handleB@test.io", "password2", "password2")

	response, _ := req.PostSessions("handleB", "password2")
	sessionid := helper.GetSessionFromResponse(response)

	req.PostCircles("handleB", sessionid, "CircleOfHandleB", true)

	req.PostBlock("handleB", sessionid, "handleA")

	response, _ = req.PostSessions("handleA", "password1")
	sessionid = helper.GetSessionFromResponse(response)

	response, err := req.PostJoin("handleA", sessionid, "handleB", "CircleOfHandleB")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Server refusal to comply with join request")
	c.Assert(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestPostJoinCircleNoExist(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	response, _ := req.PostSessions("handleA", "password1")
	sessionid := helper.GetSessionFromResponse(response)

	response, err := req.PostJoin("handleA", sessionid, "handleB", "NonExistentCircle")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Could not find target circle, join failed")
	c.Assert(response.StatusCode, Equals, 404)
}

func (s *TestSuite) TestPostJoinCreated(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	response_B, _ := req.PostSessions("handleB", "password2")
	sessionid_B := helper.GetSessionFromResponse(response_B)

	req.PostCircles("handleB", sessionid_B, "MyCircle", true)

	response_A, _ := req.PostSessions("handleA", "password1")
	sessionid_A := helper.GetSessionFromResponse(response_A)

	response, err := req.PostJoin("handleA", sessionid_A, "handleB", "MyCircle")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Join request successful!")
	c.Assert(response.StatusCode, Equals, 201)
}
