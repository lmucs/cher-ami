package api_test

import (
	"./helper"
	. "gopkg.in/check.v1"
)

//
// Post Circles Tests:
//

func (s *TestSuite) TestPostCirclesUserNoExist(c *C) {
	res, err := req.PostCircles("SomeSessionId", "SomeCircle", true)
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(res), Equals, "Failed to authenticate user request")
	c.Check(res.StatusCode, Equals, 401)
}

func (s *TestSuite) TestPostCirclesUserNoSession(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")

	sessionid := req.PostSessionGetSessionId("handleA", "password1")

	req.DeleteSessions(sessionid)

	res, err := req.PostCircles(sessionid, "SomeCircle", true)
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(res), Equals, "Failed to authenticate user request")
	c.Check(res.StatusCode, Equals, 401)
}

func (s *TestSuite) TestPostCirclesNameReservedGold(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")

	sessionid := req.PostSessionGetSessionId("handleA", "password1")

	res1, err := req.PostCircles(sessionid, "Gold", false)
	if err != nil {
		c.Error(err)
	}
	res2, err := req.PostCircles(sessionid, "Gold", true)
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

	sessionid := req.PostSessionGetSessionId("handleA", "password1")

	res1, err := req.PostCircles(sessionid, "Broadcast", false)
	if err != nil {
		c.Error(err)
	}
	res2, err := req.PostCircles(sessionid, "Broadcast", true)
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

	sessionid := req.PostSessionGetSessionId("handleA", "password1")

	res, err := req.PostCircles(sessionid, "MyPublicCircle", true)
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(res), Equals, "Created new circle MyPublicCircle for handleA")
	c.Check(res.StatusCode, Equals, 201)
}

func (s *TestSuite) TestPostPrivateCircleOK(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")

	sessionid := req.PostSessionGetSessionId("handleA", "password1")

	res, err := req.PostCircles(sessionid, "MyPrivateCircle", true)
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
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid := req.PostSessionGetSessionId("handleA", "password1")

	req.DeleteUser("handleA", "password1", sessionid)

	response, err := req.PostBlock(sessionid, "handleB")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Check(response.StatusCode, Equals, 401)
}

func (s *TestSuite) TestPostBlockTargetNoExist(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")

	sessionid := req.PostSessionGetSessionId("handleA", "password1")

	response, err := req.PostBlock(sessionid, "handleB")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Bad request, user handleB wasn't found")
	c.Check(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostBlockUserNoSession(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid := req.PostSessionGetSessionId("handleA", "password1")

	req.DeleteSessions(sessionid)

	response, err := req.PostBlock(sessionid, "handleB")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Check(response.StatusCode, Equals, 401)
}

func (s *TestSuite) TestPostBlockOK(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid := req.PostSessionGetSessionId("handleA", "password1")

	response, err := req.PostBlock(sessionid, "handleB")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "User handleB has been blocked")
	c.Check(response.StatusCode, Equals, 200)
}

//
// Post Join Default Tests:
//

func (s *TestSuite) TestPostJoinDefaultUserNoSession(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid := req.PostSessionGetSessionId("handleA", "password1")

	req.DeleteSessions(sessionid)

	response, err := req.PostBlock(sessionid, "handleB")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Check(response.StatusCode, Equals, 401)
}

func (s *TestSuite) TestPostJoinDefaultTargetNoExist(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")

	response, _ := req.PostSessions("handleA", "password1")
	sessionid := helper.GetSessionFromResponse(response)

	response, err := req.PostJoinDefault(sessionid, "handleB")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Bad request, user handleB wasn't found")
	c.Check(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostJoinDefaultUserBlocked(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid_B := req.PostSessionGetSessionId("handleB", "password2")

	req.PostBlock(sessionid_B, "handleA")

	sessionid_A := req.PostSessionGetSessionId("handleA", "password1")

	response, err := req.PostJoinDefault(sessionid_A, "handleB")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Server refusal to comply with join request")
	c.Check(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestPostJoinDefaultCreated(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid := req.PostSessionGetSessionId("handleA", "password1")

	response, err := req.PostJoinDefault(sessionid, "handleB")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "JoinDefault request successful!")
	c.Check(response.StatusCode, Equals, 201)
}

//
// Post Join Tests:
//

func (s *TestSuite) TestPostJoinUserNoSession(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid_B := req.PostSessionGetSessionId("handleB", "password2")

	req.PostCircles(sessionid_B, "handleB", true)

	sessionid_A := req.PostSessionGetSessionId("handleA", "password1")

	req.DeleteSessions(sessionid_A)

	response, err := req.PostJoin(sessionid_A, "handleB", "CircleOfB")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Check(response.StatusCode, Equals, 401)
}

func (s *TestSuite) TestPostJoinUserNoExist(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid_B := req.PostSessionGetSessionId("handleB", "password2")

	req.PostCircles(sessionid_B, "CircleOfB", true)

	sessionid_A := req.PostSessionGetSessionId("handleA", "password1")

	response, err := req.PostJoin(sessionid_A, "handleC", "CircleOfB")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Bad request, user handleC wasn't found")
	c.Check(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostJoinUserBlocked(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid_B := req.PostSessionGetSessionId("handleB", "password2")

	req.PostCircles(sessionid_B, "CircleOfHandleB", true)
	req.PostBlock(sessionid_B, "handleA")

	sessionid_A := req.PostSessionGetSessionId("handleA", "password1")

	response, err := req.PostJoin(sessionid_A, "handleB", "CircleOfHandleB")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Server refusal to comply with join request")
	c.Check(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestPostJoinCircleNoExist(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid_A := req.PostSessionGetSessionId("handleA", "password1")

	response, err := req.PostJoin(sessionid_A, "handleB", "NonExistentCircle")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Could not find target circle, join failed")
	c.Check(response.StatusCode, Equals, 404)
}

func (s *TestSuite) TestPostJoinCreated(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid_B := req.PostSessionGetSessionId("handleB", "password2")

	req.PostCircles(sessionid_B, "MyCircle", true)

	sessionid_A := req.PostSessionGetSessionId("handleA", "password1")

	response, err := req.PostJoin(sessionid_A, "handleB", "MyCircle")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Join request successful!")
	c.Check(response.StatusCode, Equals, 201)
}
