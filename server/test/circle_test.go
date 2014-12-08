package api_test

import (
	"../types"
	"./helper"
	. "gopkg.in/check.v1"
	"time"
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
	token := req.PostSessionGetAuthToken("handleA", "password1")

	if res, err := req.PostCircles(token, "MyPublicCircle", true); err != nil {
		c.Error(err)
	} else {
		data := types.CircleResponse{}
		helper.Unmarshal(res, &data)
		c.Check(res.StatusCode, Equals, 201)
		c.Check(data.Name, Equals, "MyPublicCircle")
		// [TODO] This can be improved, not best way to assure correct url
		c.Check(data.Url, Not(Equals), "")
		c.Check(data.Owner, Equals, "handleA")
		c.Check(data.Description, Equals, "")
		c.Check(data.Visibility, Equals, "public")
		// [TODO] This can be improved, not best way to assure correct members url
		c.Check(data.Members, Not(Equals), "")
		// [TODO] This can be improved, not best way to assure correct date
		c.Check(data.Created, Not(Equals), time.Time{})
	}
}

func (s *TestSuite) TestPostPrivateCircleOK(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")
	token := req.PostSessionGetAuthToken("handleA", "password1")

	if res, err := req.PostCircles(token, "MyPrivateCircle", false); err != nil {
		c.Error(err)
	} else {
		data := types.CircleResponse{}
		helper.Unmarshal(res, &data)
		c.Check(res.StatusCode, Equals, 201)
		c.Check(data.Name, Equals, "MyPrivateCircle")
		// [TODO] This can be improved, not best way to assure correct url
		c.Check(data.Url, Not(Equals), "")
		c.Check(data.Owner, Equals, "handleA")
		c.Check(data.Description, Equals, "")
		c.Check(data.Visibility, Equals, "private")
		// [TODO] This can be improved, not best way to assure correct members url
		c.Check(data.Members, Not(Equals), "")
		// [TODO] This can be improved, not best way to assure correct date
		c.Check(data.Created, Not(Equals), time.Time{})
	}
}

//
// Search Circles Tests:
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
		"token": token,
		"user":  "",
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
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")
	token_A := req.PostSessionGetAuthToken("handleA", "password1")
	token_B := req.PostSessionGetAuthToken("handleB", "password2")

	// One public and private circle by handleA
	// Test from perspective of both
	if _, err := req.PostCircles(token_A, "TheFirstCircle", true); err != nil {
		c.Error(err)
	}
	if _, err := req.PostCircles(token_A, "TheSecondCircle", false); err != nil {
		c.Error(err)
	}
	if res, err := req.GetCircles(types.Json{
		"token": token_A,
		"user":  "handleA",
	}); err != nil {
		c.Error(err)
	} else {
		body := types.SearchCirclesResponse{}
		helper.Unmarshal(res, &body)
		c.Check(res.StatusCode, Equals, 200)
		c.Check(len(body.Results), Equals, 4)
		c.Check(body.Count, Equals, 4)
	}
	// handleB will see the same json
	if res, err := req.GetCircles(types.Json{
		"token": token_B,
		"user":  "handleA",
	}); err != nil {
		c.Error(err)
	} else {
		body := types.SearchCirclesResponse{}
		helper.Unmarshal(res, &body)
		c.Check(res.StatusCode, Equals, 200)
		c.Check(len(body.Results), Equals, 4)
		c.Check(body.Count, Equals, 4)
	}
}

// Circles join by the user are from a variety of creators
func (s *TestSuite) TestSearchCirclesOfTargetVarietyOK(c *C) {
	req.PostSignup("Alpha", "testA@test.io", "password1", "password1")
	req.PostSignup("Bravo", "testB@test.io", "password2", "password2")
	req.PostSignup("Charlie", "testC@test.io", "password3", "password3")
	req.PostSignup("Delta", "testD@test.io", "password4", "password4")
	token_A := req.PostSessionGetAuthToken("Alpha", "password1")
	token_B := req.PostSessionGetAuthToken("Bravo", "password2")

	// Alpha will be the joiner. His joined circles will be tested from the perspective
	// of himself and Bravo
	if _, err := req.PostJoin(token_A, "Bravo", types.BROADCAST); err != nil {
		c.Error(err)
	}
	if _, err := req.PostJoin(token_A, "Charlie", types.BROADCAST); err != nil {
		c.Error(err)
	}
	if _, err := req.PostJoin(token_A, "Delta", types.BROADCAST); err != nil {
		c.Error(err)
	}

	if res, err := req.GetCircles(types.Json{
		"token": token_A,
	}); err != nil {
		c.Error(err)
	} else {
		body := types.SearchCirclesResponse{}
		helper.Unmarshal(res, &body)
		c.Check(res.StatusCode, Equals, 200)
		c.Check(len(body.Results), Equals, 5)
		c.Check(body.Count, Equals, 5)
	}

	if res, err := req.GetCircles(types.Json{
		"token": token_B,
		"user":  "Alpha",
	}); err != nil {
		c.Error(err)
	} else {
		body := types.SearchCirclesResponse{}
		helper.Unmarshal(res, &body)
		c.Check(res.StatusCode, Equals, 200)
		c.Check(len(body.Results), Equals, 5)
		c.Check(body.Count, Equals, 5)
	}
}

// This is essentially getting the circles one is part of
// without supplying an user parameter.
func (s *TestSuite) TestSearchCirclesNoSpecificUserOK(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")
	token_A := req.PostSessionGetAuthToken("handleA", "password1")
	token_B := req.PostSessionGetAuthToken("handleB", "password2")

	// One public and private circle by handleA
	// One circle by handleB
	// Test from unsupplied user parameter perspective of both
	if _, err := req.PostCircles(token_A, "CharmanderCircle", true); err != nil {
		c.Error(err)
	}
	if _, err := req.PostCircles(token_A, "CharmeleonCircle", false); err != nil {
		c.Error(err)
	}
	if _, err := req.PostCircles(token_B, "BulbasaurCircle", false); err != nil {
		c.Error(err)
	}
	// Test unsupplied user param and user="" inputs (they work the same)
	if res, err := req.GetCircles(types.Json{
		"token": token_A,
	}); err != nil {
		c.Error(err)
	} else {
		body := types.SearchCirclesResponse{}
		helper.Unmarshal(res, &body)
		c.Check(res.StatusCode, Equals, 200)
		c.Check(len(body.Results), Equals, 4)
		c.Check(body.Count, Equals, 4)
	}

	if res, err := req.GetCircles(types.Json{
		"token": token_A,
		"user":  "",
	}); err != nil {
		c.Error(err)
	} else {
		body := types.SearchCirclesResponse{}
		helper.Unmarshal(res, &body)
		c.Check(res.StatusCode, Equals, 200)
		c.Check(len(body.Results), Equals, 4)
		c.Check(body.Count, Equals, 4)
	}

	if res, err := req.GetCircles(types.Json{
		"token": token_B,
	}); err != nil {
		c.Error(err)
	} else {
		body := types.SearchCirclesResponse{}
		helper.Unmarshal(res, &body)
		c.Check(res.StatusCode, Equals, 200)
		c.Check(len(body.Results), Equals, 3)
		c.Check(body.Count, Equals, 3)
	}

	if res, err := req.GetCircles(types.Json{
		"token": token_B,
		"user":  "",
	}); err != nil {
		c.Error(err)
	} else {
		body := types.SearchCirclesResponse{}
		helper.Unmarshal(res, &body)
		c.Check(res.StatusCode, Equals, 200)
		c.Check(len(body.Results), Equals, 3)
		c.Check(body.Count, Equals, 3)
	}
}

func (s *TestSuite) TestSearchCirclesBeforeWorks(c *C) {
	// [TODO] stub, lower priority, but should be verified to work...
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

	req.PostCircles(sessionid_B, "CircleOfB", true)

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

func (s *TestSuite) TestPostJoinOK(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid_B := req.PostSessionGetAuthToken("handleB", "password2")

	req.PostCircles(sessionid_B, "CircleOfB", true)

	sessionid_A := req.PostSessionGetAuthToken("handleA", "password1")

	response, err := req.PostJoin(sessionid_A, "handleB", "CircleOfB")
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonResponseMessage(response), Equals, "Join request successful!")
	c.Check(response.StatusCode, Equals, 201)
}
