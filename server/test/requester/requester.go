package requester

//
// This package is used for sending REST requests
// for api unit testing.
//

import (
	helper "../helper/"
	"fmt"
	"net/http"
)

// Routes stored in a struct
type Routes struct {
	signupURL      string
	changePassURL  string
	sessionsURL    string
	userURL        string
	usersURL       string
	messagesURL    string
	publishURL     string
	joindefaultURL string
	joinURL        string
	blockURL       string
	circlesURL     string
}

type Requester struct {
	Routes *Routes
}

type json map[string]interface{}

//
// Constructor --- Use this!
//
func NewRequester(apiURL string) *Requester {
	routes := &Routes{
		fmt.Sprintf("%s/signup", apiURL),
		fmt.Sprintf("%s/changepassword", apiURL),
		fmt.Sprintf("%s/sessions", apiURL),
		fmt.Sprintf("%s/users/user", apiURL),
		fmt.Sprintf("%s/users", apiURL),
		fmt.Sprintf("%s/messages", apiURL),
		fmt.Sprintf("%s/publish", apiURL),
		fmt.Sprintf("%s/joindefault", apiURL),
		fmt.Sprintf("%s/join", apiURL),
		fmt.Sprintf("%s/block", apiURL),
		fmt.Sprintf("%s/circles", apiURL),
	}
	req := &Requester{
		routes,
	}
	return req
}

//
// Send/Receive Calls to API:
//

func (req Requester) DeleteUser(handle string, password string, sessionid string) (*http.Response, error) {
	payload := json{
		"handle":    handle,
		"password":  password,
		"sessionid": sessionid,
	}
	deleteURL := req.Routes.usersURL + "/" + payload["handle"].(string)
	return helper.Execute("DELETE", deleteURL, payload)
}

func (req Requester) DeleteSessions(handle string) (*http.Response, error) {
	payload := json{
		"handle": handle,
	}

	return helper.Execute("DELETE", req.Routes.sessionsURL, payload)
}

func (req Requester) GetAuthoredMessages(handle string, sessionid string) (*http.Response, error) {
	payload := json{
		"handle":    handle,
		"sessionid": sessionid,
	}

	return helper.GetWithQueryParams(req.Routes.messagesURL, payload)
}

func (req Requester) GetMessageById(id, sessionid string) (*http.Response, error) {
	payload := json{
		"sessionid": sessionid,
	}

	return helper.GetWithQueryParams(req.Routes.messagesURL+"/"+id, payload)
}

func (req Requester) GetUser(handle string) (*http.Response, error) {
	payload := json{
		"handle": handle,
	}

	return helper.GetWithQueryParams(req.Routes.usersURL+"/"+handle, payload)
}

func (req Requester) PostBlock(handle string, sessionid string, target string) (*http.Response, error) {
	payload := json{
		"handle":    handle,
		"sessionid": sessionid,
		"target":    target,
	}

	return helper.Execute("POST", req.Routes.blockURL, payload)
}

func (req Requester) PostChangePassword(handle string, sessionid string, password string, newPassword string, confirmNewPassword string) (*http.Response, error) {
	payload := json{
		"handle":             handle,
		"sessionid":          sessionid,
		"password":           password,
		"newpassword":        newPassword,
		"confirmnewpassword": confirmNewPassword,
	}

	return helper.Execute("POST", req.Routes.changePassURL, payload)
}

func (req Requester) PostCircles(handle string, sessionid string, circleName string, public bool) (*http.Response, error) {
	payload := json{
		"handle":     handle,
		"sessionid":  sessionid,
		"circlename": circleName,
		"public":     public,
	}

	return helper.Execute("POST", req.Routes.circlesURL, payload)
}

func (req Requester) PostJoin(handle string, sessionid string, target string, circle string) (*http.Response, error) {
	payload := json{
		"handle":    handle,
		"sessionid": sessionid,
		"target":    target,
		"circle":    circle,
	}

	return helper.Execute("POST", req.Routes.joinURL, payload)
}

func (req Requester) PostMessages(content string, sessionid string) (*http.Response, error) {
	payload := json{
		"content":   content,
		"sessionid": sessionid,
	}

	return helper.Execute("POST", req.Routes.messagesURL, payload)
}

func (req Requester) PostMessageGetMessageId(content, sessionid string) string {
	payload := json{
		"content":   content,
		"sessionid": sessionid,
	}
	res, err := helper.Execute("POST", req.Routes.messagesURL, payload)

	if err != nil {
		panic(err)
	}

	return helper.GetIdFromResponse(res)
}

func (req Requester) PostSessions(handle string, password string) (*http.Response, error) {
	payload := json{
		"handle":   handle,
		"password": password,
	}

	return helper.Execute("POST", req.Routes.sessionsURL, payload)
}

func (req Requester) PostSessionGetSessionId(handle string, password string) (sessionid string) {
	res, err := req.PostSessions(handle, password)
	if err != nil {
		panic("Unexpected failure to post session.")
	}
	return helper.GetSessionFromResponse(res)
}

func (req Requester) PostSignup(handle string, email string, password string, confirmPassword string) (*http.Response, error) {
	proposal := json{
		"handle":          handle,
		"email":           email,
		"password":        password,
		"confirmpassword": confirmPassword,
	}

	return helper.Execute("POST", req.Routes.signupURL, proposal)
}

func (req Requester) PostJoinDefault(handle string, sessionid string, target string) (*http.Response, error) {
	payload := json{
		"handle":    handle,
		"sessionid": sessionid,
		"target":    target,
	}

	return helper.Execute("POST", req.Routes.joindefaultURL, payload)
}

func (req Requester) SearchForUsers(circle, nameprefix string, skip, limit int, sort string) (*http.Response, error) {
	payload := json{
		"circle":     circle,
		"nameprefix": nameprefix,
		"skip":       skip,
		"limit":      limit,
		"sort":       sort,
	}

	return helper.GetWithQueryParams(req.Routes.usersURL, payload)
}
