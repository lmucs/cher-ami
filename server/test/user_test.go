package api_test

import (
	"../types"
	"./helper"
	"encoding/json"
	. "gopkg.in/check.v1"
)

//
// Get User Tests:
//

// Unimplemented

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
// Edit User Tests:
//

func (s *TestSuite) TestEditUserAllOK(c *C) {
	req.PostSignup("handleA", "testA@test.io", "password1", "password1")
	token := req.PostSessionGetAuthToken("handleA", "password1")

	patchObj1 := types.Json{
		"resource": "firstname",
		"value":    "Thomas",
	}

	patchObj2 := types.Json{
		"resource": "lastname",
		"value":    "Anderson",
	}

	patchObj3 := types.Json{
		"resource": "gender",
		"value":    "Male",
	}

	patchObj4 := types.Json{
		"resource": "birthday",
		"value":    "2002-10-02",
	}

	patchObj5 := types.Json{
		"resource": "bio",
		"value":    "I entered the Matrix once.",
	}

	patchObj6 := types.Json{
		"resource": "interests",
		"value":    "Slowing down time, freeing the human race, blowing up elevators.",
	}

	patchObj7 := types.Json{
		"resource": "languages",
		"value":    "English",
	}

	patchObj8 := types.Json{
		"resource": "location",
		"value":    "San Francisco, CA",
	}

	res1, _ := req.EditUser(types.JsonArray{patchObj1, patchObj2, patchObj3, patchObj4, patchObj5, patchObj6, patchObj7, patchObj8}, "handleA", token)
	c.Check(res1.StatusCode, Equals, 200)
	c.Check(helper.GetJsonResponseMessage(res1), Equals, "Successfully updated user handleA")
	// TODO: need get users, then make this section more extensive

	patchObj9 := types.Json{
		"resource": "languages",
		"value":    "Pig Latin",
	}

	res2, _ := req.EditUser(types.JsonArray{patchObj9}, "handleA", token)
	c.Check(res2.StatusCode, Equals, 200)
	// TODO: check that updating a field doesn't delete the other fields.
}

//
// Delete User Tests:
//

// Unimplemented
