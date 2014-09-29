package api

import (
	"errors"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/dchest/uniuri"
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
	Db *neoism.Database
}

// Circle constants
const (
	GOLD   = "Gold"
	PUBLIC = "Public"
)

//
// API util
//

// Errors if authentication is not successful
func (a Api) authenticate(w rest.ResponseWriter, handle string, sessionid string) {
	found := []struct {
		Handle string `json:"user.handle"`
	}{}
	dberr := a.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (user:User {handle:{handle}, sessionid:{sessionid}})
            RETURN user.handle
        `,
		Parameters: neoism.Props{
			"handle":    handle,
			"sessionid": sessionid,
		},
		Result: &found,
	})
	panicErr(dberr)

	if len(found) == 0 {
		rest.Error(w, "Could not authenticate user "+handle, 400)
	}
}

func (a Api) userExists(handle string) bool {
	found := []struct {
		Handle string `json"user.handle"`
	}{}
	err := a.Db.Cypher(&neoism.CypherQuery{
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

func (a Api) circleExists(handle string, circleName string) bool {
	found := []struct {
		Name string `json"circle.name"`
	}{}
	err := a.Db.Cypher(&neoism.CypherQuery{
		Statement: `
			MATCH (user:User {handle: {handle}})
			OPTIONAL MATCH (circle:Circle {name: {name}})
			RETURN circle.name
		`,
		Parameters: neoism.Props{
			"handle": handle,
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
	err := a.Db.Cypher(&neoism.CypherQuery{
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

//
// API
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
	err := r.DecodeJsonPayload(&proposal)
	fmt.Println("email: " + proposal.Email)
	fmt.Printf("err was nil: %t\n", err == nil)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("%+v\n", r.Body)
	fmt.Println("email: " + proposal.Email)

	// Handle and Email checks
	if proposal.Handle == "" {
		rest.Error(w, "Handle is a required field for signup", 400)
		return
	} else if proposal.Email == "" {
		rest.Error(w, "Email is a required field for signup", 400)
		return
	}

	// Password checks
	minPasswordLength := 8
	if proposal.Password != proposal.ConfirmPassword {
		rest.Error(w, "Passwords do not match", 400)
		return
	} else if len(proposal.Password) < minPasswordLength {
		rest.Error(w, "Passwords must be at least 8 characters long", 400)
		return
	}

	// Ensure unique handle
	foundUsers := []struct {
		Handle string `json:"user.handle"`
	}{}
	err = a.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (user:User {handle: {handle}})
            RETURN user.handle
        `,
		Parameters: neoism.Props{
			"handle": proposal.Handle,
		},
		Result: &foundUsers,
	})
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(foundUsers) > 0 {
		rest.Error(w, "Sorry, "+proposal.Handle+" is already taken", 400)
		return
	}

	// Ensure unique email
	foundEmails := []struct {
		Email string `json:"user.email"`
	}{}
	err = a.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (user:User {email: {email}})
            RETURN user.email
        `,
		Parameters: neoism.Props{
			"email": proposal.Email,
		},
		Result: &foundEmails,
	})
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(foundEmails) > 0 {
		rest.Error(w, "Sorry, "+proposal.Email+" is already taken", 400)
		return
	}

	newUser := []struct {
		Handle string    `json:"user.handle"`
		Email  string    `json:"user.email"`
		Joined time.Time `json:"user.joined"`
	}{}
	err = a.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            CREATE (user:User { handle: {handle}, email: {email}, password: {password}, joined: {joined} })
            RETURN user.handle, user.email, user.joined
        `,
		Parameters: neoism.Props{
			"handle":   proposal.Handle,
			"email":    proposal.Email,
			"password": proposal.Password,
			"joined":   time.Now().Local(),
		},
		Result: &newUser,
	})
	panicErr(err)

	if len(newUser) != 1 {
		panic(fmt.Sprintf("Incorrect results len in query1()\n\tgot %d, expected 1\n", len(newUser)))
	}

	// Add 'Gold' and 'Public' circles
	a.makeDefaultCircles(proposal.Handle)

	var handle string = newUser[0].Handle
	var email string = newUser[0].Email
	var joined string = newUser[0].Joined.Format(time.RFC1123)

	w.WriteHeader(201)
	w.WriteJson(map[string]string{
		"Response": "Signed up a new user!",
		"Handle":   handle,
		"Email":    email,
		"Joined":   joined,
	})
}

func (a Api) Login(w rest.ResponseWriter, r *rest.Request) {
	credentials := struct {
		Handle   string
		Password string
	}{}
	err := r.DecodeJsonPayload(&credentials)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	found := []struct {
		Handle string `json:"user.handle"`
	}{}
	err = a.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (user:User {handle:{handle}, password:{password}})
            RETURN user.handle
        `,
		Parameters: neoism.Props{
			"handle":   credentials.Handle,
			"password": credentials.Password,
		},
		Result: &found,
	})
	panicErr(err)

	// Add session hash (16 chars) to node and return it to client in json res
	if len(found) > 0 {
		sessionHash := uniuri.New()

		idResponse := []struct {
			SessionId string `json:"user.sessionid"`
		}{}
		a.Db.Cypher(&neoism.CypherQuery{
			Statement: `
                MATCH (user:User {handle:{handle}, password:{password}})
                SET user.sessionid = {sessionid}
                return user.sessionid
            `,
			Parameters: neoism.Props{
				"handle":    credentials.Handle,
				"password":  credentials.Password,
				"sessionid": sessionHash,
			},
			Result: &idResponse,
		})
		if len(idResponse) != 1 {
			panic(fmt.Sprintf("Incorrect results len in query1()\n\tgot %d, expected 1\n", len(idResponse)))
		}

		w.WriteJson(map[string]string{
			"Response":  "Logged in " + credentials.Handle + ". Note your session id.",
			"SessionId": sessionHash,
		})
		return
	} else {
		rest.Error(w, "Invaid username or password, please try again.", 400)
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

	loggedOut := []struct {
		Handle string `json:"user.handle"`
	}{}
	a.Db.Cypher(&neoism.CypherQuery{
		Statement: `
			MATCH (user:User {handle: {handle}})
			REMOVE user.sessionid
			RETURN user
		`,
		Parameters: neoism.Props{
			"handle": user.Handle,
		},
		Result: &loggedOut,
	})

	if len(loggedOut) == 0 {
		w.WriteHeader(403)
		w.WriteJson(map[string]string{
			"Response": "No user was logged out",
		})
		return
	}

	w.WriteHeader(200)
	w.WriteJson(map[string]string{
		"Response": "Logged out " + user.Handle + ", have a nice day",
	})
}

func (a Api) GetUser(w rest.ResponseWriter, r *rest.Request) {
	querymap := r.URL.Query()

	// Get by handle
	if handle, ok := querymap["handle"]; ok {
		stmt := `MATCH (user:User)
                 WHERE user.handle = {handle}
                 RETURN user`
		params := neoism.Props{
			"handle": handle[0],
		}
		res := []struct {
			User neoism.Node
		}{}

		err := a.Db.Cypher(&neoism.CypherQuery{
			Statement:  stmt,
			Parameters: params,
			Result:     &res,
		})
		panicErr(err)

		u := res[0].User.Data

		w.WriteJson(u)
		return
	}

	// All users
	stmt := `MATCH (user:User)
             RETURN user.handle, user.joined
             ORDER BY user.handle`
	res := []struct {
		Handle string    `json:"user.handle"`
		Joined time.Time `json:"user.joined"`
	}{}

	err := a.Db.Cypher(&neoism.CypherQuery{
		Statement:  stmt,
		Parameters: neoism.Props{},
		Result:     &res,
	})
	panicErr(err)

	if len(res) > 0 {
		w.WriteJson(res)
	} else {
		w.WriteJson(map[string]string{
			"Response": "No results found",
		})
	}
}

func (a Api) DeleteUser(w rest.ResponseWriter, r *rest.Request) {
	querymap := r.URL.Query()

	if handle, ok := querymap["handle"]; ok {
		if password, okok := querymap["password"]; okok {

			var handle = handle[0]
			var password = password[0]

			res := []struct {
				HandleToBeDeleted string `json:"user.handle"`
			}{}
			err := a.Db.Cypher(&neoism.CypherQuery{
				Statement: `
                    MATCH (user:User {handle:{handle}, password:{password}})
                    RETURN user.handle
                `,
				Parameters: neoism.Props{
					"handle":   handle,
					"password": password,
				},
				Result: &res,
			})
			panicErr(err)

			if len(res) > 0 {
				err := a.Db.Cypher(&neoism.CypherQuery{
					// Delete user node
					Statement: `
                        MATCH (u:User {handle: {handle}})
                        DELETE u
                    `,
					Parameters: neoism.Props{
						"handle": handle,
					},
					Result: nil,
				})
				panicErr(err)

				w.WriteJson(map[string]string{
					"Response": "Deleted " + handle,
				})
				return
			} else {
				w.WriteHeader(403)
				w.WriteJson(map[string]string{
					"Response": "Could not delete user with supplied credentials",
				})
				return
			}
		}
	}
	w.WriteHeader(403)
	w.WriteJson(map[string]string{
		"Error": "Bad request parameters for delete, expected handle:String, password:String",
	})
}

func (a Api) makeCircleForUser(handle string, circleName string) (err error) {
	if circleName == GOLD || circleName == PUBLIC {
		return errors.New(circleName + " is a reserved circle name")
	}

	made := []struct {
		Handle string `json:"user.name"`
		Name   string `json:"circle.name"`
	}{}
	dberr := a.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (user:User)
            WHERE user.handle = {handle}
            CREATE (circle:Circle {name: {name}})
            CREATE (user)-[:CHIEF_OF]->(circle)
            RETURN user.name, circle.name
        `,
		Parameters: neoism.Props{
			"handle": handle,
			"name":   circleName,
		},
		Result: &made,
	})
	panicErr(dberr)
	if len(made) != 1 {
		panic(fmt.Sprintf("Incorrect results len in query1()\n\tgot %d, expected 1\n", len(made)))
	}

	return nil
}

func (a Api) makeDefaultCircles(handle string) {
	made := []struct {
		Handle string `json:"user.handle"`
		G      string `json:"g.name"`
		P      string `json:"p.name"`
	}{}
	dberr := a.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (user:User)
            WHERE user.handle = {handle}
            CREATE (g:Circle {name: {gold}})
            CREATE (p:Circle {name: {public}})
            CREATE (user)-[:CHIEF_OF]->(g)
            CREATE (user)-[:CHIEF_OF]->(p)
            RETURN user.handle, g.name, p.name
        `,
		Parameters: neoism.Props{
			"handle": handle,
			"gold":   GOLD,
			"public": PUBLIC,
		},
		Result: &made,
	})
	panicErr(dberr)
	if len(made) != 1 {
		panic(fmt.Sprintf("Incorrect results len in query1()\n\tgot %d, expected 1\n", len(made)))
	}
}

/**
 * Expects a json post with "handle", "sessionid", "content"
 */
func (a Api) NewMessage(w rest.ResponseWriter, r *rest.Request) {
	payload := struct {
		Handle    string
		Sessionid string
		Content   string
	}{}
	err := r.DecodeJsonPayload(&payload)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.authenticate(w, payload.Handle, payload.Sessionid)

	if payload.Content == "" {
		rest.Error(w, "Please enter some content for your message", 400)
		return
	}

	created := []struct {
		Content  string      `json:"message.content"`
		Relation neoism.Node `json:"r"`
	}{}
	createdAt := time.Now().Local()
	err = a.Db.Cypher(&neoism.CypherQuery{
		Statement: `
			MATCH (user:User {handle: {handle}, sessionid: {sessionid}})
			CREATE (message:Message {content: {content}, created: {date}, lastsaved: {date}})
			CREATE (user)-[r:WROTE]->(message)
			RETURN message.content, r
		`,
		Parameters: neoism.Props{
			"handle":    payload.Handle,
			"sessionid": payload.Sessionid,
			"content":   payload.Content,
			"date":      createdAt,
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
			"Response":  "Successfully created message for " + payload.Handle,
			"Published": false,
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

	// Unmarshall
	payload := struct {
		Handle    string
		Sessionid string
	}{
		querymap["handle"][0],
		querymap["sessionid"][0],
	}

	a.authenticate(w, payload.Handle, payload.Sessionid)

	messages := []struct {
		Content   string    `json:"message.content"`
		LastSaved time.Time `json:"message.lastsaved"`
	}{}
	err := a.Db.Cypher(&neoism.CypherQuery{
		Statement: `
			MATCH (user:User {handle: {handle}})-[:WROTE]->(message:Message)
			RETURN message.content, message.lastsaved	
		`,
		Parameters: neoism.Props{
			"handle": payload.Handle,
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
	author := r.PathParam("handle")
	querymap := r.URL.query()

	a.authenticate(w, handle, sessionid)

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

	if !userExists(author) {
		w.WriteHeader(404)
		w.WriteJson(map[string]string{
			"Response": "Bad request, user doesn't exist",
		})
		return
	}

	messages := []struct {
		Content   string    `json:"message.content"`
		Published time.Time `json:"message.published"`
	}{}
	err := a.Db.Cypher(&neoism.CypherQuery{
		Statement: `
			MATCH (author:User {handle: {author}}), (user:User {handle: {handle}})
			OPTIONAL MATCH (user)-[r:MEMBER_OF]->(circle:Circle)
			OPTIONAL MATCH (author)-[w:WROTE]-(visible:Message)-[p:PUB_TO]->(circle)
			RETURN visible.content, visible.published_at
		`,
		Parameters: neoism.Props{
			"author": author,
			"handle": querymap["handle"][0],
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
		Sessionid string
		Lastsaved time.Time
	}{}
	r.DecodeJsonPayload(&payload)

	a.authenticate(w, payload.Handle, payload.Sessionid)

	deleted := []struct {
		Count int `json:"count(m)"`
	}{}
	err := a.Db.Cypher(&neoism.CypherQuery{
		Statement: `
		MATCH (user:User {handle: {handle}})
		OPTIONAL MATCH (user)-[r:WROTE]->(m:Message {lastsaved: {lastsaved}})
		DELETE r, m
		RETURN count(m)
		`,
		Parameters: neoism.Props{
			"handle":    payload.Handle,
			"lastsaved": payload.Lastsaved,
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
