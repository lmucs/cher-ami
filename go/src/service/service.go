package service

import (
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/jmcvetta/neoism"
	"log"
	"time"
)

//
// Constants
//

const (
	// Reserved Circles
	GOLD      = "Gold"
	BROADCAST = "Broadcast"
)

//
// Service Types
//

type Svc struct {
	Db *neoism.Database
}

//
// Utility Functions
//

/**
 * Service instances must be initialized using this method in
 * order to ensure data integrity. Do not instantiate Svc directly.
 */
func NewService(uri string) *Svc {
	neo4jdb, err := neoism.Connect(uri)
	if err != nil {
		log.Fatal(err)
	}
	s := &Svc{neo4jdb}
	s.databaseInit()
	return s
}

func (s Svc) databaseInit() {
	var publicdomain *neoism.Node
	// Initialize PublicDomain node
	// Nodes must have at least one property to allow unique creation
	publicdomain, _, err := s.Db.GetOrCreateNode("PublicDomain", "iam", neoism.Props{
		"iam": "PublicDomain",
	})
	panicErr(err)
	// Label (has to be) added separately
	err = publicdomain.AddLabel("PublicDomain")
	panicErr(err)

	if publicdomain != nil {
		fmt.Println("Public Domain available")
	} else {
		fmt.Println("Unexpected database state, possible lack of PublicDomain")
	}
}

//
// Checks
//

func (s Svc) UserExists(handle string) bool {
	found := []struct {
		Handle string `json:"u.handle"`
	}{}
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (u:User)
            WHERE u.handle = {handle}
            RETURN u.handle
        `,
		Parameters: neoism.Props{
			"handle": handle,
		},
		Result: &found,
	}); err != nil {
		panicErr(err)
	}

	return len(found) > 0
}

func (s Svc) CircleExists(target string, circleName string) bool {
	found := []struct {
		Name string `json:"c.name"`
	}{}
	if err := s.Db.Cypher(&neoism.CypherQuery{
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
	}); err != nil {
		panicErr(err)
	}

	return len(found) > 0
}

func (s Svc) MessageExists(handle string, lastSaved time.Time) bool {
	found := []struct {
		Id int `json:"id(m)"`
	}{}
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (u:User)
            WHERE u.handle = {handle}
            MATCH (u)-[:WROTE]->(m:Message)
            WHERE m.lastsaved = {lastsaved}
            RETURN id(m)
        `,
		Parameters: neoism.Props{
			"handle":    handle,
			"lastsaved": lastSaved,
		},
		Result: &found,
	}); err != nil {
		panicErr(err)
	}

	return len(found) > 0
}

func (s Svc) HandleIsUnique(handle string) (bool, error) {
	found := []struct {
		Handle string `json:"u.handle"`
	}{}
	err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (u:User {handle: {handle}})
            RETURN u.handle
        `,
		Parameters: neoism.Props{
			"handle": handle,
		},
		Result: &found,
	})

	return len(found) == 0, err
}

func (s Svc) EmailIsUnique(email string) (bool, error) {
	found := []struct {
		Email string `json:"u.email"`
	}{}
	err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (u:User {email: {email}})
            RETURN u.email
        `,
		Parameters: neoism.Props{
			"email": email,
		},
		Result: &found,
	})

	return len(found) == 0, err
}

func (s Svc) GoodSessionCredentials(handle string, sessionid string) bool {
	found := []struct {
		Handle string `json:"u.handle"`
	}{}
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (u:User)
            WHERE u.handle = {handle}
            AND   u.sessionid = {sessionid}
            RETURN u.handle
        `,
		Parameters: neoism.Props{
			"handle":    handle,
			"sessionid": sessionid,
		},
		Result: &found,
	}); err != nil {
		panicErr(err)
	}

	return len(found) > 0
}

func (s Svc) GoodLoginCredentials(handle string, password string) (bool, error) {
	found := []struct {
		Handle string `json:"user.handle"`
	}{}
	err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (user:User {handle:{handle}, password:{password}})
            RETURN user.handle
        `,
		Parameters: neoism.Props{
			"handle":   handle,
			"password": password,
		},
		Result: &found,
	})

	return len(found) == 1, err
}

func (s Svc) BlockExistsFromTo(handle string, target string) bool {
	found := []struct {
		Relation int `json:"r"`
	}{}
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (u:User), (t:User)
            WHERE u.handle = {handle}
            AND   t.handle = {target}
            MATCH (u)-[r:BLOCKED]->(t)
            RETURN r
        `,
		Parameters: neoism.Props{
			"handle": handle,
			"target": target,
		},
		Result: &found,
	}); err != nil {
		panicErr(err)
	}

	return len(found) > 0
}

//
// Creation
//

func (s Svc) CreateNewUser(handle string, email string, password string) error {
	newUser := []struct {
		Handle string    `json:"user.handle"`
		Email  string    `json:"user.email"`
		Joined time.Time `json:"user.joined"`
	}{}
	err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            CREATE (user:User {
                handle:   {handle},
                name:     "I AM A NAME!!!!!!",
                email:    {email},
                password: {password},
                joined:   {joined}
            })
            RETURN user.handle, user.email, user.joined
        `,
		Parameters: neoism.Props{
			"handle":   handle,
			"email":    email,
			"password": password,
			"joined":   time.Now().Local(),
		},
		Result: &newUser,
	})
	return err
}

func (s Svc) MakeDefaultCirclesFor(handle string) error {
	made := []struct {
		Handle    string `json:"u.handle"`
		Gold      string `json:"g.name"`
		Broadcast string `json:"br.name"`
	}{}
	err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (p:PublicDomain)
            WHERE p.iam = "PublicDomain"
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
	if len(made) != 1 {
		panic(fmt.Sprintf("Incorrect results len in query1()\n\tgot %d, expected 1\n", len(made)))
	}
	return err
}

func (s Svc) NewCircle(handle string, circle_name string, is_public bool) error {
	query := `
        MATCH (u:User)
        WHERE u.handle = {handle}
        CREATE (u)-[:CHIEF_OF]->(c:Circle)
        SET c.name = {name}
    `
	if is_public {
		query = query + `
            WITH u, c
            MATCH (p:PublicDomain)
            WHERE p.iam = "PublicDomain"
            CREATE (c)-[:PART_OF]->(p)
        `
	}
	query = query + `
        RETURN c.name
    `

	made := []struct {
		CircleName string `json:"c.name"`
	}{}
	err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: query,
		Parameters: neoism.Props{
			"handle": handle,
			"name":   circle_name,
		},
		Result: &made,
	})
	if err != nil {
		panicErr(err)
	}
	if len(made) != 1 {
		panic(fmt.Sprintf("Incorrect results len in query1()\n\tgot %d, expected 1\n", len(made)))
	}

	return err
}

func (s Svc) JoinCircle(handle string, target string, target_circle string) (at time.Time, did_join bool) {
	joined := []struct {
		At time.Time `json:"r.at"`
	}{}
	now := time.Now().Local()
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (u:User)
            WHERE u.handle = {handle}
            MATCH (t:User)-[:CHIEF_OF]->(c:Circle)
            WHERE t.handle = {target} AND c.name = {circle}
            CREATE (u)-[r:MEMBER_OF {at: {now}}]->(c)
            RETURN r.at
        `,
		Parameters: neoism.Props{
			"handle": handle,
			"target": target,
			"circle": target_circle,
			"now":    now,
		},
		Result: &joined,
	}); err != nil {
		panicErr(err)
	} else if len(joined) != 1 {
		return time.Time{}, false
	}

	return joined[0].At, len(joined) > 0
}

func (s Svc) JoinBroadcast(handle string, target string) (at time.Time, did_join bool) {
	joined := []struct {
		Target string    `json:"t.handle"`
		At     time.Time `json:"r.at"`
	}{}
	now := time.Now().Local()
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (u:User)
            WHERE u.handle = {handle}
            MATCH (t:User)-[:CHIEF_OF]->(c:Circle)
            WHERE t.handle = {target} AND c.name = {broadcast}
            CREATE UNIQUE (u)-[r:MEMBER_OF]->(c)
            SET r.at = {now}
            RETURN r.at
        `,
		Parameters: neoism.Props{
			"handle":    handle,
			"broadcast": BROADCAST,
			"target":    target,
			"now":       now,
		},
		Result: &joined,
	}); err != nil {
		panicErr(err)
	} else if len(joined) != 1 {
		return time.Time{}, false
	}

	return now, true
}

func (s Svc) AddBlockedRelation(handle string, target string) (block_occured bool, err error) {
	res := []struct {
		Handle string      `json:"u.handle"`
		Target string      `json:"t.handle"`
		R      neoism.Node `json:"r"`
	}{}
	err = s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (u:User)
            WHERE u.handle = {handle}
            MATCH (t:User)
            WHERE t.handle = {target}
            CREATE UNIQUE (t)-[r:BLOCKED]->(u)
            RETURN u.handle, t.handle, r
        `,
		Parameters: neoism.Props{
			"handle": handle,
			"target": target,
		},
		Result: &res,
	})

	return len(res) > 0, err
}

//
// Deletion
//

func (s Svc) DeleteAllNodesAndRelations() {
	s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (n)
            OPTIONAL MATCH (n)-[r]-()
            DELETE n, r
        `,
	})
}

func (s Svc) FreshInitialState() {
	s.DeleteAllNodesAndRelations()
	s.databaseInit()
}

func (s Svc) RevokeMembershipBetween(handle string, target string) {
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (u:User)
            WHERE u.handle={handle}
            MATCH (t:User)
            WHERE t.handle={target}
            OPTIONAL MATCH (u)-[:CHIEF_OF]->(c:Circle)
            OPTIONAL MATCH (t)-[r:MEMBER_OF]->(c)
            DELETE r
        `,
		Parameters: neoism.Props{
			"handle": handle,
			"target": target,
		},
	}); err != nil {
		panicErr(err)
	}
}

func (s Svc) DeleteUser(handle string) error {
	return s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
                MATCH (user:User)-[r]->()
                WHERE user.handle = {handle}
                DELETE user, r
            `,
		Parameters: neoism.Props{
			"handle": handle,
		},
	})
}

//
// Get
//

func (s Svc) GetHandleAndNameOf(user string) (handle string, name string, found bool) {
	res := []struct {
		Handle string `json:"u.handle"`
		Name   string `json:"u.name"`
	}{}
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (u:User)
            WHERE u.handle = {handle}
            RETURN u.handle, u.name
        `,
		Parameters: neoism.Props{
			"handle": user,
		},
		Result: &res,
	}); err != nil {
		panicErr(err)
	} else if len(res) != 1 {
		return "", "", len(res) > 0
	}

	return res[0].Handle, res[0].Name, len(res) > 0
}

func (s Svc) GetPasswordHash(user string) (password_hash []byte, found bool) {
	res := []struct {
		PasswordHash string `json:"u.password"`
	}{}
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (u:User)
            WHERE u.handle = {handle}
            RETURN u.password
        `,
		Parameters: neoism.Props{
			"handle": user,
		},
		Result: &res,
	}); err != nil {
		panicErr(err)
	} else if len(res) != 1 {
		return []byte{}, len(res) > 0
	}

	return []byte(res[0].PasswordHash), len(res) > 0
}

//
// Node Attributes
//

func (s Svc) SetGetNewSessionId(handle string) (sessionid string, err error) {
	sessionHash := uniuri.New()

	created := []struct {
		SessionId string `json:"u.sessionid"`
	}{}
	err = s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
                MATCH (u:User)
                WHERE u.handle = {handle}
                SET u.sessionid = {sessionid}
                return u.sessionid
            `,
		Parameters: neoism.Props{
			"handle":    handle,
			"sessionid": sessionHash,
		},
		Result: &created,
	})
	if len(created) != 1 {
		panic(fmt.Sprintf("Incorrect results len in query1()\n\tgot %d, expected 1\n", len(created)))
	}

	return created[0].SessionId, err
}

func (s Svc) SetNewPassword(handle string, password string) bool {
	user := []struct {
		Password string
	}{}
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (u:User)
            WHERE u.handle = {handle}
            SET u.password = {password}
            RETURN u.password
        `,
		Parameters: neoism.Props{
			"handle":   handle,
			"password": password,
		},
		Result: &user,
	}); err != nil {
		panicErr(err)
	} else if len(user) != 1 {
		panic(fmt.Sprintf("Incorrect results len in query1()\n\tgot %d, expected 1\n", len(user)))
	}

	return len(user) > 0
}

func (s Svc) UnsetSessionId(handle string) error {
	err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (u:User)
            WHERE u.handle = {handle}
            REMOVE u.sessionid
        `,
		Parameters: neoism.Props{
			"handle": handle,
		},
		Result: nil,
	})

	return err
}

func (s Svc) SetGetName(handle string, name string) string {
	user := []struct {
		Name string
	}{}
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (u:User)
            WHERE u.handle = {handle}
            SET u.name = {name}
            RETURN u.name
        `,
		Parameters: neoism.Props{
			"handle": handle,
			"name":   name,
		},
		Result: &user,
	}); err != nil {
		panicErr(err)
	} else if len(user) != 1 {
		panic(fmt.Sprintf("Incorrect results len in query1()\n\tgot %d, expected 1\n", len(user)))
	}

	return user[0].Name
}

//
// Errors
//

func panicErr(err error) {
	if err != nil {
		panic(err)
		return
	}
}
