package service

import (
	"encoding/json"
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
// Return types
//
type Message struct {
	Id      string    `json:"m.id"`
	Author  string    `json:"u.handle"`
	Content string    `json:"m.content"`
	Created time.Time `json:"m.created"`
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

func (s Svc) CircleExistsInPublicDomain(circleid string) bool {
	found := []struct {
		Id string `json:"c.id"`
	}{}
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (c:Circle)-[:PART_OF]->(p:PublicDomain)
            WHERE c.id = {id}
            RETURN c.id
        `,
		Parameters: neoism.Props{
			"id": circleid,
		},
		Result: &found,
	}); err != nil {
		panicErr(err)
	}

	return len(found) > 0
}

func (s Svc) CanSeeCircle(fromPerspectiveOf string, circleid string) bool {
	if s.CircleExistsInPublicDomain(circleid) {
		return true
	}

	found := []struct {
		Id string `json:"c.id"`
	}{}
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
			MATCH   (u:User)-[:MEMBER_OF|CHIEF_OF]->(c:Circle)
			WHERE   u.handle = {handle}
			AND     c.id     = {id}
			RETURN  c.id
		`,
		Parameters: neoism.Props{
			"handle": fromPerspectiveOf,
			"id":     circleid,
		},
		Result: &found,
	}); err != nil {
		panicErr(err)
	}

	return len(found) > 0
}

func (s Svc) MessageExists(messageid string) bool {
	found := []struct {
		Id int `json:"m.id"`
	}{}
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (m:Message)
            WHERE m.id = {id}
            RETURN m.id
        `,
		Parameters: neoism.Props{
			"id": messageid,
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

func (s Svc) VerifySession(sessionid string) bool {
	found := []struct {
		Handle string `json:"u.handle"`
	}{}
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH  (u:User)<-[:SESSION_OF]-(a:AuthToken)
            WHERE  a.sessionid = {sessionid}
            AND    a.expires > {now}
            RETURN u.handle
        `,
		Parameters: neoism.Props{
			"sessionid": sessionid,
			"now":       time.Now().Local(),
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

func (s Svc) UserIsMemberOf(handle string, circleid string) bool {
	found := []struct {
		Handle string `json:"u.handle"`
	}{}
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
			MATCH (u:User)-[:MEMBER_OF:CHIEF_OF]->(c:Circle)
			WHERE u.handle = {handle}
			AND   c.id     = {circleid}
			RETURN u.handle
		`,
		Parameters: neoism.Props{
			"handle":   handle,
			"circleid": circleid,
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

func (s Svc) CreateNewUser(handle string, email string, password string) bool {
	newUser := []struct {
		Handle string    `json:"user.handle"`
		Email  string    `json:"user.email"`
		Joined time.Time `json:"user.joined"`
	}{}
	if err := s.Db.Cypher(&neoism.CypherQuery{
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
	}); err != nil {
		panicErr(err)
	}

	return len(newUser) > 0
}

func (s Svc) MakeDefaultCirclesFor(handle string) bool {
	created := []struct {
		Handle    string `json:"u.handle"`
		Gold      string `json:"g.name"`
		Broadcast string `json:"br.name"`
	}{}
	if err := s.Db.Cypher(&neoism.CypherQuery{
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
		Result: &created,
	}); err != nil {
		panicErr(err)
	}

	return len(created) > 0
}

func (s Svc) NewCircle(handle string, circle_name string, isPublic bool) bool {
	query := `
        MATCH   (u:User)
        WHERE   u.handle = {handle}
        CREATE  (u)-[:CHIEF_OF]->(c:Circle)
        SET     c.name = {name}
        SET     c.id = {id}
    `
	if isPublic {
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

	created := []struct {
		CircleName string `json:"c.name"`
	}{}
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: query,
		Parameters: neoism.Props{
			"handle": handle,
			"name":   circle_name,
			"id":     uniuri.NewLen(uniuri.UUIDLen),
		},
		Result: &created,
	}); err != nil {
		panicErr(err)
	}

	return len(created) > 0
}

func (s Svc) NewMessage(handle string, content string) bool {
	created := []struct {
		Content  string      `json:"m.content"`
		Relation neoism.Node `json:"r"`
	}{}
	now := time.Now().Local()
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH  (u:User)
            WHERE  u.handle = {handle}
            CREATE (m:Message {
                content:   {content}
              , created:   {now}
              , lastsaved: {now}
              , id:        {id}
            })
            CREATE (u)-[r:WROTE]->(m)
            RETURN m.content, r
        `,
		Parameters: neoism.Props{
			"handle":  handle,
			"content": content,
			"now":     now,
			"id":      uniuri.NewLen(uniuri.UUIDLen),
		},
		Result: &created,
	}); err != nil {
		panicErr(err)
	}
	return len(created) > 0
}

func (s Svc) PublishMessage(messageid, circleid string) bool {
	created := []struct {
		R *neoism.Relationship `json:"r"`
	}{}
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (m:Message)
            WHERE m.id={messageid}
            MATCH (c:Circle)
            WHERE c.id={circleid}
            CREATE (m)-[r:PUB_TO]->(c)
            SET r.published_at={now}
            RETURN r
        `,
		Parameters: neoism.Props{
			"messageid": messageid,
			"circleid":  circleid,
			"now":       time.Now().Local(),
		},
		Result: &created,
	}); err != nil {
		panicErr(err)
	}
	return len(created) > 0
}

func (s Svc) JoinCircle(handle string, circleid string) bool {
	joined := []struct {
		At time.Time `json:"r.at"`
	}{}
	now := time.Now().Local()
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH   (u:User), (c:Circle)
            WHERE   u.handle = {handle}
            AND     c.id     = {id}
            CREATE  (u)-[r:MEMBER_OF]->(c)
            SET     r.at     = {now}
            RETURN  r.at
        `,
		Parameters: neoism.Props{
			"handle": handle,
			"id":     circleid,
			"now":    now,
		},
		Result: &joined,
	}); err != nil {
		panicErr(err)
	}

	return len(joined) > 0
}

func (s Svc) JoinBroadcast(handle string, target string) bool {
	created := []struct {
		At time.Time `json:"r.at"`
	}{}
	now := time.Now().Local()
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH          (u:User)
            WHERE          u.handle = {handle}
            MATCH          (t:User)-[:CHIEF_OF]->(c:Circle)
            WHERE          t.handle = {target}
            AND            c.name = {broadcast}
            CREATE UNIQUE  (u)-[r:MEMBER_OF]->(c)
            SET            r.at = {now}
            RETURN         r.at
        `,
		Parameters: neoism.Props{
			"handle":    handle,
			"broadcast": BROADCAST,
			"target":    target,
			"now":       now,
		},
		Result: &created,
	}); err != nil {
		panicErr(err)
	}

	return len(created) > 0
}

func (s Svc) CreateBlockFromTo(handle string, target string) bool {
	res := []struct {
		Handle string      `json:"u.handle"`
		Target string      `json:"t.handle"`
		R      neoism.Node `json:"r"`
	}{}
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH (u:User), (t:User)
            WHERE u.handle = {handle}
            AND   t.handle = {target}
            CREATE UNIQUE (u)-[r:BLOCKED]->(t)
            RETURN u.handle, t.handle, r
        `,
		Parameters: neoism.Props{
			"handle": handle,
			"target": target,
		},
		Result: &res,
	}); err != nil {
		panicErr(err)
	}

	return len(res) > 0
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

func (s Svc) DeleteUser(handle string) bool {
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
                MATCH (u:User)
                WHERE u.handle = {handle}
                WITH  u
                OPTIONAL MATCH (a:AuthToken)-[r:SESSION_OF]->(u)
                DELETE a, r
                WITH u
                MATCH (u)-[wr:WROTE]->(m:Message)-[pt:PUB_TO]->(:Circle)
                DELETE pt, m, wr
                WITH u
                MATCH (u)-[mo:MEMBER_OF]->(:Circle)
                DELETE mo
                WITH u
                MATCH (u)-[b:BLOCKED]->(:User)
                DELETE b
                WITH u
                MATCH (u)-[co:CHIEF_OF]->(c:Circle)-[po:PART_OF]->(:PublicDomain)
                MATCH (c)<-[mo:MEMBER_OF]-(:User)
                MATCH (c)<-[pt:PUB_TO]-(:Message)
                DELETE pt, mo, co, po, c, u
            `,
		Parameters: neoism.Props{
			"handle": handle,
		},
	}); err != nil {
		panicErr(err)
	}
	return true
}

//
// Get
//

func (s Svc) SearchForUsers(
	circle string,
	nameprefix string,
	skip int,
	limit int,
	sort string,
) (results string, count int) {
	res := []struct {
		Handle string `json:"u.handle"`
		Name   string `json:"u.name"`
		Id     int    `json:"id(u)"`
	}{}

	var statement string
	props := neoism.Props{}

	regex := "(?i)" + nameprefix + ".*"

	if circle != "" {
		statement = `
			MATCH (u:User)-[]->(c:Circle)
			WHERE c.name = {circle}
			AND   u.handle =~ {regex}
		`
		props = neoism.Props{
			"circle": circle,
			"regex":  regex,
			"skip":   skip,
			"limit":  limit,
			"sort":   sort,
		}
	} else {
		statement = `
			MATCH  (u:User)
			WHERE  u.handle =~ {regex}
		`
		props = neoism.Props{
			"regex": regex,
			"skip":  skip,
			"limit": limit,
			"sort":  sort,
		}
	}

	statement = statement + `
        RETURN u.handle, u.name, id(u)
        SKIP   {skip}
        LIMIT  {limit}
	`

	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement:  statement,
		Parameters: props,
		Result:     &res,
	}); err != nil {
		panicErr(err)
	} else if len(res) == 0 {
		return "", 0
	}

	bytes, err := json.Marshal(res)
	if err != nil {
		panicErr(err)
	}

	return string(bytes), len(res)
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

func (s Svc) GetCircleId(handle string, circle string) string {
	found := []struct {
		Id string `json:"c.id"`
	}{}
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
			MATCH  (u:User)-[:CHIEF_OF]->(c:Circle)
			WHERE  u.handle = {handle}
			AND    c.name   = {circle}
			RETURN c.id
		`,
		Parameters: neoism.Props{
			"handle": handle,
			"circle": circle,
		},
		Result: &found,
	}); err != nil {
		panicErr(err)
	}

	if len(found) > 0 {
		return found[0].Id
	} else {
		return ""
	}
}

func (s Svc) GetMessagesByHandle(handle string) []Message {
	messages := make([]Message, 0)
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH     (u:User)-[:WROTE]->(m:Message)
            WHERE     u.handle = {handle}
            RETURN    m.id
                    , u.handle
                    , m.content
                    , m.created
            ORDER BY  m.created
        `,
		Parameters: neoism.Props{
			"handle": handle,
		},
		Result: &messages,
	}); err != nil {
		panicErr(err)
	}

	return messages
}

func (s Svc) GetHandleFromAuthorization(token string) (string, bool) {
	found := []struct {
		Handle string `json:"u.handle"`
	}{}
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
			MATCH   (u:User)<-[:SESSION_OF]-(a:AuthToken)
			WHERE   a.sessionid = {sessionid}
			AND     {now} < a.expires
			RETURN  u.handle
		`,
		Parameters: neoism.Props{
			"sessionid": token,
			"now":       time.Now().Local(),
		},
		Result: &found,
	}); err != nil {
		panicErr(err)
	}

	if success := len(found) > 0; success {
		return found[0].Handle, success
	} else {
		return "", success
	}
}

//
// Node Attributes
//

// Sets a session id on an AuthToken node that points to a particular user
func (s Svc) SetGetNewSessionId(handle string) string {
	created := []struct {
		SessionId string `json:"a.sessionid"`
	}{}

	sessionDuration := time.Hour
	now := time.Now().Local()

	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
                MATCH  (u:User)
                WHERE  u.handle = {handle}
                WITH   u
                OPTIONAL MATCH (u)<-[s:SESSION_OF]-(a:AuthToken)
                DELETE s, a
                WITH   u
                CREATE (u)<-[r:SESSION_OF]-(a:AuthToken)
                SET    r.created_at = {now}
                SET    a.sessionid  = {sessionid}
                SET    a.expires    = {time}
                RETURN a.sessionid
            `,
		Parameters: neoism.Props{
			"handle":    handle,
			"sessionid": uniuri.NewLen(uniuri.UUIDLen),
			"time":      now.Add(sessionDuration),
			"now":       now,
		},
		Result: &created,
	}); err != nil {
		panicErr(err)
	}
	if len(created) != 1 {
		panic(fmt.Sprintf("Incorrect results len in query1()\n\tgot %d, expected 1\n", len(created)))
	}

	return created[0].SessionId
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

func (s Svc) UnsetSessionId(handle string) bool {
	unset := []struct {
		Handle string `json:"u.handle"`
	}{}
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH          (u:User)
            WHERE          u.handle = {handle}
            WITH           u
            OPTIONAL MATCH (u)<-[so:SESSION_OF]-(a:AuthToken)
            DELETE         so, a
            RETURN         u.handle
        `,
		Parameters: neoism.Props{
			"handle": handle,
		},
		Result: &unset,
	}); err != nil {
		panicErr(err)
	}
	return len(unset) > 0
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
