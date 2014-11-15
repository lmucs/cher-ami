package query

import (
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/jmcvetta/neoism"
	"time"
)

type Query struct {
	Db *neoism.Database
}

//
// Initialization
//

// Constructor, use this when creating a new Query struct.
func NewQuery(uri string) *Query {
	neo4jdb, err := neoism.Connect(uri)
	if err != nil {
		panicIfErr(err)
	}
	query := Query{neo4jdb}
	query.databaseInit()
	return &query
}

// Initializes the Neo4j Database
func (q Query) databaseInit() {
	if publicDomain := q.CreateUniquePublicDomain(); publicDomain == nil {
		fmt.Println("Unexpected database state, possible lack of PublicDomain")
	}
}

//
// Private Utilities
//

// Preforms a Cypher query, catching any unexpected behavior in a panic.
// It is ok to panic in this case as a panic at the db query level almost
// always indicates an incorrectly constructed query.
func (q Query) cypherOrPanic(query *neoism.CypherQuery) {
	panicIfErr(q.Db.Cypher(query))
}

// Asserts that err is non-nil then panics if so
func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Constants
const (
	// Reserved Circles
	GOLD      = "Gold"
	BROADCAST = "Broadcast"
)

//
// Create
//

func (q Query) CreateUniquePublicDomain() *neoism.Node {
	// Initialize PublicDomain node
	// Nodes must have at least one property to allow unique creation
	if publicDomain, _, err := q.Db.GetOrCreateNode("PublicDomain", "iam", neoism.Props{
		"iam": "PublicDomain",
	}); err != nil {
		panic(err)
	} else {
		// Label (has to be) added separately
		panicIfErr(publicDomain.AddLabel("PublicDomain"))

		return publicDomain
	}
}

func (q Query) NewUser(handle, email, passwordHash string) bool {
	newUser := []struct {
		Handle string    `json:"user.handle"`
		Email  string    `json:"user.email"`
		Joined time.Time `json:"user.joined"`
	}{}
	q.cypherOrPanic(&neoism.CypherQuery{
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
			"password": passwordHash,
			"joined":   time.Now().Local(),
		},
		Result: &newUser,
	})
	return len(newUser) > 0
}

func (q Query) CreateDefaultCirclesForUser(handle string) bool {
	created := []struct {
		Handle    string `json:"u.handle"`
		Gold      string `json:"g.name"`
		Broadcast string `json:"br.name"`
	}{}
	q.cypherOrPanic(&neoism.CypherQuery{
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
	})
	return len(created) > 0
}

func (q Query) CreateCircle(handle, circleName string, isPublic bool,
) (circleid string, ok bool) {
	created := []struct {
		CircleName string `json:"c.name"`
		CircleId   string `json:"c.id"`
	}{}

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
        RETURN c.name, c.id
    `

	q.cypherOrPanic(&neoism.CypherQuery{
		Statement: query,
		Parameters: neoism.Props{
			"handle": handle,
			"name":   circleName,
			"id":     uniuri.NewLen(uniuri.UUIDLen),
		},
		Result: &created,
	})

	if ok = len(created) > 0; ok {
		return created[0].CircleId, ok
	} else {
		return "", ok
	}
}

func (q Query) CreateMessage(handle, content string) (messageid string, ok bool) {
	created := []struct {
		Content string `json:"m.content"`
		Id      string `json:"m.id"`
	}{}
	q.cypherOrPanic(&neoism.CypherQuery{
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
            RETURN m.content, m.id
        `,
		Parameters: neoism.Props{
			"handle":  handle,
			"content": content,
			"now":     time.Now().Local(),
			"id":      uniuri.NewLen(uniuri.UUIDLen),
		},
		Result: &created,
	})

	if ok = len(created) > 0; ok {
		return created[0].Id, ok
	} else {
		return "", ok
	}
}

func (q Query) CreatePublishedRelation(messageid, circleid string) bool {
	created := []struct {
		R *neoism.Relationship `json:"r"`
	}{}
	q.cypherOrPanic(&neoism.CypherQuery{
		Statement: `
            MATCH   (m:Message), (c:Circle)
            WHERE   m.id = {messageid}
            AND     c.id = {circleid}
            CREATE  (m)-[r:PUB_TO]->(c)
            SET     r.published_at = {now}
            RETURN  r
        `,
		Parameters: neoism.Props{
			"messageid": messageid,
			"circleid":  circleid,
			"now":       time.Now().Local(),
		},
		Result: &created,
	})
	return len(created) > 0
}

func (q Query) CreateMemberOfRelation(handle, circleid string) bool {
	joined := []struct {
		At time.Time `json:"r.at"`
	}{}
	q.cypherOrPanic(&neoism.CypherQuery{
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
			"now":    time.Now().Local(),
		},
		Result: &joined,
	})
	return len(joined) > 0
}

func (q Query) JoinBroadcastCircleOfUser(handle, target string) bool {
	created := []struct {
		At time.Time `json:"r.at"`
	}{}
	q.cypherOrPanic(&neoism.CypherQuery{
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
			"now":       time.Now().Local(),
		},
		Result: &created,
	})
	return len(created) > 0
}

func (q Query) CreateBlockRelationFromTo(handle, target string) bool {
	res := []struct {
		Handle string      `json:"u.handle"`
		Target string      `json:"t.handle"`
		R      neoism.Node `json:"r"`
	}{}
	q.cypherOrPanic(&neoism.CypherQuery{
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
	})
	return len(res) > 0
}

//
// Read
//

func (q Query) UserExistsByHandle(handle string) bool {
	found := []struct {
		Handle string `json:"u.handle"`
	}{}
	q.cypherOrPanic(&neoism.CypherQuery{
		Statement: `
            MATCH   (u:User)
            WHERE   u.handle = {handle}
            RETURN  u.handle
        `,
		Parameters: neoism.Props{
			"handle": handle,
		},
		Result: &found,
	})
	return len(found) > 0
}

func (q Query) CircleLinkedToPublicDomain(circleid string) bool {
	found := []struct {
		Id string `json:"c.id"`
	}{}
	q.cypherOrPanic(&neoism.CypherQuery{
		Statement: `
            MATCH   (c:Circle)-[:PART_OF]->(p:PublicDomain)
            WHERE   c.id = {id}
            RETURN  c.id
        `,
		Parameters: neoism.Props{
			"id": circleid,
		},
		Result: &found,
	})
	return len(found) > 0
}

func (q Query) UserPartOfCircle(handle, circleid string) bool {
	found := []struct {
		Id string `json:"c.id"`
	}{}
	q.cypherOrPanic(&neoism.CypherQuery{
		Statement: `
			MATCH   (u:User)-[:MEMBER_OF|CHIEF_OF]->(c:Circle)
			WHERE   u.handle = {handle}
			AND     c.id     = {id}
			RETURN  c.id
		`,
		Parameters: neoism.Props{
			"handle": handle,
			"id":     circleid,
		},
		Result: &found,
	})
	return len(found) > 0
}

func (q Query) MessageIsPublished(handle, messageid, circleid string) bool {
	found := []struct {
		R *neoism.Relationship `json:"r"`
	}{}
	q.cypherOrPanic(&neoism.CypherQuery{
		Statement: `
            MATCH   (u:User)-[:WROTE]->(m:Message)-[r:PUB_TO]->(c:Circle)
            WHERE   u.handle = {handle}
            AND     m.id     = {messageid}
            AND     c.id     = {circleid}
            RETURN  r
        `,
		Parameters: neoism.Props{
			"handle":    handle,
			"messageid": messageid,
			"circleid":  circleid,
		},
		Result: &found,
	})
	return len(found) > 0
}

func (q Query) GetMessageById(messageid string) bool {
	found := []struct {
		Id int `json:"m.id"`
	}{}
	q.cypherOrPanic(&neoism.CypherQuery{
		Statement: `
            MATCH (m:Message)
            WHERE m.id = {id}
            RETURN m.id
        `,
		Parameters: neoism.Props{
			"id": messageid,
		},
		Result: &found,
	})
	return len(found) > 0
}

func (q Query) HandleExists(handle string) bool {
	found := []struct {
		Handle string `json:"u.handle"`
	}{}
	q.cypherOrPanic(&neoism.CypherQuery{
		Statement: `
            MATCH (u:User {handle: {handle}})
            RETURN u.handle
        `,
		Parameters: neoism.Props{
			"handle": handle,
		},
		Result: &found,
	})
	return len(found) > 0
}

func (q Query) EmailExists(email string) bool {
	found := []struct {
		Email string `json:"u.email"`
	}{}
	q.cypherOrPanic(&neoism.CypherQuery{
		Statement: `
            MATCH (u:User {email: {email}})
            RETURN u.email
        `,
		Parameters: neoism.Props{
			"email": email,
		},
		Result: &found,
	})
	return len(found) > 0
}

func (q Query) SessionBelongsToSomeUser(sessionid string) bool {
	found := []struct {
		Handle string `json:"u.handle"`
	}{}
	q.cypherOrPanic(&neoism.CypherQuery{
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
	})
	return len(found) == 1
}

func (q Query) BlockExistsFromTo(handle, target string) bool {
	found := []struct {
		Relation int `json:"r"`
	}{}
	q.cypherOrPanic(&neoism.CypherQuery{
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
	})
	return len(found) > 0
}

//
// Update
//

//
// Delete
//
