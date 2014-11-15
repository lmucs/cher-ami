package api

import (
	"../types"
	"./service"
	apiutil "./util"
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
	Util *apiutil.Util
}

/**
 * Constructor
 */
func NewApi(uri string) *Api {
	api := &Api{
		service.NewService(uri),
		&apiutil.Util{},
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
		a.Util.SimpleJsonResponse(w, 400, "Handle is a required field for signup")
		return
	} else if email == "" {
		a.Util.SimpleJsonResponse(w, 400, "Email is a required field for signup")
		return
	}

	// Password checks
	if password != confirm_password {
		a.Util.SimpleJsonResponse(w, 403, "Passwords do not match")
		return
	} else if len(password) < MIN_PASS_LENGTH {
		a.Util.SimpleJsonResponse(w, 403, "Passwords must be at least 8 characters long")
		return
	}

	// Ensure unique handle
	if unique := a.Svc.HandleIsUnique(handle); !unique {
		a.Util.SimpleJsonResponse(w, 409, "Sorry, handle or email is already taken")
		return
	}

	// Ensure unique email
	if unique := a.Svc.EmailIsUnique(email); !unique {
		a.Util.SimpleJsonResponse(w, 409, "Sorry, handle or email is already taken")
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
		a.Util.SimpleJsonResponse(w, http.StatusInternalServerError, "Unexpected failure to create new user")
		return
	}

	if !a.Svc.MakeDefaultCirclesFor(handle) {
		a.Util.SimpleJsonResponse(w, http.StatusInternalServerError, "Unexpected failure to make default circles")
		return
	}

	w.WriteHeader(201)
	w.WriteJson(types.Json{
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
		a.Util.SimpleJsonResponse(w, 403, "Invalid username or password, please try again.")
		return
	} else {
		// err is nil if successful, error if comparison failed
		if err := bcrypt.CompareHashAndPassword(passwordHash, password); err != nil {
			a.Util.SimpleJsonResponse(w, 403, "Invalid username or password, please try again.")
			return
		} else {
			// Create an authentication node and return it to client
			sessionid := a.Svc.SetGetNewSessionId(handle)

			w.WriteHeader(201)
			w.WriteJson(types.Json{
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
	if a.Svc.UnsetSessionId(a.getSessionId(r)) {
		w.WriteHeader(204)
		return
	} else {
		a.Util.SimpleJsonResponse(w, 403, "Cannot invalidate token because it is missing")
		return
	}
}

func (a Api) ChangePassword(w rest.ResponseWriter, r *rest.Request) {
	if !a.authenticate(r) {
		a.Util.FailedToAuthenticate(w)
		return
	}
	user := struct {
		Password           string
		NewPassword        string
		ConfirmNewPassword string
	}{}

	if err := r.DecodeJsonPayload(&user); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if handle, success := a.Svc.GetHandleFromAuthorization(a.getSessionId(r)); !success {
		w.WriteHeader(400)
		w.WriteJson(types.Json{
			"Response":  "Unexpected failure to retrieve owner of session",
			"Handle":    handle,
			"Success":   success,
			"SessionId": a.getSessionId(r),
		})
		return
	} else {
		password := []byte(user.Password)
		newPassword := user.NewPassword
		confirmNewPassword := user.ConfirmNewPassword

		// Password checks
		if newPassword != confirmNewPassword {
			w.WriteHeader(400)
			w.WriteJson(types.Json{
				"Response": "Passwords do not match",
			})
			return
		} else if len(newPassword) < MIN_PASS_LENGTH {
			a.Util.SimpleJsonResponse(w, 400, "Passwords must be at least 8 characters long")
			return
		}

		if passwordHash, ok := a.Svc.GetPasswordHash(handle); !ok {
			a.Util.SimpleJsonResponse(w, 400, "Invalid username or password, please try again.")
			return
		} else {
			// err is nil if successful, error
			if err := bcrypt.CompareHashAndPassword(passwordHash, password); err != nil {
				w.WriteHeader(400)
				w.WriteJson(types.Json{
					"Response": "Invalid username or password, please try again.",
				})
				return
			} else if err := bcrypt.CompareHashAndPassword(passwordHash, []byte(newPassword)); err == nil {
				a.Util.SimpleJsonResponse(w, 400, "Current/new password are same, please provide a new password.")
				return
			} else {
				if hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10); err != nil {
					rest.Error(w, err.Error(), http.StatusInternalServerError)
					return
				} else {
					hashed_new_pass := string(hash)
					// Now set the new password
					if !a.Svc.SetNewPassword(handle, hashed_new_pass) {
						a.Util.SimpleJsonResponse(w, 400, "Password change unsuccessful")
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
			w.WriteJson(types.Json{
				"Results":  nil,
				"Response": "Search failed",
				"Reason":   "Malformed limit",
				"Count":    0,
			})
			return

		} else {
			if intval > 100 || intval < 1 {
				w.WriteHeader(400)
				w.WriteJson(types.Json{
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
			w.WriteJson(types.Json{
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
		w.WriteJson(types.Json{
			"Results":  nil,
			"Response": "Search failed",
			"Reason":   "Missing required sort parameter",
			"Count":    0,
		})
		return
	} else if sortType[0] != "handle" && sortType[0] != "joined" {
		w.WriteHeader(200)
		w.WriteJson(types.Json{
			"Results":  nil,
			"Response": "Search failed",
			"Reason":   "No such sort " + sortType[0],
			"Count":    0,
		})
		return
	}

	results, count := a.Svc.SearchForUsers(circle, nameprefix, skip, limit, sort)

	w.WriteHeader(200)
	w.WriteJson(types.Json{
		"Results":  results,
		"Response": "Search complete",
		"Count":    count,
	})
}

func (a Api) DeleteUser(w rest.ResponseWriter, r *rest.Request) {
	if !a.authenticate(r) {
		a.Util.FailedToAuthenticate(w)
		return
	}
	credentials := struct {
		Password string
	}{}
	if err := r.DecodeJsonPayload(&credentials); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if handle, success := a.Svc.GetHandleFromAuthorization(a.getSessionId(r)); !success {
		w.WriteHeader(400)
		w.WriteJson(types.Json{
			"Response":  "Unexpected failure to retrieve owner of session",
			"Handle":    handle,
			"Success":   success,
			"SessionId": a.getSessionId(r),
		})
		return
	} else {
		password := []byte(credentials.Password)

		if passwordHash, ok := a.Svc.GetPasswordHash(handle); !ok {
			w.WriteHeader(400)
			w.WriteJson(types.Json{
				"Response": "Invalid username or password, please try again.",
			})
			return
		} else {
			// err is nil if successful, error
			if err := bcrypt.CompareHashAndPassword(passwordHash, password); err != nil {
				w.WriteHeader(400)
				w.WriteJson(types.Json{
					"Response": "Invalid username or password, please try again.",
				})
				return
			} else {
				if deleted := a.Svc.DeleteUser(handle); !deleted {
					w.WriteHeader(400)
					w.WriteJson(types.Json{
						"Response": "Unexpected failure to delete user",
					})
					return
				} else {
					w.WriteHeader(204)
				}
			}
		}
	}
}

//
// Circles
//

func (a Api) NewCircle(w rest.ResponseWriter, r *rest.Request) {
	if !a.authenticate(r) {
		a.Util.FailedToAuthenticate(w)
		return
	}
	payload := struct {
		CircleName string
		Public     bool
	}{}
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if handle, success := a.Svc.GetHandleFromAuthorization(a.getSessionId(r)); !success {
		w.WriteHeader(500)
		w.WriteJson(types.Json{
			"Response":  "Unexpected failure to retrieve owner of session",
			"Handle":    handle,
			"Success":   success,
			"SessionId": a.getSessionId(r),
		})
		return
	} else {
		circleName := payload.CircleName
		isPublic := payload.Public

		if circleName == GOLD || circleName == BROADCAST {
			a.Util.SimpleJsonResponse(w, 403, circleName+" is a reserved circle name")
			return
		}

		if circleid, ok := a.Svc.NewCircle(handle, circleName, isPublic); !ok {
			a.Util.SimpleJsonResponse(w, 400, "Unexpected failure to create circle")
			return
		} else {
			w.WriteHeader(201)
			w.WriteJson(types.Json{
				"response": "Created new circle!",
				"chief":    handle,
				"name":     circleName,
				"public":   isPublic,
				"id":       circleid,
			})
		}
	}
}

//
// Messages
//

type MessageData struct {
	Id      string
	Url     string
	Author  string
	Content string
	Created time.Time
}

/**
 * Create a new, unpublished message
 */
func (a Api) NewMessage(w rest.ResponseWriter, r *rest.Request) {
	if !a.authenticate(r) {
		a.Util.FailedToAuthenticate(w)
		return
	}
	payload := struct {
		Content string
		Circles []string
	}{}
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handle, ok := a.Svc.GetHandleFromAuthorization(a.getSessionId(r))
	if !ok {
		w.WriteHeader(500)
		w.WriteJson(types.Json{
			"Response": "Unexpected failure to retrieve owner of session",
		})
		return
	}

	content := payload.Content
	circles := payload.Circles

	if payload.Content == "" {
		a.Util.SimpleJsonResponse(w, 400, "Please enter some content for your message")
		return
	}

	if messageid, success := a.Svc.NewMessage(handle, content); !success {
		a.Util.SimpleJsonResponse(w, 400, "No message created")
		return
	} else {
		if len(circles) > 0 {
			for _, circleid := range circles {
				if !a.Svc.PublishMessageToCircle(messageid, circleid) {
					a.Util.SimpleJsonResponse(w, 400, "Failed to publish to one of circles provided")
					return
				}
			}
		}
		w.WriteHeader(201)
		w.WriteJson(types.Json{
			"Response":    "Successfully created message for " + handle,
			"Id":          messageid,
			"PublishedTo": circles,
		})
	}
}

/**
 * Get messages authored by user
 */
func (a Api) GetAuthoredMessages(w rest.ResponseWriter, r *rest.Request) {
	if !a.authenticate(r) {
		a.Util.FailedToAuthenticate(w)
		return
	}

	if author, success := a.Svc.GetHandleFromAuthorization(a.getSessionId(r)); !success {
		a.Util.FailedToDetermineHandleFromSession(w)
		return
	} else {
		messages := a.Svc.GetMessagesByHandle(author)
		messageData := make([]MessageData, len(messages))

		for i := 0; i < len(messages); i++ {
			messageData[i] = MessageData{
				messages[i].Id,
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
		w.WriteJson(types.Json{
			"Response": "Found messages for user " + author,
			"Objects":  string(b),
			"Count":    len(messageData),
		})
	}
}

func (a Api) GetMessageById(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")

	if !a.authenticate(r) {
		a.Util.FailedToAuthenticate(w)
		return
	}

	handle, success := a.Svc.GetHandleFromAuthorization(a.getSessionId(r))
	if !success {
		w.WriteHeader(400)
		w.WriteJson(types.Json{
			"Response":  "Unexpected failure to retrieve owner of session",
			"Handle":    handle,
			"Success":   success,
			"SessionId": a.getSessionId(r),
		})
		return
	}

	if message, success := a.Svc.GetMessageById(handle, id); success {
		data := MessageData{
			message.Id,
			"<url>:<port>/api/messages/" + message.Id, // hard-coded url/port...
			message.Author,
			message.Content,
			message.Created,
		}

		if b, err := encoding.Marshal(data); err != nil {
			panicErr(err)
		} else {
			w.WriteHeader(200)
			w.WriteJson(types.Json{
				"Response": "Found message!",
				"Object":   string(b),
			})
		}
	} else {
		a.Util.SimpleJsonResponse(w, 404, "No such message in any circle you can see")
		return
	}
}

func (a Api) GetMessagesByHandle(w rest.ResponseWriter, r *rest.Request) {
	w.WriteHeader(405)
	w.WriteJson(types.Json{
		"Response": "Unimplemented",
	})
}

func (a Api) EditMessage(w rest.ResponseWriter, r *rest.Request) {
	if !a.authenticate(r) {
		a.Util.FailedToAuthenticate(w)
		return
	}

	payload := make(types.JsonArray, 0)
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	messageid := r.PathParam("id")
	handle, ok := a.Svc.GetHandleFromAuthorization(a.getSessionId(r))
	if !ok {
		a.Util.FailedToDetermineHandleFromSession(w)
		return
	}

	// Validate input of patch objects
	for i, obj := range payload {
		index := strconv.Itoa(i)
		if op, ok := obj["op"].(string); !ok {
			a.Util.SimpleJsonReason(w, 400, "missing `op` parameter in object "+index)
			return
		} else if resource, ok := obj["resource"].(string); !ok {
			a.Util.SimpleJsonReason(w, 400, "missing `resource` parameter in object "+index)
			return
		} else if value, ok := obj["value"].(string); !ok {
			a.Util.SimpleJsonReason(w, 400, "missing `value` parameter in object "+index)
			return
		} else {
			if op == "update" {
				if resource != "content" && resource != "image" {
					a.Util.SimpleJsonReason(w, 400, "Message only allows update to (content|image) at object "+index)
					return
				} else if resource == "content" && value == "" {
					a.Util.SimpleJsonReason(w, 400, "Cannot update message content to empty at "+index)
					return
				} else if resource == "image" {
					a.Util.SimpleJsonResponse(w, 405, "Edit message image value has yet to be implemented")
					return
				}
			} else if op == "publish" && resource == "circle" {
				if !a.Svc.UserCanPublishTo(handle, value) {
					a.Util.SimpleJsonReason(w, 400, "Could not publish message to circle "+value)
					return
				}
			} else if op == "unpublish" && resource == "circle" {
				if !a.Svc.UserCanRetractPublication(handle, messageid, value) {
					a.Util.SimpleJsonReason(w, 400, "Cannot unpublish message, specified published relation not found")
					return
				}
			} else {
				a.Util.SimpleJsonReason(w, 400, "Malformed patch request at object "+index)
				return
			}
		}
	}

	// Service requests
	for i, obj := range payload {
		op, _ := obj["op"].(string)
		resource, _ := obj["resource"].(string)
		value, _ := obj["value"].(string)

		if op == "update" {
			if resource == "content" {
				a.Svc.UpdateContentOfMessage(messageid, value)
			}
		} else if op == "publish" {
			a.Svc.PublishMessageToCircle(messageid, value)
		} else if op == "unpublish" {
			a.Svc.UnpublishMessageFromCircle(messageid, value)
		} else {
			a.Util.SimpleJsonResponse(w, 500, "Unexpected failure to fulfill service request at "+strconv.Itoa(i))
			return
		}
	}

	w.WriteHeader(200)
	w.WriteJson(types.Json{
		"response": "Successfully patched message " + messageid,
		"changes":  len(payload),
	})
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
		a.Util.FailedToAuthenticate(w)
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
	w.WriteJson(types.Json{
		"Response": "Success!",
		"Deleted":  deleted[0].Count,
	})
}

//
// Social
//

func (a Api) BlockUser(w rest.ResponseWriter, r *rest.Request) {
	if !a.authenticate(r) {
		a.Util.FailedToAuthenticate(w)
		return
	}

	payload := struct {
		Target string
	}{}
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if handle, success := a.Svc.GetHandleFromAuthorization(a.getSessionId(r)); !success {
		a.Util.FailedToDetermineHandleFromSession(w)
		return
	} else {
		target := payload.Target

		if !a.Svc.UserExists(target) {
			a.Util.SimpleJsonResponse(w, 400, "Bad request, user "+target+" wasn't found")
			return
		}

		a.Svc.KickTargetFromCircles(handle, target)

		if !a.Svc.CreateBlockFromTo(handle, target) {
			a.Util.SimpleJsonResponse(w, 400, "Unexpected failure to block user")
		} else {
			a.Util.SimpleJsonResponse(w, 200, "User "+target+" has been blocked")
		}
	}
}

func (a Api) JoinDefault(w rest.ResponseWriter, r *rest.Request) {
	if !a.authenticate(r) {
		a.Util.FailedToAuthenticate(w)
		return
	}

	payload := struct {
		Target string
	}{}
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if handle, success := a.Svc.GetHandleFromAuthorization(a.getSessionId(r)); !success {
		w.WriteHeader(400)
		w.WriteJson(types.Json{
			"Response":  "Unexpected failure to retrieve owner of session",
			"Handle":    handle,
			"Success":   success,
			"SessionId": a.getSessionId(r),
		})
		return
	} else {
		target := payload.Target

		if !a.Svc.UserExists(target) {
			a.Util.SimpleJsonResponse(w, 400, "Bad request, user "+target+" wasn't found")
			return
		}

		if a.Svc.BlockExistsFromTo(target, handle) {
			a.Util.SimpleJsonResponse(w, 403, "Server refusal to comply with join request")
			return
		}

		if !a.Svc.JoinBroadcast(handle, target) {
			a.Util.SimpleJsonResponse(w, 400, "Unexpected failure to join Broadcast")
		} else {
			a.Util.SimpleJsonResponse(w, 201, "JoinDefault request successful!")
		}
	}
}

/**
 * Allows joining by (target, circlename) or (circleid) candidate keys
 */
func (a Api) Join(w rest.ResponseWriter, r *rest.Request) {
	if !a.authenticate(r) {
		a.Util.FailedToAuthenticate(w)
		return
	}

	payload := struct {
		Target   string
		Circle   string
		CircleId string
	}{}
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if handle, success := a.Svc.GetHandleFromAuthorization(a.getSessionId(r)); !success {
		w.WriteHeader(400)
		w.WriteJson(types.Json{
			"Response":  "Unexpected failure to retrieve owner of session",
			"Handle":    handle,
			"Success":   success,
			"SessionId": a.getSessionId(r),
		})
		return
	} else {
		target := payload.Target
		circle := payload.Circle
		circleid := payload.CircleId

		if !a.Svc.UserExists(target) {
			a.Util.SimpleJsonResponse(w, 400, "Bad request, user "+target+" wasn't found")
			return
		}

		if a.Svc.BlockExistsFromTo(target, handle) {
			a.Util.SimpleJsonResponse(w, 403, "Server refusal to comply with join request")
			return
		}

		if circleid == "" {
			if id := a.Svc.GetCircleId(target, circle); id == "" {
				a.Util.SimpleJsonResponse(w, 404, "Could not find target circle, join failed")
				return
			} else {
				circleid = id
			}
		}

		if !a.Svc.CanSeeCircle(handle, circleid) {
			a.Util.SimpleJsonResponse(w, 404, "Could not find target circle, join failed")
			return
		}

		if a.Svc.JoinCircle(handle, circleid) {
			a.Util.SimpleJsonResponse(w, 201, "Join request successful!")
		} else {
			a.Util.SimpleJsonResponse(w, 400, "Unexpected failure to join circle, join failed")
		}
	}
}
