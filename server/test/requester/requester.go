package requester

//
// This package is used for sending REST requests
// for api unit testing.
//

import (
	"../../types"
	helper "../helper/"
	"fmt"
	"net/http"
)

// Routes stored in a struct
type Routes struct {
	signupURL      string
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

//
// Constructor --- Use this!
//
func NewRequester(apiURL string) *Requester {
	routes := &Routes{
		fmt.Sprintf("%s/signup", apiURL),
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

func (req Requester) DeleteUser(handle string, password string, token string) (*http.Response, error) {
	payload := types.Json{
		"handle":   handle,
		"password": password,
		"token":    token,
	}
	deleteURL := req.Routes.usersURL + "/" + payload["handle"].(string)
	return helper.Execute("DELETE", deleteURL, payload)
}

func (req Requester) DeleteSessions(token string) (*http.Response, error) {
	payload := types.Json{
		"token": token,
	}

	return helper.Execute("DELETE", req.Routes.sessionsURL, payload)
}

func (req Requester) GetMessages(payload types.Json) (*http.Response, error) {
	return helper.GetWithQueryParams(req.Routes.messagesURL, payload)
}

func (req Requester) GetMessageById(id, token string) (*http.Response, error) {
	payload := types.Json{
		"token": token,
	}

	return helper.GetWithQueryParams(req.Routes.messagesURL+"/"+id, payload)
}

func (req Requester) GetMessageByUrl(url, token string) (*http.Response, error) {
	id := helper.GetIdFromUrlString(url)
	return req.GetMessageById(id, token)
}

func (req Requester) GetUser(handle string) (*http.Response, error) {
	payload := types.Json{
		"handle": handle,
	}

	return helper.GetWithQueryParams(req.Routes.usersURL+"/"+handle, payload)
}

func (req Requester) EditUser(payload types.JsonArray, handle, token string) (*http.Response, error) {
	return helper.ExecutePatch(token, req.Routes.usersURL+"/"+handle, payload)
}

func (req Requester) PostBlock(token string, target string) (*http.Response, error) {
	payload := types.Json{
		"token":  token,
		"target": target,
	}

	return helper.Execute("POST", req.Routes.blockURL, payload)
}

func (req Requester) PostCircles(token string, circleName string, public bool) (*http.Response, error) {
	payload := types.Json{
		"token":      token,
		"circlename": circleName,
		"public":     public,
	}

	return helper.Execute("POST", req.Routes.circlesURL, payload)
}

func (req Requester) PostCircleGetCircleId(token string, circleName string, public bool) string {
	payload := types.Json{
		"token":      token,
		"circlename": circleName,
		"public":     public,
	}
	res, err := helper.Execute("POST", req.Routes.circlesURL, payload)

	if err != nil {
		panic(err)
	}

	return helper.GetIdFromUrlField(res)
}

func (req Requester) GetCircles(payload types.Json) (*http.Response, error) {
	return helper.GetWithQueryParams(req.Routes.circlesURL, payload)
}

func (req Requester) PostJoin(token string, target string, circle string) (*http.Response, error) {
	payload := types.Json{
		"token":  token,
		"target": target,
		"circle": circle,
	}

	return helper.Execute("POST", req.Routes.joinURL, payload)
}

func (req Requester) PostMessage(content string, token string) (*http.Response, error) {
	payload := types.Json{
		"content": content,
		"token":   token,
	}

	return helper.Execute("POST", req.Routes.messagesURL, payload)
}

func (req Requester) PostMessageWithCircles(content string, token string, circles []string) (*http.Response, error) {
	payload := types.Json{
		"content": content,
		"token":   token,
		"circles": circles,
	}

	return helper.Execute("POST", req.Routes.messagesURL, payload)
}

func (req Requester) PostMessageGetMessageUrl(content, token string) string {
	payload := types.Json{
		"content": content,
		"token":   token,
	}
	res, err := helper.Execute("POST", req.Routes.messagesURL, payload)

	if err != nil {
		panic(err)
	}

	return helper.GetUrlFromResponse(res)
}

func (req Requester) PostMessageGetMessageId(content, token string) string {
	return helper.GetIdFromUrlString(req.PostMessageGetMessageUrl(content, token))
}

func (req Requester) PostMessageWithCirclesGetMessageUrl(content string, token string, circles []string) string {
	payload := types.Json{
		"content": content,
		"token":   token,
		"circles": circles,
	}
	if res, err := helper.Execute("POST", req.Routes.messagesURL, payload); err != nil {
		panic(err)
	} else {
		return helper.GetUrlFromResponse(res)
	}
}

func (req Requester) PostMessageWithCirclesGetMessageId(content, token string, circles []string) string {
	return helper.GetIdFromUrlString(req.PostMessageWithCirclesGetMessageUrl(content, token, circles))
}

func (req Requester) EditMessage(payload types.JsonArray, id string, token string) (*http.Response, error) {
	return helper.ExecutePatch(token, req.Routes.messagesURL+"/"+id, payload)
}

func (req Requester) PostSessions(handle string, password string) (*http.Response, error) {
	payload := types.Json{
		"handle":   handle,
		"password": password,
	}

	return helper.Execute("POST", req.Routes.sessionsURL, payload)
}

func (req Requester) PostSessionGetAuthToken(handle string, password string) (token string) {
	res, err := req.PostSessions(handle, password)
	if err != nil {
		panic("Unexpected failure to post session.")
	}
	return helper.GetAuthTokenFromResponse(res)
}

func (req Requester) PostSignup(handle string, email string, password string, confirmPassword string) (*http.Response, error) {
	proposal := types.Json{
		"handle":          handle,
		"email":           email,
		"password":        password,
		"confirmpassword": confirmPassword,
	}

	return helper.Execute("POST", req.Routes.signupURL, proposal)
}

func (req Requester) PostJoinDefault(token string, target string) (*http.Response, error) {
	payload := types.Json{
		"token":  token,
		"target": target,
	}

	return helper.Execute("POST", req.Routes.joindefaultURL, payload)
}

func (req Requester) SearchForUsers(circle, nameprefix string, skip, limit int, sort string) (*http.Response, error) {
	payload := types.Json{
		"circle":     circle,
		"nameprefix": nameprefix,
		"skip":       skip,
		"limit":      limit,
		"sort":       sort,
	}

	return helper.GetWithQueryParams(req.Routes.usersURL, payload)
}
