package query

import (
	"fmt"
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
