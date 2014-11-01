package api

import (
	"../service"
	"./responses"
	encoding "encoding/json"
	"github.com/ChimeraCoder/go.crypto/bcrypt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jmcvetta/neoism"
	"net/http"
	"strconv"
	"time"
)

func panicErr(err error) {
	if err != nil {
		panic(err)
		return
	}
}

//
// Application Types
//

type Api struct {
	Svc  *service.Svc
	Resp *responses.Resp
}

/**
 * Constructor
 */
func NewApi(uri string) *Api {
	api := &Api{
		service.NewService(uri),
		&responses.Resp{},
	}
	return api
}

// Constants
const (
	GOLD            = "Gold"
	BROADCAST       = "Broadcast"
	MIN_PASS_LENGTH = 8
)

//
// API util
//

type json map[string]interface{}

func (a Api) authenticate(r *rest.Request) (success bool) {
	if sessionid := r.Header.Get("Authorization"); sessionid != "" {
		return a.Svc.VerifySession(sessionid)
	} else {
		return false
	}
}

func (a Api) getSessionId(r *rest.Request) string {
	return r.Header.Get("Authorization")
}

//
// API
//

//
// Credentials
//

/**
 * Expects a json POST with "username", "email", "password", "confirmpassword"
 */
func (a Api) Signup(w rest.ResponseWriter, r *rest.Request) {
	type Proposal struct {
		Handle          string
		Email           string
		Password        string
		ConfirmPassword string
	}
	proposal := Proposal{}
	if err := r.DecodeJsonPayload(&proposal); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handle := proposal.Handle
	email := proposal.Email
	password := proposal.Password
	confirm_password := proposal.ConfirmPassword

	// Handle and Email checks
	if handle == "" {
		a.Resp.SimpleJsonResponse(w, 400, "Handle is a required field for signup")
		return
	} else if email == "" {
		a.Resp.SimpleJsonResponse(w, 400, "Email is a required field for signup")
		return
	}

	// Password checks
	if password != confirm_password {
		a.Resp.SimpleJsonResponse(w, 400, "Passwords do not match")
		return
	} else if len(password) < MIN_PASS_LENGTH {
		a.Resp.SimpleJsonResponse(w, 400, "Passwords must be at least 8 characters long")
		return
	}

	// Ensure unique handle
	if unique, err := a.Svc.HandleIsUnique(handle); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !unique {
		a.Resp.SimpleJsonResponse(w, 400, "Sorry, handle or email is already taken")
		return
	}

	// Ensure unique email
	if unique, err := a.Svc.EmailIsUnique(email); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !unique {
		a.Resp.SimpleJsonResponse(w, 400, "Sorry, handle or email is already taken")
		return
	}

	var hashed_pass string
	if hash, err := bcrypt.GenerateFromPassword([]byte(password), 10); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		hashed_pass = string(hash)
	}
	if !a.Svc.CreateNewUser(handle, email, hashed_pass) {
		a.Resp.SimpleJsonResponse(w, 400, "Unexpected failure to create new user")
		return
	}

	if !a.Svc.MakeDefaultCirclesFor(handle) {
		a.Resp.SimpleJsonResponse(w, 400, "Unexpected failure to make default circles")
		return
	}

	w.WriteHeader(201)
	w.WriteJson(json{
		"Response": "Signed up a new user!",
		"Handle":   handle,
		"Email":    email,
	})
}

func (a Api) Login(w rest.ResponseWriter, r *rest.Request) {
	credentials := struct {
		Handle   string
		Password string
	}{}
	if err := r.DecodeJsonPayload(&credentials); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handle := credentials.Handle
	password := []byte(credentials.Password)

	if passwordHash, ok := a.Svc.GetPasswordHash(handle); !ok {
		a.Resp.SimpleJsonResponse(w, 403, "Invalid username or password, please try again.")
		return
	} else {
		// err is nil if successful, error if comparison failed
		if err := bcrypt.CompareHashAndPassword(passwordHash, password); err != nil {
			a.Resp.SimpleJsonResponse(w, 403, "Invalid username or password, please try again.")
			return
		} else {
			// Create an authentication node and return it to client
			sessionid := a.Svc.SetGetNewSessionId(handle)

			w.WriteHeader(200)
			w.WriteJson(json{
				"Response":  "Logged in " + handle + ". Note your session id.",
				"sessionid": sessionid,
			})
			return
		}
	}
}

/**
 * Expects a json post with "handle"
 */
func (a Api) Logout(w rest.ResponseWriter, r *rest.Request) {
	user := struct {
		Handle string
	}{}

	if err := r.DecodeJsonPayload(&user); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handle := user.Handle

	if a.Svc.UnsetSessionId(handle) {
		w.WriteHeader(204)
		return
	} else {

		a.Resp.SimpleJsonResponse(w, 403, "Cannot invalidate token because it is missing")
		return
	}
}

func (a Api) ChangePassword(w rest.ResponseWriter, r *rest.Request) {
	user := struct {
		Handle             string
		Password           string
		NewPassword        string
		ConfirmNewPassword string
	}{}

	if err := r.DecodeJsonPayload(&user); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handle := user.Handle
	password := []byte(user.Password)
	newPassword := user.NewPassword
	confirmNewPassword := user.ConfirmNewPassword

	if !a.authenticate(r) {
		a.Resp.FailedToAuthenticate(w)
		return
	}

	// Password checks
	if newPassword != confirmNewPassword {
		w.WriteHeader(400)
		w.WriteJson(json{
			"Response": "Passwords do not match",
		})
		return
	} else if len(newPassword) < MIN_PASS_LENGTH {
		a.Resp.SimpleJsonResponse(w, 400, "Passwords must be at least 8 characters long")
		return
	}

	if passwordHash, ok := a.Svc.GetPasswordHash(handle); !ok {
		a.Resp.SimpleJsonResponse(w, 400, "Invalid username or password, please try again.")
		return
	} else {
		// err is nil if successful, error
		if err := bcrypt.CompareHashAndPassword(passwordHash, password); err != nil {
			w.WriteHeader(400)
			w.WriteJson(json{
				"Response": "Invalid username or password, please try again.",
			})
			return
		} else if err := bcrypt.CompareHashAndPassword(passwordHash, []byte(newPassword)); err == nil {
			a.Resp.SimpleJsonResponse(w, 400, "Current/new password are same, please provide a new password.")
			return
		} else {
			if hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10); err != nil {
				rest.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				hashed_new_pass := string(hash)
				// Now set the new password
				if !a.Svc.SetNewPassword(handle, hashed_new_pass) {
					a.Resp.SimpleJsonResponse(w, 400, "Password change unsuccessful")
					return
				} else {
					// No JSON is written.
					w.WriteHeader(204)
					return
				}
			}
		}
	}
}

//
// User
//

func (a Api) SearchForUsers(w rest.ResponseWriter, r *rest.Request) {
	querymap := r.URL.Query()

	var circle string
	var nameprefix string
	var skip int
	var limit int
	var sort string

	if val, ok := querymap["limit"]; !ok {
		limit = 10
	} else {
		if intval, err := strconv.Atoi(val[0]); err != nil {
			w.WriteHeader(400)
			w.WriteJson(json{
				"Results":  nil,
				"Response": "Search failed",
				"Reason":   "Malformed limit",
				"Count":    0,
			})
			return

		} else {
			if intval > 100 || intval < 1 {
				w.WriteHeader(400)
				w.WriteJson(json{
					"Results":  nil,
					"Response": "Search failed",
					"Reason":   "Limit out of range",
					"Count":    0,
				})
			} else {
				limit = intval
			}
		}
	}

	if val, ok := querymap["nameprefix"]; !ok {
		nameprefix = ""
	} else {
		nameprefix = val[0]
	}

	if val, ok := querymap["circle"]; !ok {
		circle = ""
	} else {
		circle = val[0]
	}

	if val, ok := querymap["skip"]; !ok {
		skip = 0
	} else {
		if intval, err := strconv.Atoi(val[0]); err != nil {
			w.WriteHeader(400)
			w.WriteJson(json{
				"Results":  nil,
				"Response": "Search failed",
				"Reason":   "Malformed skip",
				"Count":    0,
			})
			return
		} else {
			skip = intval
		}
	}

	if sortType, ok := querymap["sort"]; !ok {
		w.WriteHeader(400)
		w.WriteJson(json{
			"Results":  nil,
			"Response": "Search failed",
			"Reason":   "Missing required sort parameter",
			"Count":    0,
		})
		return
	} else if sortType[0] != "handle" && sortType[0] != "joined" {
		w.WriteHeader(200)
		w.WriteJson(json{
			"Results":  nil,
			"Response": "Search failed",
			"Reason":   "No such sort " + sortType[0],
			"Count":    0,
		})
		return
	}

	results, count := a.Svc.SearchForUsers(circle, nameprefix, skip, limit, sort)

	w.WriteHeader(200)
	w.WriteJson(json{
		"Results":  results,
		"Response": "Search complete",
		"Count":    count,
	})
}

func (a Api) DeleteUser(w rest.ResponseWriter, r *rest.Request) {
	credentials := struct {
		Handle   string
		Password string
	}{}
	if err := r.DecodeJsonPayload(&credentials); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handle := credentials.Handle
	password := []byte(credentials.Password)

	if !a.authenticate(r) {
		a.Resp.FailedToAuthenticate(w)
		return
	}

	if passwordHash, ok := a.Svc.GetPasswordHash(handle); !ok {
		w.WriteHeader(400)
		w.WriteJson(json{
			"Response": "Invalid username or password, please try again.",
		})
		return
	} else {
		// err is nil if successful, error
		if err := bcrypt.CompareHashAndPassword(passwordHash, password); err != nil {
			w.WriteHeader(400)
			w.WriteJson(json{
				"Response": "Invalid username or password, please try again.",
			})
			return
		} else {
			if deleted := a.Svc.DeleteUser(handle); !deleted {
				w.WriteHeader(400)
				w.WriteJson(json{
					"Response": "Unexpected failure to delete user",
				})
				return
			} else {
				w.WriteHeader(204)
			}
		}
	}
}

//
// Circles
//

func (a Api) NewCircle(w rest.ResponseWriter, r *rest.Request) {
	payload := struct {
		Handle     string
		CircleName string
		Public     bool
	}{}
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handle := payload.Handle
	circleName := payload.CircleName
	isPublic := payload.Public

	if !a.authenticate(r) {
		a.Resp.FailedToAuthenticate(w)
		return
	}

	if circleName == GOLD || circleName == BROADCAST {
		a.Resp.SimpleJsonResponse(w, 403, circleName+" is a reserved circle name")
		return
	}

	if !a.Svc.NewCircle(handle, circleName, isPublic) {
		a.Resp.SimpleJsonResponse(w, 400, "Unexpected failure to create circle")
		return
	}

	a.Resp.SimpleJsonResponse(w, 201, "Created new circle "+circleName+" for "+handle)
}

//
// Messages
//

/**
 * Create a new, unpublished message
 */
func (a Api) NewMessage(w rest.ResponseWriter, r *rest.Request) {
	payload := struct {
		Handle  string
		Content string
	}{}
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handle := ""
	if payload.Handle == "" {
		if h, ok := a.Svc.GetHandleFromAuthorization(a.getSessionId(r)); !ok {
			w.WriteHeader(400)
			w.WriteJson(json{
				"Response": "Unexpected failure to retrieve owner of session",
			})
			return
		} else {
			handle = h
		}
	} else {
		handle = payload.Handle
	}
	content := payload.Content

	if !a.authenticate(r) {
		a.Resp.FailedToAuthenticate(w)
		return
	}

	if payload.Content == "" {
		a.Resp.SimpleJsonResponse(w, 400, "Please enter some content for your message")
		return
	}

	if !a.Svc.NewMessage(handle, content) {
		a.Resp.SimpleJsonResponse(w, 400, "No message created")
	} else {
		w.WriteHeader(201)
		w.WriteJson(json{
			"Response":  "Successfully created message for " + handle,
			"Published": false,
		})
	}
}

/**
 * Publishes a message identified by it's lastSaved time to a specific circle owned
 * by the user.
 */
func (a Api) PublishMessage(w rest.ResponseWriter, r *rest.Request) {
	payload := struct {
		Handle    string
		CircleId  string
		MessageId string
	}{}
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handle := payload.Handle
	circleid := payload.CircleId
	messageid := payload.MessageId

	if !a.authenticate(r) {
		a.Resp.FailedToAuthenticate(w)
		return
	}

	if !a.Svc.UserIsMemberOf(handle, circleid) {
		w.WriteHeader(401)
		w.WriteJson(json{
			"Response": "Refusal to comply with request",
			"Reason":   "You are not a member or owner of the specified circle",
		})
		return
	}

	if !a.Svc.CanSeeCircle(handle, circleid) {
		a.Resp.SimpleJsonResponse(w, 400, "Could not find specified circle to publish to")
		return
	} else if !a.Svc.MessageExists(messageid) {
		a.Resp.SimpleJsonResponse(w, 400, "Could not find intended message for publishing")
		return
	}

	if !a.Svc.PublishMessage(messageid, circleid) {
		a.Resp.SimpleJsonResponse(w, 400, "Bad request, no message published")
	} else {
		a.Resp.SimpleJsonResponse(w, 201, "Success! Published message to "+circleid)
	}
}

/**
 * Get messages authored by user
 */
func (a Api) GetAuthoredMessages(w rest.ResponseWriter, r *rest.Request) {
	if !a.authenticate(r) {
		a.Resp.FailedToAuthenticate(w)
		return
	}

	if author, success := a.Svc.GetHandleFromAuthorization(a.getSessionId(r)); !success {
		w.WriteHeader(400)
		w.WriteJson(json{
			"Response":  "Unexpected failure to retrieve owner of session",
			"Author":    author,
			"Success":   success,
			"SessionId": a.getSessionId(r),
		})
		return
	} else {
		type MessageData struct {
			Url     string
			Author  string
			Content string
			Date    time.Time
		}
		messages := a.Svc.GetMessagesByHandle(author)
		messageData := make([]MessageData, len(messages))

		for i := 0; i < len(messages); i++ {
			messageData[i] = MessageData{
				"<url>:<port>/api/messages/" + messages[i].Id, // hard-coded url/port...
				messages[i].Author,
				messages[i].Content,
				messages[i].Created,
			}
		}

		b, err := encoding.Marshal(messageData)
		if err != nil {
			panicErr(err)
		}

		w.WriteHeader(200)
		w.WriteJson(json{
			"Response": "Found messages for user " + author,
			"Objects":  string(b),
			"Count":    len(messageData),
		})
	}
}

/**
 * Get messages authored by a User that are visible to the authenticated
 * user. This means from all shared circles that the queried User has published to.
 */
func (a Api) GetMessagesByHandle(w rest.ResponseWriter, r *rest.Request) {
	author := r.PathParam("author")
	querymap := r.URL.Query()

	// check query parameters
	if _, ok := querymap["handle"]; !ok {
		w.WriteHeader(400)
		w.WriteJson(json{
			"Response": "Bad Request, not enough parameters to authenticate user",
		})
		return
	}

	handle := querymap["handle"][0]

	if !a.authenticate(r) {
		a.Resp.FailedToAuthenticate(w)
		return
	}

	if !a.Svc.UserExists(author) {
		w.WriteHeader(400)
		w.WriteJson(json{
			"Response": "Bad request, user doesn't exist",
		})
		return
	}

	messages := []struct {
		Content   string    `json:"message.content"`
		Published time.Time `json:"message.published"`
	}{}
	err := a.Svc.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (author:User {handle: {author}}), (user:User {handle: {handle}})
            OPTIONAL MATCH (user)-[r:MEMBER_OF]->(circle:Circle)
            OPTIONAL MATCH (author)-[w:WROTE]-(visible:Message)-[p:PUB_TO]->(circle)
            RETURN visible.content, visible.published_at
        `,
		Parameters: neoism.Props{
			"author": author,
			"handle": handle,
		},
		Result: &messages,
	})
	panicErr(err)

	w.WriteHeader(200)
	w.WriteJson(messages)
}

/**
 * Deletes an unpublished message
 */
func (a Api) DeleteMessage(w rest.ResponseWriter, r *rest.Request) {
	payload := struct {
		Handle    string
		LastSaved time.Time
	}{}
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handle := payload.Handle
	lastsaved := payload.LastSaved

	if !a.authenticate(r) {
		a.Resp.FailedToAuthenticate(w)
		return
	}

	deleted := []struct {
		Count int `json:"count(m)"`
	}{}
	if err := a.Svc.Db.Cypher(&neoism.CypherQuery{
		Statement: `
        MATCH (user:User {handle: {handle}})
        OPTIONAL MATCH (user)-[r:WROTE]->(m:Message {lastsaved: {lastsaved}})
        DELETE r, m
        RETURN count(m)
        `,
		Parameters: neoism.Props{
			"handle":    handle,
			"lastsaved": lastsaved,
		},
		Result: &deleted,
	}); err != nil {
		panicErr(err)
	}

	w.WriteHeader(200)
	w.WriteJson(json{
		"Response": "Success!",
		"Deleted":  deleted[0].Count,
	})
}

//
// Social
//

func (a Api) BlockUser(w rest.ResponseWriter, r *rest.Request) {
	payload := struct {
		Handle string
		Target string
	}{}
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handle := payload.Handle
	target := payload.Target

	if !a.authenticate(r) {
		a.Resp.FailedToAuthenticate(w)
		return
	}

	if !a.Svc.UserExists(target) {
		a.Resp.SimpleJsonResponse(w, 400, "Bad request, user "+target+" wasn't found")
		return
	}

	a.Svc.RevokeMembershipBetween(handle, target)

	if !a.Svc.CreateBlockFromTo(handle, target) {
		a.Resp.SimpleJsonResponse(w, 400, "Unexpected failure to block user")
	} else {
		a.Resp.SimpleJsonResponse(w, 200, "User "+target+" has been blocked")
	}
}

func (a Api) JoinDefault(w rest.ResponseWriter, r *rest.Request) {
	payload := struct {
		Handle string
		Target string
	}{}
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handle := payload.Handle
	target := payload.Target

	if !a.authenticate(r) {
		a.Resp.FailedToAuthenticate(w)
		return
	}

	if !a.Svc.UserExists(target) {
		a.Resp.SimpleJsonResponse(w, 400, "Bad request, user "+target+" wasn't found")
		return
	}

	if a.Svc.BlockExistsFromTo(target, handle) {
		a.Resp.SimpleJsonResponse(w, 403, "Server refusal to comply with join request")
		return
	}

	if !a.Svc.JoinBroadcast(handle, target) {
		a.Resp.SimpleJsonResponse(w, 400, "Unexpected failure to join Broadcast")
	} else {
		a.Resp.SimpleJsonResponse(w, 201, "JoinDefault request successful!")
	}
}

/**
 * Allows joining by (target, circlename) or (circleid) candidate keys
 */
func (a Api) Join(w rest.ResponseWriter, r *rest.Request) {
	payload := struct {
		Handle   string
		Target   string
		Circle   string
		CircleId string
	}{}
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handle := payload.Handle
	target := payload.Target
	circle := payload.Circle
	circleid := payload.CircleId

	if !a.authenticate(r) {
		a.Resp.FailedToAuthenticate(w)
		return
	}

	if !a.Svc.UserExists(target) {
		a.Resp.SimpleJsonResponse(w, 400, "Bad request, user "+target+" wasn't found")
		return
	}

	if a.Svc.BlockExistsFromTo(target, handle) {
		a.Resp.SimpleJsonResponse(w, 403, "Server refusal to comply with join request")
		return
	}

	if circleid == "" {
		if id := a.Svc.GetCircleId(target, circle); id == "" {
			a.Resp.SimpleJsonResponse(w, 404, "Could not find target circle, join failed")
			return
		} else {
			circleid = id
		}
	}

	if !a.Svc.CanSeeCircle(handle, circleid) {
		a.Resp.SimpleJsonResponse(w, 404, "Could not find target circle, join failed")
		return
	}

	if a.Svc.JoinCircle(handle, circleid) {
		a.Resp.SimpleJsonResponse(w, 201, "Join request successful!")
	} else {
		a.Resp.SimpleJsonResponse(w, 400, "Unexpected failure to join circle, join failed")
	}
}
