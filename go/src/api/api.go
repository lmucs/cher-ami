package api

import (
	service "../service"
	"fmt"
	"github.com/ChimeraCoder/go.crypto/bcrypt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jmcvetta/neoism"
	"net/http"
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
	GOLD      = "Gold"
	BROADCAST = "Broadcast"
)

//
// API util
//

/**
 * Requests
 */
func (a Api) authenticate(handle string, sessionid string) bool {
	ok, err := a.Svc.GoodSessionCredentials(handle, sessionid)
	if err != nil {
		panicErr(err)
	}

	return ok
}

func (a Api) userExists(handle string) bool {
	found := []struct {
		Handle string `json:"user.handle"`
	}{}
	err := a.Svc.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (user:User {handle: {handle}})
            RETURN user.handle
        `,
		Parameters: neoism.Props{
			"handle": handle,
		},
		Result: &found,
	})
	panicErr(err)

	return len(found) > 0
}

func (a Api) circleExists(target string, circleName string) bool {
	found := []struct {
		Name string `json:"c.name"`
	}{}
	err := a.Svc.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (t:User)
            WHERE t.handle = {target}
            MATCH (t)-[:CHIEF_OF]->(c:Circle)
            WHERE c.name = {name}
            RETURN c.name
        `,
		Parameters: neoism.Props{
			"target": target,
			"name":   circleName,
		},
		Result: &found,
	})
	panicErr(err)

	return len(found) > 0
}

func (a Api) messageExists(handle string, lastSaved time.Time) bool {
	count := []struct {
		Count int `json:"count(m)"`
	}{}
	err := a.Svc.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (u:User {handle: {handle}})
            OPTIONAL MATCH (u)-[:WROTE]->(m:Message {lastsaved: {lastsaved}})
            RETURN count(m)
        `,
		Parameters: neoism.Props{
			"handle":    handle,
			"lastsaved": lastSaved,
		},
		Result: &count,
	})
	panicErr(err)

	return count[0].Count > 0
}

func (a Api) hasBlocked(handle string, target string) bool {
	blocked := []struct {
		Count int `json:"count(r)"`
	}{}
	err := a.Svc.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (u:User {handle: {handle}})
            MATCH (t:User {handle: {target}})
            OPTIONAL MATCH (u)-[r:BLOCKED]->(t)
            RETURN count(r)
        `,
		Parameters: neoism.Props{
			"handle": handle,
			"target": target,
		},
		Result: &blocked,
	})
	panicErr(err)

	return blocked[0].Count > 0
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
		rest.Error(w, err.Error(), http.StatusInternalServerError)
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
	const min_pass_length = 8
	if password != confirm_password {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Passwords do not match",
		})
		return
	} else if len(password) < min_pass_length {
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
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handle := credentials.Handle
	password := []byte(credentials.Password)

	if password_hash, ok := a.Svc.GetPasswordHash(handle); !ok {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Invalid username or password, please try again.",
		})
		return
	} else {
		// err is nil if successful, error
		if err := bcrypt.CompareHashAndPassword(password_hash, password); err != nil {
			w.WriteHeader(400)
			w.WriteJson(map[string]string{
				"Response": "Invalid username or password, please try again.",
			})
			return
		} else {
			// Add session hash to node and return it to client
			if sessionid, err := a.Svc.SetGetNewSessionId(handle); err != nil {
				panicErr(err)
			} else {
				w.WriteJson(map[string]string{
					"Response":  "Logged in " + handle + ". Note your session id.",
					"SessionId": sessionid,
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
	user := struct {
		Handle string
	}{}
	err := r.DecodeJsonPayload(&user)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !a.userExists(user.Handle) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "That user doesn't exist",
		})
		return
	}

	if err := a.Svc.UnsetSessionId(user.Handle); err != nil {
		panicErr(err)
	}

	w.WriteHeader(200)
	w.WriteJson(map[string]string{
		"Response": "Goodbye " + user.Handle + ", have a nice day",
	})
}

//
// User
//

func (a Api) GetUser(w rest.ResponseWriter, r *rest.Request) {
	user := struct {
		Handle string
	}{}
	if err := r.DecodeJsonPayload(&user); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if handle, name, found := a.Svc.GetHandleAndNameOf(user.Handle); found {
		w.WriteHeader(200)
		w.WriteJson(map[string]string{
			"handle": handle,
			"name":   name,
		})
	} else {
		w.WriteHeader(404)
		w.WriteJson(map[string]string{
			"Response": "No results found",
		})
	}
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
		Handle    string
		Password  string
		SessionId string
	}{}
	if err := r.DecodeJsonPayload(&credentials); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handle := credentials.Handle
	password := []byte(credentials.Password)
	sessionid := credentials.SessionId

	if !a.authenticate(handle, sessionid) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Failed to authenticate user request",
		})
		return
	}

	if password_hash, ok := a.Svc.GetPasswordHash(handle); !ok {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Invalid username or password, please try again.",
		})
		return
	} else {
		// err is nil if successful, error
		if err := bcrypt.CompareHashAndPassword(password_hash, password); err != nil {
			w.WriteHeader(400)
			w.WriteJson(map[string]string{
				"Response": "Invalid username or password, please try again.",
			})
			return
		} else {
			if err := a.Svc.DeleteUser(handle); err != nil {
				panicErr(err)
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
		SessionId  string
		CircleName string
		Public     bool
	}{}
	r.DecodeJsonPayload(&payload)

	handle := payload.Handle
	sessionid := payload.SessionId
	circle_name := payload.CircleName
	is_public := payload.Public

	if !a.authenticate(handle, sessionid) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Failed to authenticate user request",
		})
		return
	}

	if circle_name == GOLD || circle_name == BROADCAST {
		w.WriteHeader(403)
		w.WriteJson(map[string]string{
			"Response": circle_name + " is a reserved circle name",
		})
		return
	}

	err := a.Svc.NewCircle(handle, circle_name, is_public)
	panicErr(err)

	w.WriteHeader(201)
	w.WriteJson(map[string]string{
		"Response": "Created new circle " + circle_name + " for " + handle,
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
	err := r.DecodeJsonPayload(&payload)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handle := payload.Handle
	sessionid := payload.SessionId
	content := payload.Content

	if !a.authenticate(handle, sessionid) {
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
	err = a.Svc.Db.Cypher(&neoism.CypherQuery{
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
		SessionId string
		LastSaved time.Time
		Circle    string
	}{}
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handle := payload.Handle
	sessionid := payload.SessionId
	lastsaved := payload.LastSaved
	circle := payload.Circle

	if !a.authenticate(handle, sessionid) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Failed to authenticate user request",
		})
		return
	}

	if !a.circleExists(handle, circle) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Bad Request, could not find specified circle to publish to",
		})
		return
	}

	if !a.messageExists(handle, lastsaved) {
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
	if _, ok := querymap["sessionid"]; !ok {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Bad Request, not enough parameters to authenticate user",
		})
		return
	}

	handle := querymap["handle"][0]
	sessionid := querymap["sessionid"][0]

	if !a.authenticate(handle, sessionid) {
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
	if _, ok := querymap["sessionid"]; !ok {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Bad Request, not enough parameters to authenticate user",
		})
		return
	}
	handle := querymap["handle"][0]
	sessionid := querymap["sessionid"][0]

	if !a.authenticate(handle, sessionid) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Failed to authenticate user request",
		})
		return
	}

	if !a.userExists(author) {
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
		SessionId string
		LastSaved time.Time
	}{}
	err := r.DecodeJsonPayload(&payload)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handle := payload.Handle
	sessionid := payload.SessionId
	lastsaved := payload.LastSaved

	if !a.authenticate(handle, sessionid) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Failed to authenticate user request",
		})
		return
	}

	deleted := []struct {
		Count int `json:"count(m)"`
	}{}
	a.Svc.Db.Cypher(&neoism.CypherQuery{
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
	})
	panicErr(err)

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
		Handle    string
		SessionId string
		Target    string
	}{}
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handle := payload.Handle
	sessionid := payload.SessionId
	target := payload.Target

	if !a.authenticate(handle, sessionid) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Failed to authenticate user request",
		})
		return
	}

	if !a.userExists(target) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Bad request, user " + target + " wasn't found",
		})
		return
	}

	a.Svc.RevokeMembershipBetween(handle, target)

	// Block user
	if block_occured, err := a.Svc.AddBlockedRelation(handle, target); err != nil {
		panicErr(err)
	} else if block_occured {
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
		Handle    string
		SessionId string
		Target    string
	}{}
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handle := payload.Handle
	sessionid := payload.SessionId
	target := payload.Target

	if !a.authenticate(handle, sessionid) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Failed to authenticate user request",
		})
		return
	}

	if !a.userExists(target) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Bad request, user " + target + " wasn't found",
		})
		return
	}

	if a.hasBlocked(handle, target) {
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
		Handle    string
		SessionId string
		Target    string
		Circle    string
	}{}
	if err := r.DecodeJsonPayload(&payload); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handle := payload.Handle
	sessionid := payload.SessionId
	target := payload.Target
	circle := payload.Circle

	if !a.authenticate(handle, sessionid) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Failed to authenticate user request",
		})
		return
	}

	if !a.userExists(target) {
		w.WriteHeader(400)
		w.WriteJson(map[string]string{
			"Response": "Bad request, user " + payload.Target + " wasn't found",
		})
		return
	}

	if a.hasBlocked(handle, target) {
		w.WriteHeader(403)
		w.WriteJson(map[string]string{
			"Response": "Server refusal to comply with join request",
		})
		return
	}

	if !a.circleExists(target, circle) {
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
