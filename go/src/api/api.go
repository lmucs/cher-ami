package api

import (
	service "../service"
	"fmt"
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
	Svc *service.Svc
}

/**
 * Constructor
 */
func NewApi(uri string) *Api {
	api := &Api{
		service.NewService(uri),
	}
	return api
}

// Circle constants
const (
	GOLD            = "Gold"
	BROADCAST       = "Broadcast"
	MIN_PASS_LENGTH = 8
)

//
// API util
//

func (a Api) authenticate(r *rest.Request) (success bool) {
	if sessionid := r.Header.Get("Authentication"); sessionid != "" {
		return a.Svc.GoodSessionCredentials(sessionid)
	} else {
		return false
	}
}

//
// API
//

//
// Credentials
//

/*
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
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Handle is a required field for signup",
		})
		return
	} else if email == "" {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Email is a required field for signup",
		})
		return
	}

	// Password checks
	if password != confirm_password {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Passwords do not match",
		})
		return
	} else if len(password) < MIN_PASS_LENGTH {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Passwords must be at least 8 characters long",
		})
		return
	}

	// Ensure unique handle
	if unique, err := a.Svc.HandleIsUnique(handle); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !unique {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Sorry, handle or email is already taken",
		})
		return
	}

	// Ensure unique email
	if unique, err := a.Svc.EmailIsUnique(email); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !unique {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Sorry, handle or email is already taken",
		})
		return
	}

	var hashed_pass string
	if hash, err := bcrypt.GenerateFromPassword([]byte(password), 10); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		hashed_pass = string(hash)
	}
	if err := a.Svc.CreateNewUser(
		handle,
		email,
		hashed_pass,
	); err != nil {
		panicErr(err)
	}

	if err := a.Svc.MakeDefaultCirclesFor(handle); err != nil {
		panicErr(err)
	}

	w.WriteHeader(201)
	w.WriteJson(map[string]string{
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
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Invalid username or password, please try again.",
		})
		return
	} else {
		// err is nil if successful, error if comparison failed
		if err := bcrypt.CompareHashAndPassword(passwordHash, password); err != nil {
			w.WriteHeader(400)
			w.WriteJson(map[string]string{
				"Response": "Invalid username or password, please try again.",
			})
			return
		} else {
			// Create an authentication node and return it to client
			sessionid := a.Svc.SetGetNewSessionId(handle)

			w.WriteHeader(200)
			w.WriteJson(map[string]string{
				"Response":  "Logged in " + handle + ". Note your session id.",
				"sessionid": sessionid,
			})
			w.Header().Add("Authentication", sessionid)
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
		w.WriteHeader(200)
		w.WriteJson(map[string]string{
			"Response": "Goodbye " + handle + ", have a nice day",
		})
		return
	} else {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "That user doesn't exist",
		})
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

	err := r.DecodeJsonPayload(&user)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handle := user.Handle
	password := []byte(user.Password)
	newPassword := user.NewPassword
	confirmNewPassword := user.ConfirmNewPassword

	if !a.authenticate(r) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Failed to authenticate user request",
		})
		return
	}

	// Password checks
	if newPassword != confirmNewPassword {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Passwords do not match",
		})
		return
	} else if len(newPassword) < MIN_PASS_LENGTH {
		fmt.Println("")
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Passwords must be at least 8 characters long",
		})
		return
	}

	if passwordHash, ok := a.Svc.GetPasswordHash(handle); !ok {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Invalid username or password, please try again.",
		})
		return
	} else {
		// err is nil if successful, error
		if err := bcrypt.CompareHashAndPassword(passwordHash, password); err != nil {
			w.WriteHeader(400)
			w.WriteJson(map[string]string{
				"Response": "Invalid username or password, please try again.",
			})
			return
		} else if err := bcrypt.CompareHashAndPassword(passwordHash, []byte(newPassword)); err == nil {
			w.WriteHeader(400)
			w.WriteJson(map[string]string{
				"Response": "Current/new password are same, please provide a new password.",
			})
			return
		} else {
			if hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10); err != nil {
				rest.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				hashed_new_pass := string(hash)
				// Now set the new password
				if !a.Svc.SetNewPassword(handle, hashed_new_pass) {
					w.WriteHeader(400)
					w.WriteJson(map[string]string{
						"Response": "Password change unsuccessful",
					})
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
			w.WriteJson(map[string]interface{}{
				"Results":  nil,
				"Response": "Search failed",
				"Reason":   "Malformed limit",
				"Count":    0,
			})
			return

		} else {
			if intval > 100 || intval < 1 {
				w.WriteHeader(400)
				w.WriteJson(map[string]interface{}{
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
			w.WriteJson(map[string]interface{}{
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
		w.WriteJson(map[string]interface{}{
			"Results":  nil,
			"Response": "Search failed",
			"Reason":   "Missing required sort parameter",
			"Count":    0,
		})
		return
	} else if sortType[0] != "handle" && sortType[0] != "joined" {
		w.WriteHeader(200)
		w.WriteJson(map[string]interface{}{
			"Results":  nil,
			"Response": "Search failed",
			"Reason":   "No such sort " + sortType[0],
			"Count":    0,
		})
		return
	}

	results, count := a.Svc.SearchForUsers(circle, nameprefix, skip, limit, sort)

	w.WriteHeader(200)
	w.WriteJson(map[string]interface{}{
		"Results":  results,
		"Response": "Search complete.",
		"Count":    count,
	})
}

func (a Api) GetUsers(w rest.ResponseWriter, r *rest.Request) {
	res := []struct {
		Handle string    `json:"user.handle"`
		Joined time.Time `json:"user.joined"`
	}{}

	err := a.Svc.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (user:User)
            RETURN user.handle, user.joined
            ORDER BY user.handle
        `,
		Parameters: neoism.Props{},
		Result:     &res,
	})
	panicErr(err)

	if len(res) > 0 {
		users := []map[string]string{}

		for i := range res {
			user := map[string]string{
				"Handle": res[i].Handle,
				"Joined": res[i].Joined.Format("Jan 2, 2006 at 3:04 PM (MST)"),
			}
			users = append(users, user)
		}

		w.WriteHeader(200)
		w.WriteJson(users)
	} else {
		w.WriteHeader(404)
		w.WriteJson(map[string]string{
			"Response": "No results found",
		})
	}
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
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Failed to authenticate user request",
		})
		return
	}

	if passwordHash, ok := a.Svc.GetPasswordHash(handle); !ok {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Invalid username or password, please try again.",
		})
		return
	} else {
		// err is nil if successful, error
		if err := bcrypt.CompareHashAndPassword(passwordHash, password); err != nil {
			w.WriteHeader(400)
			w.WriteJson(map[string]string{
				"Response": "Invalid username or password, please try again.",
			})
			return
		} else {
			if deleted := a.Svc.DeleteUser(handle); !deleted {
				w.WriteHeader(400)
				w.WriteJson(map[string]string{
					"Response": "Unexpected failure to delete user",
				})
				return
			}
			w.WriteHeader(200)
			w.WriteJson(map[string]string{
				"Response": "Deleted " + handle,
			})
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
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Failed to authenticate user request",
		})
		return
	}

	if circleName == GOLD || circleName == BROADCAST {
		w.WriteHeader(403)
		w.WriteJson(map[string]string{
			"Response": circleName + " is a reserved circle name",
		})
		return
	}

	err := a.Svc.NewCircle(handle, circleName, isPublic)
	panicErr(err)

	w.WriteHeader(201)
	w.WriteJson(map[string]string{
		"Response": "Created new circle " + circleName + " for " + handle,
	})
}

func (a Api) makeDefaultCircles(handle string) {
	made := []struct {
		Handle    string `json:"u.handle"`
		Gold      string `json:"g.name"`
		Broadcast string `json:"br.name"`
	}{}
	dberr := a.Svc.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (p:PublicDomain {u:true})
            MATCH (u:User)
            WHERE u.handle = {handle}
            CREATE (g:Circle {name: {gold}})
            CREATE (br:Circle {name: {broadcast}})
            CREATE (u)-[:CHIEF_OF]->(g)
            CREATE (u)-[:CHIEF_OF]->(br)
            CREATE UNIQUE (br)-[:PART_OF]->(p)
            RETURN u.handle, g.name, br.name
        `,
		Parameters: neoism.Props{
			"handle":    handle,
			"gold":      GOLD,
			"broadcast": BROADCAST,
		},
		Result: &made,
	})
	panicErr(dberr)
	if len(made) != 1 {
		panic(fmt.Sprintf("Incorrect results len in query1()\n\tgot %d, expected 1\n", len(made)))
	}
}

//
// Messages
//

/**
 * Expects a json post with "handle", "sessionid", "content"
 */
func (a Api) NewMessage(w rest.ResponseWriter, r *rest.Request) {
	payload := struct {
		Handle    string
		SessionId string
		Content   string
	}{}
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handle := payload.Handle
	sessionid := payload.SessionId
	content := payload.Content

	if !a.authenticate(r) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Failed to authenticate user request",
		})
		return
	}

	if payload.Content == "" {
		rest.Error(w, "Please enter some content for your message", 400)
		return
	}

	created := []struct {
		Content  string      `json:"message.content"`
		Relation neoism.Node `json:"r"`
	}{}
	now := time.Now().Local()
	err := a.Svc.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (user:User {handle: {handle}, sessionid: {sessionid}})
            CREATE (message:Message {content: {content}, created: {now}, lastsaved: {now}})
            CREATE (user)-[r:WROTE]->(message)
            RETURN message.content, r
        `,
		Parameters: neoism.Props{
			"handle":    handle,
			"sessionid": sessionid,
			"content":   content,
			"now":       now,
		},
		Result: &created,
	})
	panicErr(err)

	if len(created) != 1 {
		w.WriteHeader(500)
		w.WriteJson(map[string]string{
			"Response": "No message created",
		})
	} else {
		w.WriteHeader(201)
		w.WriteJson(map[string]interface{}{
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
		LastSaved time.Time
		Circle    string
	}{}
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handle := payload.Handle
	lastsaved := payload.LastSaved
	circle := payload.Circle

	if !a.authenticate(r) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Failed to authenticate user request",
		})
		return
	}

	if !a.Svc.CircleExists(handle, circle) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Bad Request, could not find specified circle to publish to",
		})
		return
	}

	if !a.Svc.MessageExists(handle, lastsaved) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Bad Request, could not find intended message for publishing",
		})
		return
	}

	created := []struct {
		Count int `json:"count(r)"`
	}{}
	err := a.Svc.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (u:User)
            WHERE u.handle={handle}
            MATCH (u)-[:CHIEF_OF]->(c:Circle)
            WHERE c.name={name}
            MATCH (u)-[:WROTE]->(m:Message)
            WHERE m.lastsaved={lastsaved}
            CREATE (m)-[r:PUB_TO]->(c)
            SET r.publishedat={date}
            RETURN count(r)
        `,
		Parameters: neoism.Props{
			"handle":    handle,
			"name":      circle,
			"lastsaved": lastsaved,
			"date":      time.Now().Local(),
		},
		Result: &created,
	})
	panicErr(err)

	if created[0].Count > 0 {
		w.WriteHeader(201)
		w.WriteJson(map[string]string{
			"Response": "Success! Published message to " + circle,
		})
	} else {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Bad request, no message published",
		})
	}
}

/**
 * Get messages authored by user
 * Expects query parameters "handle" and "sessionid"
 */
func (a Api) GetAuthoredMessages(w rest.ResponseWriter, r *rest.Request) {
	querymap := r.URL.Query()

	// Check query parameters
	if _, ok := querymap["handle"]; !ok {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Bad Request, not enough parameters to authenticate user",
		})
		return
	}

	handle := querymap["handle"][0]

	if !a.authenticate(r) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Failed to authenticate user request",
		})
		return
	}

	messages := []struct {
		Content   string    `json:"message.content"`
		LastSaved time.Time `json:"message.lastsaved"`
	}{}
	err := a.Svc.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (user:User {handle: {handle}})-[:WROTE]->(message:Message)
            RETURN message.content, message.lastsaved
        `,
		Parameters: neoism.Props{
			"handle": handle,
		},
		Result: &messages,
	})
	panicErr(err)

	w.WriteHeader(200)
	w.WriteJson(messages)
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
		w.WriteJson(map[string]string{
			"Response": "Bad Request, not enough parameters to authenticate user",
		})
		return
	}

	handle := querymap["handle"][0]

	if !a.authenticate(r) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Failed to authenticate user request",
		})
		return
	}

	if !a.Svc.UserExists(author) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
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
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Failed to authenticate user request",
		})
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
	w.WriteJson(map[string]interface{}{
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
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Failed to authenticate user request",
		})
		return
	}

	if !a.Svc.UserExists(target) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Bad request, user " + target + " wasn't found",
		})
		return
	}

	a.Svc.RevokeMembershipBetween(handle, target)

	// Block user
	if success, err := a.Svc.AddBlockedRelation(handle, target); err != nil {
		panicErr(err)
	} else if success {
		w.WriteHeader(200)
		w.WriteJson(map[string]string{
			"Response": "User " + target + " has been blocked",
		})
	} else {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Unexpected failure to block user",
		})
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
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Failed to authenticate user request",
		})
		return
	}

	if !a.Svc.UserExists(target) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Bad request, user " + target + " wasn't found",
		})
		return
	}

	if a.Svc.BlockExistsFromTo(handle, target) {
		w.WriteHeader(403)
		w.WriteJson(map[string]string{
			"Response": "Server refusal to comply with join request",
		})
		return
	}

	if at, did_join := a.Svc.JoinBroadcast(handle, target); did_join {
		w.WriteHeader(201)
		w.WriteJson(map[string]string{
			"Response": "JoinDefault request successful!",
			"Info":     handle + " added to " + target + "'s broadcast at " + at.Format(time.RFC1123),
		})
	} else {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Unexpected failure to join Broadcast",
		})
	}
}

func (a Api) Join(w rest.ResponseWriter, r *rest.Request) {
	payload := struct {
		Handle string
		Target string
		Circle string
	}{}
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handle := payload.Handle
	target := payload.Target
	circle := payload.Circle

	if !a.authenticate(r) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Failed to authenticate user request",
		})
		return
	}

	if !a.Svc.UserExists(target) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Bad request, user " + payload.Target + " wasn't found",
		})
		return
	}

	if a.Svc.BlockExistsFromTo(handle, target) {
		w.WriteHeader(403)
		w.WriteJson(map[string]string{
			"Response": "Server refusal to comply with join request",
		})
		return
	}

	if !a.Svc.CircleExists(target, circle) {
		w.WriteHeader(404)
		w.WriteJson(map[string]string{
			"Response": "Could not find target circle, join failed",
		})
		return
	}

	if at, did_join := a.Svc.JoinCircle(handle, target, circle); did_join {
		w.WriteHeader(201)
		w.WriteJson(map[string]string{
			"Response": "Join request successful!",
			"Info":     handle + " joined " + circle + " of " + target + " at " + at.Format(time.RFC1123),
		})
	} else {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Unexpected failure to join circle, join failed",
		})
	}
}
