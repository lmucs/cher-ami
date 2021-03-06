package api

import (
	"../types"
	"./service"
	apiutil "./util"
	"github.com/ChimeraCoder/go.crypto/bcrypt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/mccoyst/validate"
	"net/http"
	"strconv"
	"time"
)

//
// Application Types
//

type Api struct {
	Svc       *service.Svc
	Util      *apiutil.Util
	Validator *validate.V
}

/**
 * Constructor
 */
func NewApi(uri string) *Api {
	api := &Api{
		Svc:       service.NewService(uri),
		Util:      &apiutil.Util{},
		Validator: types.NewValidator(),
	}
	return api
}

//
// API util
//

func (a Api) authenticate(r *rest.Request) (success bool) {
	if token := a.getTokenFromHeader(r); token != "" {
		return a.Svc.VerifyAuthToken(token)
	} else {
		return false
	}
}

func (a Api) getTokenFromHeader(r *rest.Request) string {
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
	proposal := types.SignupProposal{}
	if err := r.DecodeJsonPayload(&proposal); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.Validator.ValidateAndTag(proposal, "json"); err != nil {
		a.Util.SimpleJsonValidationReason(w, 400, err)
		return
	}

	handle := proposal.Handle
	email := proposal.Email
	password := proposal.Password
	confirm_password := proposal.ConfirmPassword

	// Password checks
	if password != confirm_password {
		a.Util.SimpleJsonReason(w, 403, "Passwords do not match")
		return
	}

	// Ensure unique handle
	if unique := a.Svc.HandleIsUnique(handle); !unique {
		a.Util.SimpleJsonReason(w, 409, "Sorry, handle or email is already taken")
		return
	}

	// Ensure unique email
	if unique := a.Svc.EmailIsUnique(email); !unique {
		a.Util.SimpleJsonReason(w, 409, "Sorry, handle or email is already taken")
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
		a.Util.SimpleJsonReason(w, http.StatusInternalServerError, "Unexpected failure to create new user")
		return
	}

	if !a.Svc.MakeDefaultCirclesFor(handle) {
		a.Util.SimpleJsonReason(w, http.StatusInternalServerError, "Unexpected failure to make default circles")
		return
	}

	w.WriteHeader(201)
	w.WriteJson(types.Json{
		"response": "Signed up a new user!",
		"handle":   handle,
		"email":    email,
	})
}

func (a Api) Login(w rest.ResponseWriter, r *rest.Request) {
	credentials := types.LoginCredentials{}
	if err := r.DecodeJsonPayload(&credentials); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.Validator.ValidateAndTag(credentials, "json"); err != nil {
		a.Util.SimpleJsonValidationReason(w, 400, err)
		return
	}

	handle := credentials.Handle
	password := []byte(credentials.Password)

	if passwordHash, ok := a.Svc.GetPasswordHash(handle); !ok {
		a.Util.SimpleJsonReason(w, 403, "Invalid username or password, please try again.")
		return
	} else {
		// err is nil if successful, error if comparison failed
		if err := bcrypt.CompareHashAndPassword(passwordHash, password); err != nil {
			a.Util.SimpleJsonReason(w, 403, "Invalid username or password, please try again.")
			return
		} else {
			// Create an Authentication token and return it to client
			if token, ok := a.Svc.SetGetNewAuthToken(handle); !ok {
				a.Util.SimpleJsonReason(w, 500, "Unexpected failure to produce new Authorization token")
			} else {
				w.WriteHeader(201)
				w.WriteJson(types.Json{
					"handle":   handle,
					"response": "Logged in " + handle + ". Note your Authorization token.",
					"token":    token,
				})
				return
			}
		}
	}
}

/**
 * Expects a json post with "handle"
 */
func (a Api) Logout(w rest.ResponseWriter, r *rest.Request) {
	if a.Svc.DestroyAuthToken(a.getTokenFromHeader(r)) {
		w.WriteHeader(204)
		return
	} else {
		a.Util.SimpleJsonReason(w, 403, "Cannot invalidate token because it is missing")
		return
	}
}

//
// User
//
func (a Api) GetUser(w rest.ResponseWriter, r *rest.Request) {
	if !a.authenticate(r) {
		a.Util.FailedToAuthenticate(w)
		return
	}

	target := r.PathParam("handle")

	if handle, ok := a.Svc.GetHandleFromAuthorization(a.getTokenFromHeader(r)); !ok {
		a.Util.FailedToDetermineHandleFromAuthToken(w)
		return
	} else {
		if user, ok := a.Svc.GetVisibleUser(handle, target); !ok {
			a.Util.SimpleJsonReason(w, 500, "Failed to get user "+target)
			return
		} else {
			w.WriteHeader(200)
			w.WriteJson(user)
		}
	}
}

func (a Api) EditUser(w rest.ResponseWriter, r *rest.Request) {
	if !a.authenticate(r) {
		a.Util.FailedToAuthenticate(w)
		return
	}

	attributes := []types.UserPatch{}
	if err := r.DecodeJsonPayload(&attributes); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handle := r.PathParam("handle")
	if h, ok := a.Svc.GetHandleFromAuthorization(a.getTokenFromHeader(r)); !ok {
		a.Util.FailedToDetermineHandleFromAuthToken(w)
		return
	} else if h != handle {
		a.Util.SimpleJsonReason(w, 401, "You are not authorized to modify user "+handle)
		return
	}

	// Validate input of patch objects
	for index, obj := range attributes {
		if err := a.Validator.ValidateAndTag(obj, "json"); err != nil {
			a.Util.PatchValidationReason(w, 400, err, index)
			return
		}
	}

	// Service requests
	for i, obj := range attributes {
		resource := obj.Resource
		value := obj.Value

		if !a.Svc.UpdateUserAttribute(handle, resource, value) {
			a.Util.SimpleJsonReason(w, 500, "Unexpected failure to fulfill service request at "+strconv.Itoa(i))
			return
		}
	}

	w.WriteHeader(200)
	w.WriteJson(types.Json{
		"response": "Successfully updated user " + handle,
		"changes":  len(attributes),
	})

}

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
				"results":  nil,
				"response": "Search failed",
				"reason":   "Malformed limit",
				"count":    0,
			})
			return

		} else {
			if intval > 100 || intval < 1 {
				w.WriteHeader(400)
				w.WriteJson(types.Json{
					"results":  nil,
					"response": "Search failed",
					"reason":   "Limit out of range",
					"count":    0,
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
				"results":  nil,
				"response": "Search failed",
				"reason":   "Malformed skip",
				"count":    0,
			})
			return
		} else {
			skip = intval
		}
	}

	if sortType, ok := querymap["sort"]; !ok {
		w.WriteHeader(400)
		w.WriteJson(types.Json{
			"results":  nil,
			"response": "Search failed",
			"reason":   "Missing required sort parameter",
			"count":    0,
		})
		return
	} else if sortType[0] != "handle" && sortType[0] != "joined" {
		w.WriteHeader(200)
		w.WriteJson(types.Json{
			"results":  nil,
			"response": "Search failed",
			"reason":   "No such sort " + sortType[0],
			"count":    0,
		})
		return
	}

	results, count := a.Svc.SearchForUsers(circle, nameprefix, skip, limit, sort)

	w.WriteHeader(200)
	w.WriteJson(types.Json{
		"results":  results,
		"response": "Search complete",
		"count":    count,
	})
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
		CircleName string `json:"circlename"`
		Public     bool   `json:"public"`
	}{}
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	circleName := payload.CircleName
	isPublic := payload.Public

	if circleName == "" {
		a.Util.SimpleJsonReason(w, 400, "Missing `circlename` parameter")
		return
	} else if circleName == types.GOLD || circleName == types.BROADCAST {
		a.Util.SimpleJsonReason(w, 403, circleName+" is a reserved circle name")
		return
	}

	if handle, success := a.Svc.GetHandleFromAuthorization(a.getTokenFromHeader(r)); !success {
		a.Util.FailedToDetermineHandleFromAuthToken(w)
		return
	} else {
		if circleResponse, ok := a.Svc.NewCircle(handle, circleName, isPublic); !ok {
			a.Util.SimpleJsonReason(w, 400, "Unexpected failure to create circle")
			return
		} else {
			w.WriteHeader(201)
			w.WriteJson(circleResponse)
		}
	}
}

func (a Api) SearchCircles(w rest.ResponseWriter, r *rest.Request) {
	if !a.authenticate(r) {
		a.Util.FailedToAuthenticate(w)
		return
	}

	querymap := r.URL.Query()

	var user string
	var before time.Time
	var limit int

	if val, ok := querymap["limit"]; !ok {
		limit = 20
	} else {
		if intval, err := strconv.Atoi(val[0]); err != nil {
			a.Util.SimpleJsonReason(w, 400, "Malformed limit")
			return
		} else {
			if intval > 100 || intval < 1 {
				a.Util.SimpleJsonReason(w, 400, "Limit out of range")
				return
			} else {
				limit = intval
			}
		}
	}

	// Reveals public and private circles user is apart of. If user parameter is absent
	// will use the logged in user as the target of the query.
	// Empty assumed not to be the name of a user, DWIW
	if val, ok := querymap["user"]; !ok || val[0] == "" {
		if handle, ok := a.Svc.GetHandleFromAuthorization(a.getTokenFromHeader(r)); !ok {
			a.Util.FailedToDetermineHandleFromAuthToken(w)
			return
		} else {
			user = handle
		}
	} else {
		user = val[0]
	}

	if val, ok := querymap["before"]; !ok {
		before = time.Now().Local()
	} else {
		if millis, err := strconv.Atoi(val[0]); err != nil {
			a.Util.SimpleJsonResponse(w, 400, "Malformed duration")
			return
		} else {
			var seconds, nanoseconds int64
			seconds = int64(millis / 1000)
			nanoseconds = int64(millis % 1000 * 1000000)
			before = time.Unix(seconds, nanoseconds)
		}
	}

	results, count := a.Svc.CirclesUserIsPartOf(user, before, limit)

	w.WriteHeader(200)
	w.WriteJson(types.SearchCirclesResponse{
		Results:  results,
		Response: "Search complete",
		Count:    count,
	})
}

//
// Messages
//

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

	handle, ok := a.Svc.GetHandleFromAuthorization(a.getTokenFromHeader(r))
	if !ok {
		a.Util.FailedToDetermineHandleFromAuthToken(w)
		return
	}

	content := payload.Content
	circles := payload.Circles

	if payload.Content == "" {
		a.Util.SimpleJsonReason(w, 400, "Please enter some content for your message")
		return
	}

	if message, ok := a.Svc.NewMessage(handle, content); !ok {
		a.Util.SimpleJsonReason(w, 500, "Unexpected failure to create message")
		return
	} else {
		if len(circles) > 0 {
			for _, circleid := range circles {
				if !a.Svc.PublishMessageToCircle(message.Id, circleid) {
					a.Util.SimpleJsonReason(w, 400, "Failed to publish to one of circles provided")
					return
				}
			}
		}
		w.WriteHeader(201)
		w.WriteJson(message)
	}
}

func (a Api) GetMessages(w rest.ResponseWriter, r *rest.Request) {
	if !a.authenticate(r) {
		a.Util.FailedToAuthenticate(w)
		return
	}

	querymap := r.URL.Query()
	var target, circleid string
	var t, c bool

	if self, ok := a.Svc.GetHandleFromAuthorization(a.getTokenFromHeader(r)); !ok {
		a.Util.FailedToDetermineHandleFromAuthToken(w)
		return
	} else {
		if _, ok := querymap["all"]; ok {
			if messagesView, ok := a.Svc.GetAllMessages(self); !ok {
				a.Util.SimpleJsonReason(w, 400, "Something happened")
			} else {
				w.WriteHeader(200)
				w.WriteJson(messagesView)
				return
			}
		}
		if targetUses, ok := querymap["handle"]; ok {
			target, t = targetUses[0], ok
		}
		if circleidUses, ok := querymap["circleid"]; ok {
			circleid, c = circleidUses[0], ok
		}

		if t && c {
			if messagesView, ok := a.Svc.GetMessagesByTargetInCircle(self, target, circleid); !ok {
				a.Util.SimpleJsonReason(w, 400, "Could not find circle or you lack access rights")
			} else {
				w.WriteHeader(200)
				w.WriteJson(messagesView)
			}
		} else if t && !c && target != self {
			if messagesView, ok := a.Svc.GetPublicMessagesByHandle(self, target); !ok {
				a.Util.SimpleJsonReason(w, 400, "Could not find circle or you lack access rights")
			} else {
				w.WriteHeader(200)
				w.WriteJson(messagesView)
			}
		} else if !t && c {
			if messagesView, ok := a.Svc.GetMessagesInCircle(self, circleid); !ok {
				a.Util.SimpleJsonReason(w, 400, "Could not find circle or you lack access rights")
			} else {
				w.WriteHeader(200)
				w.WriteJson(messagesView)
			}
		} else {
			if messagesView, ok := a.Svc.GetMessageFeedOfSelf(self); !ok {
				a.Util.SimpleJsonReason(w, 500, "Unexpected failure to get feed")
			} else {
				w.WriteHeader(200)
				w.WriteJson(messagesView)
			}
		}
	}
}

func (a Api) GetMessageById(w rest.ResponseWriter, r *rest.Request) {
	if !a.authenticate(r) {
		a.Util.FailedToAuthenticate(w)
		return
	}

	id := r.PathParam("id")

	handle, ok := a.Svc.GetHandleFromAuthorization(a.getTokenFromHeader(r))
	if !ok {
		a.Util.FailedToDetermineHandleFromAuthToken(w)
		return
	}

	if message, ok := a.Svc.GetVisibleMessageById(handle, id); ok {
		w.WriteHeader(200)
		w.WriteJson(message)
	} else {
		a.Util.SimpleJsonReason(w, 404, "No such message with id "+id+" could be found")
		return
	}
}

func (a Api) EditMessage(w rest.ResponseWriter, r *rest.Request) {
	if !a.authenticate(r) {
		a.Util.FailedToAuthenticate(w)
		return
	}

	payload := make([]types.MessagePatch, 0)
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	messageid := r.PathParam("id")
	handle, ok := a.Svc.GetHandleFromAuthorization(a.getTokenFromHeader(r))
	if !ok {
		a.Util.FailedToDetermineHandleFromAuthToken(w)
		return
	}

	// Validate input of patch objects
	for index, obj := range payload {
		if err := a.Validator.ValidateAndTag(obj, "json"); err != nil {
			a.Util.PatchValidationReason(w, 400, err, index)
			return
		} else {
			op := obj.Op
			resource := obj.Resource
			value := obj.Value
			if op == "update" {
				if resource == "image" {
					a.Util.SimpleJsonReason(w, 405, "Edit message image value has yet to be implemented")
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
				a.Util.SimpleJsonReason(w, 400, "Malformed patch request at object "+strconv.Itoa(index))
				return
			}
		}
	}

	// Service requests
	for i, obj := range payload {
		op := obj.Op
		resource := obj.Resource
		value := obj.Value

		if op == "update" {
			if resource == "content" {
				a.Svc.UpdateContentOfMessage(messageid, value)
			}
		} else if op == "publish" {
			a.Svc.PublishMessageToCircle(messageid, value)
		} else if op == "unpublish" {
			a.Svc.UnpublishMessageFromCircle(messageid, value)
		} else {
			a.Util.SimpleJsonReason(w, 500, "Unexpected failure to fulfill service request at "+strconv.Itoa(i))
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
	a.Util.SimpleJsonReason(w, 405, "Unimplemented")
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
	if handle, success := a.Svc.GetHandleFromAuthorization(a.getTokenFromHeader(r)); !success {
		a.Util.FailedToDetermineHandleFromAuthToken(w)
		return
	} else {
		target := payload.Target

		if !a.Svc.UserExists(target) {
			a.Util.SimpleJsonReason(w, 400, "Bad request, user "+target+" wasn't found")
			return
		}

		a.Svc.KickTargetFromCircles(handle, target)

		if !a.Svc.CreateBlockFromTo(handle, target) {
			a.Util.SimpleJsonReason(w, 400, "Unexpected failure to block user")
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

	if handle, success := a.Svc.GetHandleFromAuthorization(a.getTokenFromHeader(r)); !success {
		a.Util.FailedToDetermineHandleFromAuthToken(w)
		return
	} else {
		target := payload.Target

		if !a.Svc.UserExists(target) {
			a.Util.SimpleJsonReason(w, 400, "Bad request, user "+target+" wasn't found")
			return
		}

		if a.Svc.BlockExistsFromTo(target, handle) {
			a.Util.SimpleJsonReason(w, 403, "Server refusal to comply with join request")
			return
		}

		if !a.Svc.JoinBroadcast(handle, target) {
			a.Util.SimpleJsonReason(w, 400, "Unexpected failure to join Broadcast")
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

	if handle, success := a.Svc.GetHandleFromAuthorization(a.getTokenFromHeader(r)); !success {
		a.Util.FailedToDetermineHandleFromAuthToken(w)
		return
	} else {
		target := payload.Target
		circle := payload.Circle
		circleid := payload.CircleId

		if !a.Svc.UserExists(target) {
			a.Util.SimpleJsonReason(w, 400, "Bad request, user "+target+" wasn't found")
			return
		}

		if a.Svc.BlockExistsFromTo(target, handle) {
			a.Util.SimpleJsonReason(w, 403, "Server refusal to comply with join request")
			return
		}

		if circleid == "" {
			if id := a.Svc.GetCircleId(target, circle); id == "" {
				a.Util.SimpleJsonReason(w, 404, "Could not find target circle, join failed")
				return
			} else {
				circleid = id
			}
		}

		if !a.Svc.CanSeeCircle(handle, circleid) {
			a.Util.SimpleJsonReason(w, 404, "Could not find target circle, join failed")
			return
		}

		if a.Svc.JoinCircle(handle, circleid) {
			a.Util.SimpleJsonResponse(w, 201, "Join request successful!")
		} else {
			a.Util.SimpleJsonReason(w, 400, "Unexpected failure to join circle, join failed")
		}
	}
}
