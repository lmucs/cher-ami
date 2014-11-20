package api_test

import (
	"./helper"
	"encoding/json"
	. "gopkg.in/check.v1"
)

//
// Get User Tests:
//

// func (s *TestSuite) TestGetUserNotFound(c *C) {
// 	response, err := getUser("testing123")
// 	if err != nil {
// 		c.Error(err)
// 	}

// 	c.Check(helper.GetJsonResponseMessage(response), Equals, "No results found")
// 	c.Check(response.StatusCode, Equals, 404)
// }

// func (s *TestSuite) TestGetUserOK(c *C) {
// 	req.PostSignup("testing123", "testing123", "testing123", "testing123")

// 	response, err := getUser("testing123")
// 	if err != nil {
// 		c.Error(err)
// 	}

// 	handle := getJsonUserData(response)

// 	c.Check(handle, Equals, "testing123")
// 	c.Check(response.StatusCode, Equals, 200)
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
		c.Check(response.StatusCode, Equals, 200)
		c.Check(data.Count, Equals, 3)
		c.Check(data.Response, Equals, "Search complete")
		c.Check(data.Reason, Equals, "")
		c.Check(len(results), Equals, 3)
	}
}

//
// Delete User Tests:
//

func (s *TestSuite) TestDeleteUserInvalidPassword(c *C) {
	req.PostSignup("handleA", "test@test.io", "password1", "password1")

	sessionid := req.PostSessionGetAuthToken("handleA", "password1")

	response, err := req.DeleteUser("handleA", "notpassword1", sessionid)
	if err != nil {
		c.Error(err)
	}

	c.Check(helper.GetJsonReasonMessage(response), Equals, "Invalid username or password, please try again.")
	c.Check(response.StatusCode, Equals, 400)
}

// func (s *TestSuite) TestDeleteUserOK(c *C) {
// 	req.PostSignup("handleA", "test@test.io", "password1", "password1")

// 	sessionid := req.PostSessionGetAuthToken("handleA", "password1")

// 	deleteUserResponse, err := req.DeleteUser("handleA", "password1", sessionid)
// 	if err != nil {
// 		c.Error(err)
// 	}
// 	// TODO check if user really deleted
// 	// getUserResponse, err := getUser("handleA")
// 	// if err != nil {
// 	// 	c.Error(err)
// 	// }

// 	c.Check(deleteUserResponse.StatusCode, Equals, 204)
// }
