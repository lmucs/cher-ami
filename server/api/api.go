package api

import (
	"../service"
	apiutil "./api-util"
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
		a.Util.SimpleJsonResponse(w, 400, "Handle is a required field for signup")
		return
	} else if email == "" {
		a.Util.SimpleJsonResponse(w, 400, "Email is a required field for signup")
		return
	}

	// Password checks
	if password != confirm_password {
		a.Util.SimpleJsonResponse(w, 400, "Passwords do not match")
		return
	} else if len(password) < MIN_PASS_LENGTH {
		a.Util.SimpleJsonResponse(w, 400, "Passwords must be at least 8 characters long")
		return
	}

	// Ensure unique handle
	if unique, err := a.Svc.HandleIsUnique(handle); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !unique {
		a.Util.SimpleJsonResponse(w, 400, "Sorry, handle or email is already taken")
		return
	}

	// Ensure unique email
	if unique, err := a.Svc.EmailIsUnique(email); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !unique {
		a.Util.SimpleJsonResponse(w, 400, "Sorry, handle or email is already taken")
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
		a.Util.SimpleJsonResponse(w, 400, "Unexpected failure to create new user")
		return
	}

	if !a.Svc.MakeDefaultCirclesFor(handle) {
		a.Util.SimpleJsonResponse(w, 400, "Unexpected failure to make default circles")
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
		w.WriteJson(json{
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
			w.WriteJson(json{
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
				w.WriteJson(json{
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
		w.WriteJson(json{
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
		w.WriteHeader(400)
		w.WriteJson(json{
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

		if !a.Svc.NewCircle(handle, circleName, isPublic) {
			a.Util.SimpleJsonResponse(w, 400, "Unexpected failure to create circle")
			return
		}

		a.Util.SimpleJsonResponse(w, 201, "Created new circle "+circleName+" for "+handle)
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
	payload := struct {
		Content string
	}{}
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !a.authenticate(r) {
		a.Util.FailedToAuthenticate(w)
		return
	}

	handle, ok := a.Svc.GetHandleFromAuthorization(a.getSessionId(r))
	if !ok {
		w.WriteHeader(400)
		w.WriteJson(json{
			"Response": "Unexpected failure to retrieve owner of session",
		})
		return
	}

	content := payload.Content

	if payload.Content == "" {
		a.Util.SimpleJsonResponse(w, 400, "Please enter some content for your message")
		return
	}

	if id, success := a.Svc.NewMessage(handle, content); !success {
		a.Util.SimpleJsonResponse(w, 400, "No message created")
	} else {
		w.WriteHeader(201)
		w.WriteJson(json{
			"Response":  "Successfully created message for " + handle,
			"Id":        id,
			"Published": false,
		})
	}
}

/**
 * Publishes a message identified by it's lastSaved time to a specific circle owned
 * by the user.
 */
func (a Api) PublishMessage(w rest.ResponseWriter, r *rest.Request) {
	if !a.authenticate(r) {
		a.Util.FailedToAuthenticate(w)
		return
	}
	payload := struct {
		CircleId  string
		MessageId string
	}{}
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
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
		circleid := payload.CircleId
		messageid := payload.MessageId

		if !a.Svc.UserIsMemberOf(author, circleid) {
			w.WriteHeader(401)
			w.WriteJson(json{
				"Response": "Refusal to comply with request",
				"Reason":   "You are not a member or owner of the specified circle",
			})
			return
		}

		if !a.Svc.CanSeeCircle(author, circleid) {
			a.Util.SimpleJsonResponse(w, 400, "Could not find specified circle to publish to")
			return
		} else if !a.Svc.MessageExists(messageid) {
			a.Util.SimpleJsonResponse(w, 400, "Could not find intended message for publishing")
			return
		}

		if !a.Svc.PublishMessage(messageid, circleid) {
			a.Util.SimpleJsonResponse(w, 400, "Bad request, no message published")
		} else {
			a.Util.SimpleJsonResponse(w, 201, "Success! Published message to "+circleid)
		}
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
		w.WriteJson(json{
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
		w.WriteJson(json{
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
			w.WriteJson(json{
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
	w.WriteJson(json{
		"Response": "Unimplemented",
	})
}

func (a Api) EditMessage(w rest.ResponseWriter, r *rest.Request) {
	payload := struct {
		Circles []string
		Content string
	}{}
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	circles := payload.Circles
	content := payload.Content
	messageid := r.PathParam("id")

	if !a.authenticate(r) {
		a.Util.FailedToAuthenticate(w)
		return
	}

	if content != "" {
		a.Svc.EditMessageContent(a.getSessionId(r), content)
	}

	for i, circleid := range circles {
		result := a.Svc.PublishMessage(circleid, messageid)
		if result == "bad-circle" {
			w.WriteHeader(400)
			w.WriteJson(json{
				"Response": "Failed to patch message",
				"Reason":   "Some specified circle did not exist, or could not be published to",
			})
		} else if result == "no-such-message" {
			w.WriteHeader(400)
			w.WriteJson(json{
				"Response": "Failed to patch message",
				"Reason":   "You are not the author of any such message",
			})
		}
	}
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
	w.WriteJson(json{
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
		w.WriteHeader(400)
		w.WriteJson(json{
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

		a.Svc.RevokeMembershipBetween(handle, target)

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
		w.WriteJson(json{
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
		w.WriteJson(json{
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
