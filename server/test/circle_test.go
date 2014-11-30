package api_test

import (
	"../types"
	"./helper"
	. "gopkg.in/check.v1"
)

//
// Post Circles Tests:
//

func (s *TestSuite) TestPostCirclesUserNoExist(c *C) {
	if res, err := req.PostCircles("SomeSessionId", "SomeCircle", true); err != nil {
		c.Error(err)
	} else {
		c.Check(helper.GetJsonResponseMessage(res), Equals, "Failed to authenticate user request")
		c.Check(res.StatusCode, Equals, 401)
	}
}

func (s *TestSuite) TestPostCirclesUserNoSession(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")
	sessionid := req.PostSessionGetAuthToken("handleA", "password1")
	req.DeleteSessions(sessionid)

	if res, err := req.PostCircles(sessionid, "SomeCircle", true); err != nil {
		c.Error(err)
	} else {
		c.Check(helper.GetJsonResponseMessage(res), Equals, "Failed to authenticate user request")
		c.Check(res.StatusCode, Equals, 401)
	}
}

func (s *TestSuite) TestPostCirclesNameReservedGold(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")
	sessionid := req.PostSessionGetAuthToken("handleA", "password1")

	if res, err := req.PostCircles(sessionid, "Gold", false); err != nil {
		c.Error(err)
	} else {
		c.Check(helper.GetJsonReasonMessage(res), Equals, "Gold is a reserved circle name")
		c.Check(res.StatusCode, Equals, 403)
	}

	if res, err := req.PostCircles(sessionid, "Gold", true); err != nil {
		c.Error(err)
	} else {
		c.Check(helper.GetJsonReasonMessage(res), Equals, "Gold is a reserved circle name")
		c.Check(res.StatusCode, Equals, 403)
	}
}

func (s *TestSuite) TestPostCirclesNameReservedBroadcast(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")
	sessionid := req.PostSessionGetAuthToken("handleA", "password1")
	if res, err := req.PostCircles(sessionid, "Broadcast", false); err != nil {
		c.Error(err)
	} else {
		c.Check(helper.GetJsonReasonMessage(res), Equals, "Broadcast is a reserved circle name")
		c.Check(res.StatusCode, Equals, 403)
	}

	if res, err := req.PostCircles(sessionid, "Broadcast", true); err != nil {
		c.Error(err)
	} else {
		c.Check(helper.GetJsonReasonMessage(res), Equals, "Broadcast is a reserved circle name")
		c.Check(res.StatusCode, Equals, 403)
	}
}

func (s *TestSuite) TestPostCirclesNameEmpty(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")
	token := req.PostSessionGetAuthToken("handleA", "password1")

	if res, err := req.PostCircles(token, "", false); err != nil {
		c.Error(err)
	} else {
		c.Check(helper.GetJsonReasonMessage(res), Equals, "Missing `circlename` parameter")
		c.Check(res.StatusCode, Equals, 400)
	}

	if res, err := req.PostCircles(token, "", true); err != nil {
		c.Error(err)
	} else {
		c.Check(helper.GetJsonReasonMessage(res), Equals, "Missing `circlename` parameter")
		c.Check(res.StatusCode, Equals, 400)
	}
}

func (s *TestSuite) TestPostPublicCircleOK(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")
	sessionid := req.PostSessionGetAuthToken("handleA", "password1")

	if res, err := req.PostCircles(sessionid, "MyPublicCircle", true); err != nil {
		c.Error(err)
	} else {
		data := struct {
			Response string
			Chief    string
			Name     string
			Public   bool
			Id       string
		}{}
		c.Check(res.StatusCode, Equals, 201)
		helper.Unmarshal(res, &data)
		c.Check(data.Response, Equals, "Created new circle!")
		c.Check(data.Chief, Equals, "handleA")
		c.Check(data.Name, Equals, "MyPublicCircle")
		c.Check(data.Public, Equals, true)
	}
}

func (s *TestSuite) TestPostPrivateCircleOK(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")
	sessionid := req.PostSessionGetAuthToken("handleA", "password1")

	if res, err := req.PostCircles(sessionid, "MyPrivateCircle", false); err != nil {
		c.Error(err)
	} else {
		data := struct {
			Response string
			Chief    string
			Name     string
			Public   bool
			Id       string
		}{}
		helper.Unmarshal(res, &data)
		c.Check(res.StatusCode, Equals, 201)
		c.Check(data.Response, Equals, "Created new circle!")
		c.Check(data.Chief, Equals, "handleA")
		c.Check(data.Name, Equals, "MyPrivateCircle")
		c.Check(data.Public, Equals, false)
	}
}

//
// Search Cicles Tests:
//

func (s *TestSuite) TestSearchCirclesTargetNoExist(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")
	token := req.PostSessionGetAuthToken("handleA", "password1")

	if res, err := req.GetCircles(types.Json{
		"token": token,
		"user":  "handleB",
	}); err != nil {
		c.Error(err)
	} else {
		body := types.SearchCirclesResponse{}
		helper.Unmarshal(res, &body)
		c.Check(res.StatusCode, Equals, 200)
		c.Check(len(body.Results), Equals, 0)
		c.Check(body.Count, Equals, 0)
	}
}

func (s *TestSuite) TestSearchCirclesDefaultCirclesOK(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")
	token := req.PostSessionGetAuthToken("handleA", "password1")

	// Empty handle
	if res, err := req.GetCircles(types.Json{
		"token":  token,
		"handle": "",
	}); err != nil {
		c.Error(err)
	} else {
		body := types.SearchCirclesResponse{}
		helper.Unmarshal(res, &body)
		c.Check(res.StatusCode, Equals, 200)
		c.Check(len(body.Results), Equals, 2)
		c.Check(body.Count, Equals, 2)
	}

	// No handle
	if res, err := req.GetCircles(types.Json{
		"token": token,
	}); err != nil {
		c.Error(err)
	} else {
		body := types.SearchCirclesResponse{}
		helper.Unmarshal(res, &body)
		c.Check(res.StatusCode, Equals, 200)
		c.Check(len(body.Results), Equals, 2)
		c.Check(body.Count, Equals, 2)
	}
}

func (s *TestSuite) TestSearchCirclesOfTargetOK(c *C) {
	// req.PostSignup("handleA", "test@test.io", "password1", "password1")
	// token := req.PostSessionGetAuthToken("handleA", "password1")

	// test all properties that circles should have

}

func (s *TestSuite) TestSearchCirclesNoSpecificUserOK(c *C) {
	// stub
}

func (s *TestSuite) TestSearchCirclesBeforeWorks(c *C) {
	// stub
}

//
// Post Block Tests:
//

func (s *TestSuite) TestPostBlockUserNoExist(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	// This will be dangling, not belonging to handleA or any other user
	sessionid := req.PostSessionGetAuthToken("handleA", "password1")
	_ = req.PostSessionGetAuthToken("handleA", "password1")

	response, err := req.PostBlock(sessionid, "handleB")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Failed to authenticate user request")
	c.Check(response.StatusCode, Equals, 401)
}

func (s *TestSuite) TestPostBlockTargetNoExist(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")

	sessionid := req.PostSessionGetAuthToken("handleA", "password1")

	response, err := req.PostBlock(sessionid, "handleB")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonReasonMessage(response), Equals, "Bad request, user handleB wasn't found")
	c.Check(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostBlockUserNoSession(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid := req.PostSessionGetAuthToken("handleA", "password1")

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

	sessionid := req.PostSessionGetAuthToken("handleA", "password1")

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

	sessionid := req.PostSessionGetAuthToken("handleA", "password1")

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
	sessionid := helper.GetAuthTokenFromResponse(response)

	response, err := req.PostJoinDefault(sessionid, "handleB")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonReasonMessage(response), Equals, "Bad request, user handleB wasn't found")
	c.Check(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostJoinDefaultUserBlocked(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid_B := req.PostSessionGetAuthToken("handleB", "password2")

	req.PostBlock(sessionid_B, "handleA")

	sessionid_A := req.PostSessionGetAuthToken("handleA", "password1")

	response, err := req.PostJoinDefault(sessionid_A, "handleB")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonReasonMessage(response), Equals, "Server refusal to comply with join request")
	c.Check(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestPostJoinDefaultCreated(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid := req.PostSessionGetAuthToken("handleA", "password1")

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

	sessionid_B := req.PostSessionGetAuthToken("handleB", "password2")

	req.PostCircles(sessionid_B, "handleB", true)

	sessionid_A := req.PostSessionGetAuthToken("handleA", "password1")

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

	sessionid_B := req.PostSessionGetAuthToken("handleB", "password2")

	req.PostCircles(sessionid_B, "CircleOfB", true)

	sessionid_A := req.PostSessionGetAuthToken("handleA", "password1")

	response, err := req.PostJoin(sessionid_A, "handleC", "CircleOfB")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonReasonMessage(response), Equals, "Bad request, user handleC wasn't found")
	c.Check(response.StatusCode, Equals, 400)
}

func (s *TestSuite) TestPostJoinUserBlocked(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid_B := req.PostSessionGetAuthToken("handleB", "password2")

	req.PostCircles(sessionid_B, "CircleOfHandleB", true)
	req.PostBlock(sessionid_B, "handleA")

	sessionid_A := req.PostSessionGetAuthToken("handleA", "password1")

	response, err := req.PostJoin(sessionid_A, "handleB", "CircleOfHandleB")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonReasonMessage(response), Equals, "Server refusal to comply with join request")
	c.Check(response.StatusCode, Equals, 403)
}

func (s *TestSuite) TestPostJoinCircleNoExist(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid_A := req.PostSessionGetAuthToken("handleA", "password1")

	response, err := req.PostJoin(sessionid_A, "handleB", "NonExistentCircle")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonReasonMessage(response), Equals, "Could not find target circle, join failed")
	c.Check(response.StatusCode, Equals, 404)
}

func (s *TestSuite) TestPostJoinCreated(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid_B := req.PostSessionGetAuthToken("handleB", "password2")

	req.PostCircles(sessionid_B, "MyCircle", true)

	sessionid_A := req.PostSessionGetAuthToken("handleA", "password1")

	response, err := req.PostJoin(sessionid_A, "handleB", "MyCircle")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Join request successful!")
	c.Check(response.StatusCode, Equals, 201)
}
