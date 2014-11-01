package api_test

import (
	"./helper"
	. "gopkg.in/check.v1"
)

//
// Post Circles Tests:
//

func (s *TestSuite) TestPostCirclesUserNoExist(c *C) {
	res, err := req.PostCircles("handleNotA", "SomeSessionId", "SomeCircle", true)
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(res), Equals, "Failed to authenticate user request")
	c.Check(res.StatusCode, Equals, 401)
}

func (s *TestSuite) TestPostCirclesUserNoSession(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")

	response, _ := req.PostSessions("handleA", "password1")
	sessionid := helper.GetSessionFromResponse(response)

	req.DeleteSessions("handleA")

	res, err := req.PostCircles("handleA", sessionid, "SomeCircle", true)
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(res), Equals, "Failed to authenticate user request")
	c.Check(res.StatusCode, Equals, 401)
}

func (s *TestSuite) TestPostCirclesNameReservedGold(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")

	response, _ := req.PostSessions("handleA", "password1")
	sessionid := helper.GetSessionFromResponse(response)

	res1, err := req.PostCircles("handleA", sessionid, "Gold", false)
	if err != nil {
		c.Error(err)
	}
	res2, err := req.PostCircles("handleA", sessionid, "Gold", true)
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(res1), Equals, "Gold is a reserved circle name")
	c.Check(res1.StatusCode, Equals, 403)
	c.Check(helper.GetJsonResponseMessage(res2), Equals, "Gold is a reserved circle name")
	c.Check(res2.StatusCode, Equals, 403)
}

func (s *TestSuite) TestPostCirclesNameReservedBroadcast(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")

	response, _ := req.PostSessions("handleA", "password1")
	sessionid := helper.GetSessionFromResponse(response)

	res1, err := req.PostCircles("handleA", sessionid, "Broadcast", false)
	if err != nil {
		c.Error(err)
	}
	res2, err := req.PostCircles("handleA", sessionid, "Broadcast", true)
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(res1), Equals, "Broadcast is a reserved circle name")
	c.Check(res1.StatusCode, Equals, 403)
	c.Check(helper.GetJsonResponseMessage(res2), Equals, "Broadcast is a reserved circle name")
	c.Check(res2.StatusCode, Equals, 403)
}

func (s *TestSuite) TestPostPublicCircleOK(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")

	response, _ := req.PostSessions("handleA", "password1")
	sessionid := helper.GetSessionFromResponse(response)

	res, err := req.PostCircles("handleA", sessionid, "MyPublicCircle", true)
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(res), Equals, "Created new circle MyPublicCircle for handleA")
	c.Check(res.StatusCode, Equals, 201)
}

func (s *TestSuite) TestPostPrivateCircleOK(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")

	response, _ := req.PostSessions("handleA", "password1")
	sessionid := helper.GetSessionFromResponse(response)

	res, err := req.PostCircles("handleA", sessionid, "MyPrivateCircle", true)
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(res), Equals, "Created new circle MyPrivateCircle for handleA")
	c.Check(res.StatusCode, Equals, 201)
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
