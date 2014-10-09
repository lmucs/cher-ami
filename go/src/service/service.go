package service

import (
	"fmt"
	// "github.com/dchest/uniuri"
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
	// Nodes must have at least one property to allow uniquely creation
	publicdomain, _, err := s.Db.GetOrCreateNode("PublicDomain", "u", neoism.Props{
		"u": true,
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

func (s Svc) GoodSessionCredentials(handle string, sessionid string) (bool, error) {
	found := []struct {
		Handle string `json:"user.handle"`
	}{}
	err := s.Db.Cypher(&neoism.CypherQuery{
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

	return len(found) == 1, err
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
	if len(made) != 1 {
		panic(fmt.Sprintf("Incorrect results len in query1()\n\tgot %d, expected 1\n", len(made)))
	}
	return err
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
