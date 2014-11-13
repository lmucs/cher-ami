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

	sessionid := req.PostSessionGetSessionId("handleA", "password1")

	req.PostMessages("Go is going gophers!", sessionid)
	req.PostMessages("Hypothesize about stuff", sessionid)
	req.PostMessages("The nearest exit may be behind you", sessionid)
	req.PostMessages("I make soap.", sessionid)

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

	sessionid := req.PostSessionGetSessionId("handleA", "password1")

	req.PostMessages("Go is going gophers!", sessionid)
	req.PostMessages("Hypothesize about stuff", sessionid)
	req.PostMessages("The nearest exit may be behind you", sessionid)
	req.PostMessages("I make soap.", sessionid)

	if res, _ := req.GetMessageById("some_id", sessionid); true {
		c.Check(res.StatusCode, Equals, 404)
		message_response := struct {
			Response string
			Object   string
		}{}
		helper.Unmarshal(res, &message_response)
		c.Check(message_response.Response, Equals, "No such message in any circle you can see")
	}

	if res, _ := req.GetMessageById("another-wrong-id", sessionid); true {
		c.Check(res.StatusCode, Equals, 404)
		message_response := struct {
			Response string
			Object   string
		}{}
		helper.Unmarshal(res, &message_response)
		c.Check(message_response.Response, Equals, "No such message in any circle you can see")
	}

	if res, _ := req.GetMessageById("2", sessionid); true {
		c.Check(res.StatusCode, Equals, 404)
		message_response := struct {
			Response string
			Object   string
		}{}
		helper.Unmarshal(res, &message_response)
		c.Check(message_response.Response, Equals, "No such message in any circle you can see")
	}
}

func (s *TestSuite) TestGetMessageByIdUserBlocked(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid_A := req.PostSessionGetSessionId("handleA", "password1")
	sessionid_B := req.PostSessionGetSessionId("handleB", "password2")

	req.PostBlock(sessionid_B, "handleA")
	message_id := req.PostMessageGetMessageId("Go is going gophers!", sessionid_B)

	// handleA attempts to retrieve
	if res, _ := req.GetMessageById(message_id, sessionid_A); true {
		c.Check(res.StatusCode, Equals, 404)
		message_response := struct {
			Response string
			Object   string
		}{}
		helper.Unmarshal(res, &message_response)
		c.Check(message_response.Response, Equals, "No such message in any circle you can see")
	}
}

func (s *TestSuite) TestGetMessageByIdPrivateCircle(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid_A := req.PostSessionGetSessionId("handleA", "password1")
	sessionid_B := req.PostSessionGetSessionId("handleB", "password2")

	message_id := req.PostMessageGetMessageId("Go is going gophers!", sessionid_B)
	req.PostCircles(sessionid_B, "SomePrivateCircle", false)

	if res, _ := req.GetMessageById(message_id, sessionid_A); true {
		c.Check(res.StatusCode, Equals, 404)
		message_response := struct {
			Response string
			Object   string
		}{}
		helper.Unmarshal(res, &message_response)
		c.Check(message_response.Response, Equals, "No such message in any circle you can see")
	}

}

// Successful retrieval by id
func (s *TestSuite) TestGetMessageByIdOK(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	req.PostSignup("handleB", "testB@test.io", "password2", "password2")

	sessionid_A := req.PostSessionGetSessionId("handleA", "password1")
	sessionid_B := req.PostSessionGetSessionId("handleB", "password2")

	messageid_1 := req.PostMessageGetMessageId("Go is going gophers!", sessionid_A)

	if res, _ := req.GetMessageById(messageid_1, sessionid_B); true {
		c.Check(res.StatusCode, Equals, 200)
		message_response := struct {
			Response string
			Object   string
		}{}
		helper.Unmarshal(res, &message_response)
		var msg MessageData
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
	sessionid := req.PostSessionGetSessionId("handleA", "password1")
	messageid := req.PostMessageGetMessageId("Hello, world!", sessionid)
	req.DeleteSessions(sessionid)

	patchObj := types.Json{
		"op":       "update",
		"resource": "content",
		"value":    "Hello, world! Again!",
	}
	// patchObj2 := types.Json{
	// 	"op":       "update",
	// 	"resource": "content",
	// 	"value":    "Hello, world! Again!",
	// }
	patch := []types.Json{patchObj}

	res, _ := req.EditMessage(patch, messageid, sessionid)
	c.Check(res.StatusCode, Equals, 401)
	c.Check(helper.GetJsonResponseMessage(res), Equals, "Failed to authenticate user request")
}

func (s *TestSuite) TestEditMessageMissingParams(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	sessionid := req.PostSessionGetSessionId("handleA", "password1")
	messageid := req.PostMessageGetMessageId("Hello, world!", sessionid)

	onlyOp := types.Json{
		"op": "update",
	}

	// onlyResource := types.Json{
	// 	"resource": "content",
	// }

	// onlyValue := types.Json{
	// 	"resource": "content",
	// }

	// onlyOpResource := types.Json{
	// 	"op":       "update",
	// 	"resource": "content",
	// }

	// onlyResourceValue := types.Json{
	// 	"resource": "content",
	// 	"value":    "Hello, world! Again!",
	// }

	// onlyOpValue := types.Json{
	// 	"op":    "update",
	// 	"value": "Hello, world! Again!",
	// }

	res, _ := req.EditMessage([]types.Json{onlyOp}, messageid, sessionid)
	c.Check(res.StatusCode, Equals, 400)
	c.Check(helper.GetJsonReasonMessage(res), Equals, "missing `resource` parameter in object 0")

}

func (s *TestSuite) TestEditMessageBadResource(c *C) {

}

func (s *TestSuite) TestEditMessageBadValue(c *C) {

}

func (s *TestSuite) TestEditMessageBadPatchObject(c *C) {

}

func (s *TestSuite) TestEditMessageUnableToPublish(c *C) {

}

func (s *TestSuite) TestEditMessageUnableToUnpublish(c *C) {

}

func (s *TestSuite) TestEditMessageContentOK(c *C) {

}

func (s *TestSuite) TestEditMessagePublishOK(c *C) {

}

func (s *TestSuite) TestEditMessageUnpublishOK(c *C) {

}

func (s *TestSuite) TestEditMessageAllOK(c *C) {

}
