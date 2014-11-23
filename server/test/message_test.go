package api_test

import (
	"../types"
	"./helper"
	encoding "encoding/json"
	. "gopkg.in/check.v1"
	"time"
)

// Testing Structs
type MessageData struct {
	Id      string
	Url     string
	Author  string
	Content string
	Created time.Time
}

type MessageResponse struct {
	Response string
	Reason   string
	Object   string
}

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

	res, _ := req.PostMessage("Go is going gophers!", sessionid)

	c.Check(res.StatusCode, Equals, 201)
	c.Check(helper.GetJsonResponseMessage(res), Equals, "Successfully created message for handleA")
}

func (s *TestSuite) TestPostMessageContentCirclesOK(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")

	sessionid := req.PostSessionGetAuthToken("handleA", "password1")

	circleid := req.PostCircleGetCircleId(sessionid, "MyPublicCircle", true)
	circleid2 := req.PostCircleGetCircleId(sessionid, "MyPublicCircle2", true)

	circles := []string{circleid, circleid2}

	res, _ := req.PostMessageWithCircles("Go is going gophers!", sessionid, circles)

	c.Check(res.StatusCode, Equals, 201)
	c.Check(helper.GetJsonResponseMessage(res), Equals, "Successfully created message for handleA")
}

//
// Get Authored Messages Tests
//
func (s *TestSuite) TestGetAuthoredMessagesInvalidAuth(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	res, _ := req.GetAuthoredMessages("")
	c.Check(res.StatusCode, Equals, 401)
}

func (s *TestSuite) TestGetAuthoredMessagesOK(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")

	sessionid := req.PostSessionGetAuthToken("handleA", "password1")

	req.PostMessage("Go is going gophers!", sessionid)
	req.PostMessage("Hypothesize about stuff", sessionid)
	req.PostMessage("The nearest exit may be behind you", sessionid)
	req.PostMessage("I make soap.", sessionid)

	res, _ := req.GetAuthoredMessages(sessionid)

	data := struct {
		Response string
		Objects  string
		Count    int
	}{}
	helper.Unmarshal(res, &data)

	objects := make([]MessageData, 0)
	encoding.Unmarshal([]byte(data.Objects), &objects)

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
		message_response := MessageResponse{}
		helper.Unmarshal(res, &message_response)
		c.Check(message_response.Reason, Equals, "No such message with id some_id could be found")
	}

	if res, _ := req.GetMessageById("another-wrong-id", sessionid); true {
		c.Check(res.StatusCode, Equals, 404)
		message_response := MessageResponse{}
		helper.Unmarshal(res, &message_response)
		c.Check(message_response.Reason, Equals, "No such message with id another-wrong-id could be found")
	}

	if res, _ := req.GetMessageById("2", sessionid); true {
		c.Check(res.StatusCode, Equals, 404)
		message_response := MessageResponse{}
		helper.Unmarshal(res, &message_response)
		c.Check(message_response.Reason, Equals, "No such message with id 2 could be found")
	}
}

func (s *TestSuite) TestGetMessageByIdUserBlocked(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid_A := req.PostSessionGetAuthToken("handleA", "password1")
	sessionid_B := req.PostSessionGetAuthToken("handleB", "password2")

	req.PostBlock(sessionid_B, "handleA")
	message_id := req.PostMessageGetMessageId("Go is going gophers!", sessionid_B)

	// handleA attempts to retrieve
	if res, _ := req.GetMessageById(message_id, sessionid_A); true {
		c.Check(res.StatusCode, Equals, 404)
		message_response := MessageResponse{}
		helper.Unmarshal(res, &message_response)
		c.Check(message_response.Reason, Equals, "No such message with id "+message_id+" could be found")
	}
}

func (s *TestSuite) TestGetMessageByIdPrivateCircle(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid_A := req.PostSessionGetAuthToken("handleA", "password1")
	sessionid_B := req.PostSessionGetAuthToken("handleB", "password2")

	message_id := req.PostMessageGetMessageId("Go is going gophers!", sessionid_B)
	req.PostCircles(sessionid_B, "SomePrivateCircle", false)

	if res, _ := req.GetMessageById(message_id, sessionid_A); true {
		c.Check(res.StatusCode, Equals, 404)
		message_response := MessageResponse{}
		helper.Unmarshal(res, &message_response)
		c.Check(message_response.Reason, Equals, "No such message with id "+message_id+" could be found")
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
	messageid_1 := req.PostMessageWithCirclesGetMessageId("Go is going gophers!", sessionid_A, []string{circleid_1})

	if res, _ := req.GetMessageById(messageid_1, sessionid_B); true {
		c.Check(res.StatusCode, Equals, 200)
		var (
			message_response MessageResponse
			msg              MessageData
		)
		helper.Unmarshal(res, &message_response)
		encoding.Unmarshal([]byte(message_response.Object), &msg)
		c.Check(message_response.Response, Equals, "Found message!")
		c.Check(msg.Id, Equals, messageid_1)
		c.Check(msg.Author, Equals, "handleA")
		c.Check(msg.Content, Equals, "Go is going gophers!")
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

}

func (s *TestSuite) TestEditMessageUnableToUnpublish(c *C) {

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
