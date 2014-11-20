package api_test

import (
	"./helper"
	. "gopkg.in/check.v1"
)

//
// Signup Tests:
//

func (s *TestSuite) TestSignupEmptyHandle(c *C) {
	response, err := req.PostSignup("", "test@test.io", "password1", "password1")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonReasonMessage(response), Equals, "Handle is a required field for signup")
	c.Check(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestSignupEmptyEmail(c *C) {
	response, err := req.PostSignup("handleA", "", "password1", "password1")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonReasonMessage(response), Equals, "Email is a required field for signup")
	c.Check(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestSignupPasswordMismatch(c *C) {
	response, err := req.PostSignup("handleA", "handleA@test.io", "password1", "password2")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonReasonMessage(response), Equals, "Passwords do not match")
	c.Check(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestSignupPasswordTooShort(c *C) {
	entry := "testing"

	for i := len(entry); i >= 0; i-- {
		pass := entry[:len(entry)-i]
		response, err := req.PostSignup("handleA", "test@test.io", pass, pass)
		if err != nil {
			c.Error(err)
		}

		c.Check(helper.GetJsonReasonMessage(response), Equals, "Passwords must be at least 8 characters long")
		c.Check(response.StatusCode, Equals, 403, Commentf("Password length = %d.", len(entry)-i))
	}
}

func (s *TestSuite) TestSignupHandleTaken(c *C) {
	req.PostSignup("handleA", "handleA@test.io", "password1", "password1")

	response, err := req.PostSignup("handleA", "handleB@test.io", "password2", "password2")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonReasonMessage(response), Equals, "Sorry, handle or email is already taken")
	c.Check(response.StatusCode, Equals, 409)
}

func (s *TestSuite) TestSignupEmailTaken(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")

	response, err := req.PostSignup("handleB", "test@test.io", "password2", "password2")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonReasonMessage(response), Equals, "Sorry, handle or email is already taken")
	c.Check(response.StatusCode, Equals, 409)
}

func (s *TestSuite) TestSignupCreated(c *C) {
	response, err := req.PostSignup("handleA", "test@test.io", "password1", "password1")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Signed up a new user!")
	c.Check(response.StatusCode, Equals, 201)
}

//
// Login Tests:
//

func (s *TestSuite) TestLoginUserNoExist(c *C) {
	response, err := req.PostSessions("handleA", "password1")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonReasonMessage(response), Equals, "Invalid username or password, please try again.")
	c.Check(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestLoginInvalidUsername(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")

	response, err := req.PostSessions("wrong_username", "password1")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonReasonMessage(response), Equals, "Invalid username or password, please try again.")
	c.Check(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestLoginInvalidPassword(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")

	response, err := req.PostSessions("handleA", "wrong_password")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonReasonMessage(response), Equals, "Invalid username or password, please try again.")
	c.Check(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestLoginOK(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")

	response, err := req.PostSessions("handleA", "password1")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Logged in handleA. Note your Authorization token.")
	c.Check(response.StatusCode, Equals, 201)
}

//
// Logout Tests:
//

func (s *TestSuite) TestLogoutUserNoExist(c *C) {
	response, err := req.DeleteSessions("token-of-noone")
	if err != nil {
		c.Error(err)
		c.Error(response)
	}

	c.Check(response.StatusCode, Equals, 403)
	c.Check(helper.GetJsonReasonMessage(response), Equals, "Cannot invalidate token because it is missing")
}

func (s *TestSuite) TestLogoutOK(c *C) {
	req.PostSignup("handleA", "handleA@test.io", "password1", "password1")

	sessionid_A := req.PostSessionGetAuthToken("handleA", "password1")

	response, err := req.DeleteSessions(sessionid_A)
	if err != nil {
		c.Error(err)
	}

	c.Check(response.StatusCode, Equals, 204)
}

//
// ChangePassword Tests:
//

func (s *TestSuite) TestChangePasswordUserNoExist(c *C) {
	response, err := req.PostChangePassword("SomeSessionId", "password1", "password123", "password123")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Check(response.StatusCode, Equals, 401)
}

func (s *TestSuite) TestChangePasswordSamePassword(c *C) {
	req.PostSignup("handleA", "handleA@test.io", "password1", "password1")

	sessionid := req.PostSessionGetAuthToken("handleA", "password1")
	response, err := req.PostChangePassword(sessionid, "password1", "password1", "password1")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonReasonMessage(response), Equals, "Current/new password are same, please provide a new password.")
	c.Check(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestChangePasswordOK(c *C) {
	req.PostSignup("handleA", "handleA@test.io", "password1", "password1")

	sessionid := req.PostSessionGetAuthToken("handleA", "password1")
	response, err := req.PostChangePassword(sessionid, "password1", "password2", "password2")
	if err != nil {
		c.Error(err)
	}

	c.Check(response.StatusCode, Equals, 204)
}
