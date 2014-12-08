package api_test

import (
	"../types"
	"./helper"
	. "gopkg.in/check.v1"
)

//
// Post Message Tests
//

func (s *TestSuite) TestPostMessageInvalidAuth(c *C) {
	res, _ := req.PostMessage("Go is going gophers!", "SomeSessionId")
	c.Check(res.StatusCode, Equals, 401)
	c.Check(helper.GetJsonResponseMessage(res), Equals, "Failed to authenticate user request")
}

func (s *TestSuite) TestPostMessageEmptyContent(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	sessionid := req.PostSessionGetAuthToken("handleA", "password1")

	res, _ := req.PostMessage("", sessionid)

	c.Check(res.StatusCode, Equals, 400)
	c.Check(helper.GetJsonReasonMessage(res), Equals, "Please enter some content for your message")
}

func (s *TestSuite) TestPostMessageContentOnlyOK(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	sessionid := req.PostSessionGetAuthToken("handleA", "password1")

	if res, _ := req.PostMessage("Go is going gophers!", sessionid); true {
		c.Check(res.StatusCode, Equals, 201)
		m := types.MessageView{}
		helper.Unmarshal(res, &m)
		c.Check(m.Id, Not(Equals), "")
		c.Check(m.Url, Not(Equals), "")
		c.Check(m.Author, Equals, "handleA")
		c.Check(m.Content, Equals, "Go is going gophers!")
	}
}

func (s *TestSuite) TestPostMessageContentCirclesOK(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")

	sessionid := req.PostSessionGetAuthToken("handleA", "password1")

	circleid := req.PostCircleGetCircleId(sessionid, "MyPublicCircle", true)
	circleid2 := req.PostCircleGetCircleId(sessionid, "MyPublicCircle2", true)

	circles := []string{circleid, circleid2}

	if res, _ := req.PostMessageWithCircles("Go is going gophers!", sessionid, circles); true {
		c.Check(res.StatusCode, Equals, 201)
		m := types.MessageView{}
		helper.Unmarshal(res, &m)
		c.Check(m.Id, Not(Equals), "")
		c.Check(m.Url, Not(Equals), "")
		c.Check(m.Author, Equals, "handleA")
		c.Check(m.Content, Equals, "Go is going gophers!")
		// [TODO] ensure that the message was published successfully, we only know
		// it was successful because there was as 201 not a 400
	}
}

//
// Get Authored Messages Tests
//
func (s *TestSuite) TestGetMessagesInvalidAuth(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	res, _ := req.GetMessages(types.Json{
		"token": "",
	})
	c.Check(res.StatusCode, Equals, 401)
}

func (s *TestSuite) TestGetMessagesByTargetInCircle(c *C) {
	// stub
}

func (s *TestSuite) TestGetPublicMessagesByHandleOK(c *C) {
	// stub
}

func (s *TestSuite) TestGetMessagesInCircleOK(c *C) {
	req.PostSignup("Alpha", "testA@test.io", "password1", "password1")
	req.PostSignup("Bravo", "testB@test.io", "password2", "password2")
	req.PostSignup("Charlie", "testC@test.io", "password3", "password3")

	sessionid_A := req.PostSessionGetAuthToken("Alpha", "password1")
	sessionid_B := req.PostSessionGetAuthToken("Bravo", "password2")
	sessionid_C := req.PostSessionGetAuthToken("Charlie", "password3")

	circleid_alta := req.PostCircleGetCircleId(sessionid_B, "Alta", true)
	circleid_baja := req.PostCircleGetCircleId(sessionid_B, "Baja", true)

	// Alpha joins the circle
	if _, err := req.PostJoin(sessionid_A, "Bravo", "Baja"); err != nil {
		c.Error(err)
	}
	if _, err := req.PostJoin(sessionid_B, "Alpha", "Alta"); err != nil {
		c.Error(err)
	}
	if _, err := req.PostJoin(sessionid_C, "Alpha", "Alta"); err != nil {
		c.Error(err)
	}

	// No messages posted yet
	if res, err := req.GetMessages(types.Json{
		"token":  sessionid_A,
		"circle": circleid_alta,
	}); err != nil {
		c.Error(err)
	} else {
		c.Check(res.StatusCode, Equals, 200)
		mrv := []types.MessageResponseView{}
		helper.Unmarshal(res, &mrv)
		c.Check(len(mrv), Equals, 0)
	}
	if res, err := req.GetMessages(types.Json{
		"token":  sessionid_A,
		"circle": circleid_baja,
	}); err != nil {
		c.Error(err)
	} else {
		c.Check(res.StatusCode, Equals, 200)
		mrv := []types.MessageResponseView{}
		helper.Unmarshal(res, &mrv)
		c.Check(len(mrv), Equals, 0)
	}

	// Post some messages
	if res, _ := req.PostMessageWithCircles("Yo, sup", sessionid_A, []string{circleid_alta, circleid_baja}); true {
		c.Check(res.StatusCode, Equals, 201)
	}
	if res, _ := req.PostMessageWithCircles("Yo, dawg", sessionid_B, []string{circleid_alta, circleid_baja}); true {
		c.Check(res.StatusCode, Equals, 201)
	}
	if res, _ := req.PostMessageWithCircles("Yo, dude", sessionid_C, []string{circleid_alta}); true {
		c.Check(res.StatusCode, Equals, 201)
	}

	// Own circle
	if res, err := req.GetMessages(types.Json{
		"token":    sessionid_A,
		"circleid": circleid_alta,
	}); err != nil {
		c.Error(err)
	} else {
		c.Check(res.StatusCode, Equals, 200)
		mrv := []types.MessageResponseView{}
		helper.Unmarshal(res, &mrv)
		c.Check(len(mrv), Equals, 3)
	}

	// Other circle
	if res, err := req.GetMessages(types.Json{
		"token":    sessionid_A,
		"circleid": circleid_baja,
	}); err != nil {
		c.Error(err)
	} else {
		c.Check(res.StatusCode, Equals, 200)
		mrv := []types.MessageResponseView{}
		helper.Unmarshal(res, &mrv)
		c.Check(len(mrv), Equals, 2)
	}
}

func (s *TestSuite) TestGetMessageFeedOfSelfOK(c *C) {
	req.PostSignup("Alpha", "testA@test.io", "password1", "password1")
	req.PostSignup("Bravo", "testB@test.io", "password2", "password2")
	req.PostSignup("Charlie", "testC@test.io", "password3", "password3")

	sessionid_A := req.PostSessionGetAuthToken("Alpha", "password1")
	sessionid_B := req.PostSessionGetAuthToken("Bravo", "password2")
	sessionid_C := req.PostSessionGetAuthToken("Charlie", "password3")

	circleid_baja := req.PostCircleGetCircleId(sessionid_B, "Baja", true)
	circleid_cabo := req.PostCircleGetCircleId(sessionid_C, "Cabo", true)

	// Alpha joins the two circles
	if _, err := req.PostJoin(sessionid_A, "Bravo", "Baja"); err != nil {
		c.Error(err)
	}
	if _, err := req.PostJoin(sessionid_A, "Charlie", "Cabo"); err != nil {
		c.Error(err)
	}

	// No messages posted yet
	if res, err := req.GetMessages(types.Json{
		"token": sessionid_A,
	}); err != nil {
		c.Error(err)
	} else {
		c.Check(res.StatusCode, Equals, 200)
		mrv := []types.MessageResponseView{}
		helper.Unmarshal(res, &mrv)
		c.Check(len(mrv), Equals, 0)
	}

	// Post some messages
	if res, _ := req.PostMessageWithCircles("Yo, sup", sessionid_A, []string{circleid_baja, circleid_cabo}); true {
		c.Check(res.StatusCode, Equals, 201)
	}
	if res, _ := req.PostMessageWithCircles("Yo, dawg", sessionid_B, []string{circleid_baja}); true {
		c.Check(res.StatusCode, Equals, 201)
	}
	if res, _ := req.PostMessageWithCircles("Yo, dude", sessionid_C, []string{circleid_cabo}); true {
		c.Check(res.StatusCode, Equals, 201)
	}

	if res, err := req.GetMessages(types.Json{
		"token": sessionid_A,
	}); err != nil {
		c.Error(err)
	} else {
		c.Check(res.StatusCode, Equals, 200)
		mrv := []types.MessageResponseView{}
		helper.Unmarshal(res, &mrv)
		c.Check(len(mrv), Equals, 4)
	}

	// case when target is self (same response expected)
	if res, err := req.GetMessages(types.Json{
		"token": sessionid_A,
		"user":  "Alpha",
	}); err != nil {
		c.Error(err)
	} else {
		c.Check(res.StatusCode, Equals, 200)
		mrv := []types.MessageResponseView{}
		helper.Unmarshal(res, &mrv)
		c.Check(len(mrv), Equals, 4)
	}
}

//
// Get Message By ID Tests
//

// Absent or incorrect session token
func (s *TestSuite) TestGetMessageByIdInvalidAuth(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")

	if res, _ := req.GetMessageById("some_id", ""); true {
		c.Check(res.StatusCode, Equals, 401)
	}
	if res, _ := req.GetMessageById("some_id", "bad_id"); true {
		c.Check(res.StatusCode, Equals, 401)
	}
	if res, _ := req.GetMessageById("some_id", "sCxs2ad213124jP1241d"); true {
		c.Check(res.StatusCode, Equals, 401)
	}
}

// Target id doesn't exist
func (s *TestSuite) TestGetMessageByIdDoesNotExist(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")

	sessionid := req.PostSessionGetAuthToken("handleA", "password1")

	req.PostMessage("Go is going gophers!", sessionid)
	req.PostMessage("Hypothesize about stuff", sessionid)
	req.PostMessage("The nearest exit may be behind you", sessionid)
	req.PostMessage("I make soap.", sessionid)

	if res, _ := req.GetMessageById("some_id", sessionid); true {
		c.Check(res.StatusCode, Equals, 404)
		catcher := types.ReasonCatcher{}
		helper.Unmarshal(res, &catcher)
		c.Check(catcher.Reason, Equals, "No such message with id some_id could be found")
	}

	if res, _ := req.GetMessageById("another-wrong-id", sessionid); true {
		c.Check(res.StatusCode, Equals, 404)
		catcher := types.ReasonCatcher{}
		helper.Unmarshal(res, &catcher)
		c.Check(catcher.Reason, Equals, "No such message with id another-wrong-id could be found")
	}

	if res, _ := req.GetMessageById("2", sessionid); true {
		c.Check(res.StatusCode, Equals, 404)
		catcher := types.ReasonCatcher{}
		helper.Unmarshal(res, &catcher)
		c.Check(catcher.Reason, Equals, "No such message with id 2 could be found")
	}
}

func (s *TestSuite) TestGetMessageByIdUserBlocked(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid_A := req.PostSessionGetAuthToken("handleA", "password1")
	sessionid_B := req.PostSessionGetAuthToken("handleB", "password2")

	req.PostBlock(sessionid_B, "handleA")
	message_url := req.PostMessageGetMessageUrl("Go is going gophers!", sessionid_B)

	// handleA attempts to retrieve
	if res, _ := req.GetMessageByUrl(message_url, sessionid_A); true {
		c.Check(res.StatusCode, Equals, 404)
		catcher := types.ReasonCatcher{}
		helper.Unmarshal(res, &catcher)
		id := helper.GetIdFromUrlString(message_url)
		c.Check(catcher.Reason, Equals, "No such message with id "+id+" could be found")
	}
}

func (s *TestSuite) TestGetMessageByIdPrivateCircle(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid_A := req.PostSessionGetAuthToken("handleA", "password1")
	sessionid_B := req.PostSessionGetAuthToken("handleB", "password2")

	message_url := req.PostMessageGetMessageUrl("Go is going gophers!", sessionid_B)
	req.PostCircles(sessionid_B, "SomePrivateCircle", false)

	if res, _ := req.GetMessageByUrl(message_url, sessionid_A); true {
		c.Check(res.StatusCode, Equals, 404)
		catcher := types.ReasonCatcher{}
		helper.Unmarshal(res, &catcher)
		id := helper.GetIdFromUrlString(message_url)
		c.Check(catcher.Reason, Equals, "No such message with id "+id+" could be found")
	}
}

// Successful retrieval by id
func (s *TestSuite) TestGetMessageByIdOK(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")
	sessionid_A := req.PostSessionGetAuthToken("handleA", "password1")
	sessionid_B := req.PostSessionGetAuthToken("handleB", "password2")

	circleid_1 := req.PostCircleGetCircleId(sessionid_A, "MyPublicCircle", true)
	req.PostJoin(sessionid_B, "handleA", "MyPublicCircle")
	message_url := req.PostMessageWithCirclesGetMessageUrl("Go is going gophers!", sessionid_A, []string{circleid_1})

	if res, _ := req.GetMessageByUrl(message_url, sessionid_B); true {
		m := types.MessageView{}
		helper.Unmarshal(res, &m)

		c.Check(res.StatusCode, Equals, 200)
		c.Check(m.Author, Equals, "handleA")
		c.Check(m.Content, Equals, "Go is going gophers!")
	}
}

//
// Edit Message Tests
//

func (s *TestSuite) TestEditMessageInvalidAuth(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	sessionid := req.PostSessionGetAuthToken("handleA", "password1")
	messageid := req.PostMessageGetMessageId("Hello, world!", sessionid)
	req.DeleteSessions(sessionid)

	patchObj := types.Json{
		"op":       "update",
		"resource": "content",
		"value":    "Hello, world! Again!",
	}

	patch := []types.Json{patchObj}

	res, _ := req.EditMessage(patch, messageid, sessionid)
	c.Check(res.StatusCode, Equals, 401)
	c.Check(helper.GetJsonResponseMessage(res), Equals, "Failed to authenticate user request")
}

func (s *TestSuite) TestEditMessageMissingParams(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	sessionid := req.PostSessionGetAuthToken("handleA", "password1")
	messageid := req.PostMessageGetMessageId("Hello, world!", sessionid)

	// Cases of missing parameters
	onlyOp := types.Json{
		"op": "update",
	}

	onlyResource := types.Json{
		"resource": "content",
	}

	onlyValue := types.Json{
		"resource": "content",
	}

	onlyOpResource := types.Json{
		"op":       "update",
		"resource": "content",
	}

	onlyResourceValue := types.Json{
		"resource": "content",
		"value":    "Hello, world! Again!",
	}

	onlyOpValue := types.Json{
		"op":    "update",
		"value": "Hello, world! Again!",
	}

	res, _ := req.EditMessage(types.JsonArray{onlyOp}, messageid, sessionid)
	c.Check(res.StatusCode, Equals, 400)
	resErr, index := helper.GetJsonPatchValidationReasonMessage(res)
	c.Check(resErr[0], Equals, "field resource is invalid: Required field for message patch")
	c.Check(index, Equals, 0)

	res, _ = req.EditMessage(types.JsonArray{onlyResource}, messageid, sessionid)
	c.Check(res.StatusCode, Equals, 400)
	resErr, index = helper.GetJsonPatchValidationReasonMessage(res)
	c.Check(resErr[0], Equals, "field op is invalid: Required field for message patch")
	c.Check(index, Equals, 0)

	res, _ = req.EditMessage(types.JsonArray{onlyValue}, messageid, sessionid)
	c.Check(res.StatusCode, Equals, 400)
	resErr, index = helper.GetJsonPatchValidationReasonMessage(res)
	c.Check(resErr[0], Equals, "field op is invalid: Required field for message patch")
	c.Check(index, Equals, 0)

	res, _ = req.EditMessage(types.JsonArray{onlyOpResource}, messageid, sessionid)
	c.Check(res.StatusCode, Equals, 400)
	resErr, index = helper.GetJsonPatchValidationReasonMessage(res)
	c.Check(resErr[0], Equals, "field value is invalid: Required field for message patch")
	c.Check(index, Equals, 0)

	res, _ = req.EditMessage(types.JsonArray{onlyResourceValue}, messageid, sessionid)
	c.Check(res.StatusCode, Equals, 400)
	resErr, index = helper.GetJsonPatchValidationReasonMessage(res)
	c.Check(resErr[0], Equals, "field op is invalid: Required field for message patch")
	c.Check(index, Equals, 0)

	res, _ = req.EditMessage(types.JsonArray{onlyOpValue}, messageid, sessionid)
	c.Check(res.StatusCode, Equals, 400)

	resErr, index = helper.GetJsonPatchValidationReasonMessage(res)
	c.Check(resErr[0], Equals, "field resource is invalid: Required field for message patch")
	c.Check(index, Equals, 0)
}

func (s *TestSuite) TestEditMessageBadOp(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	sessionid := req.PostSessionGetAuthToken("handleA", "password1")
	messageid := req.PostMessageGetMessageId("Hello, world!", sessionid)

	patchObj := types.Json{
		"op":       "change", // guess at op name
		"resource": "content",
		"value":    "Hello, world! Again!",
	}

	res, _ := req.EditMessage(types.JsonArray{patchObj}, messageid, sessionid)
	resErr, index := helper.GetJsonPatchValidationReasonMessage(res)
	c.Check(res.StatusCode, Equals, 400)
	c.Check(resErr[0], Equals, "field op is invalid: change")
	c.Check(index, Equals, 0)
}

func (s *TestSuite) TestEditMessageBadResource(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	sessionid := req.PostSessionGetAuthToken("handleA", "password1")
	messageid := req.PostMessageGetMessageId("Hello, world!", sessionid)

	patchObj := types.Json{
		"op":       "update",
		"resource": "messageText", // guess at resource name
		"value":    "Hello, world! Again!",
	}

	res, _ := req.EditMessage(types.JsonArray{patchObj}, messageid, sessionid)
	resErr, index := helper.GetJsonPatchValidationReasonMessage(res)
	c.Check(res.StatusCode, Equals, 400)
	c.Check(resErr[0], Equals, "field resource is invalid: messageText")
	c.Check(index, Equals, 0)
}

func (s *TestSuite) TestEditMessageBadPatchObject(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	sessionid := req.PostSessionGetAuthToken("handleA", "password1")
	messageid := req.PostMessageGetMessageId("Hello, world!", sessionid)

	publishContent := types.Json{
		"op":       "publish",
		"resource": "content",
		"value":    "Hello, world! Again!",
	}
	res, _ := req.EditMessage(types.JsonArray{publishContent}, messageid, sessionid)
	c.Check(res.StatusCode, Equals, 400)
	c.Check(helper.GetJsonReasonMessage(res), Equals, "Malformed patch request at object 0")

	unpublishContent := types.Json{
		"op":       "unpublish",
		"resource": "content",
		"value":    "Hello, world! Again!",
	}
	res, _ = req.EditMessage(types.JsonArray{unpublishContent}, messageid, sessionid)
	c.Check(res.StatusCode, Equals, 400)
	c.Check(helper.GetJsonReasonMessage(res), Equals, "Malformed patch request at object 0")
}

func (s *TestSuite) TestEditMessageUnableToPublish(c *C) {
	// stub
}

func (s *TestSuite) TestEditMessageUnableToUnpublish(c *C) {
	// stub
}

func (s *TestSuite) TestEditMessageUpdateContentOK(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	sessionid := req.PostSessionGetAuthToken("handleA", "password1")
	messageid := req.PostMessageGetMessageId("Hello, world!", sessionid)

	patchObj := types.Json{
		"op":       "update",
		"resource": "content",
		"value":    "Hello, world! Again!",
	}

	res, _ := req.EditMessage(types.JsonArray{patchObj}, messageid, sessionid)
	c.Check(res.StatusCode, Equals, 200)
	c.Check(helper.GetJsonResponseMessage(res), Equals, "Successfully patched message "+messageid)
}

func (s *TestSuite) TestEditMessagePublishCircleOK(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	sessionid := req.PostSessionGetAuthToken("handleA", "password1")
	messageid := req.PostMessageGetMessageId("Hello, world!", sessionid)
	circleid := req.PostCircleGetCircleId(sessionid, "MyPublicCircle", true)

	patchObj := types.Json{
		"op":       "publish",
		"resource": "circle",
		"value":    circleid,
	}

	res, _ := req.EditMessage(types.JsonArray{patchObj}, messageid, sessionid)
	c.Check(res.StatusCode, Equals, 200)
	c.Check(helper.GetJsonResponseMessage(res), Equals, "Successfully patched message "+messageid)
}

func (s *TestSuite) TestEditMessageUnpublishCircleOK(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	sessionid := req.PostSessionGetAuthToken("handleA", "password1")
	circleid := req.PostCircleGetCircleId(sessionid, "MyPublicCircle", true)
	messageid := req.PostMessageWithCirclesGetMessageId("Hello, world!", sessionid, []string{circleid})

	patchObj := types.Json{
		"op":       "unpublish",
		"resource": "circle",
		"value":    circleid,
	}

	res, _ := req.EditMessage(types.JsonArray{patchObj}, messageid, sessionid)
	c.Check(res.StatusCode, Equals, 200)
	c.Check(helper.GetJsonResponseMessage(res), Equals, "Successfully patched message "+messageid)
}

func (s *TestSuite) TestEditMessageAllOK(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	sessionid := req.PostSessionGetAuthToken("handleA", "password1")
	circleid1 := req.PostCircleGetCircleId(sessionid, "MyPublicCircle", true)
	messageid := req.PostMessageWithCirclesGetMessageId("Hello, world!", sessionid, []string{circleid1})

	patchObj1 := types.Json{
		"op":       "update",
		"resource": "content",
		"value":    "Hello, world! Again!",
	}

	circleid2 := req.PostCircleGetCircleId(sessionid, "MyPublicCircle2", true)

	patchObj2 := types.Json{
		"op":       "publish",
		"resource": "circle",
		"value":    circleid2,
	}

	patchObj3 := types.Json{
		"op":       "unpublish",
		"resource": "circle",
		"value":    circleid1,
	}

	res, _ := req.EditMessage(types.JsonArray{patchObj1, patchObj2, patchObj3}, messageid, sessionid)
	c.Check(res.StatusCode, Equals, 200)
	c.Check(helper.GetJsonResponseMessage(res), Equals, "Successfully patched message "+messageid)
}
