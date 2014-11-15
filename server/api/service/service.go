package service

import (
	"./query"
	"encoding/json"
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/jmcvetta/neoism"
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
	Query *query.Query
}

//
// Return types
//
type Message struct {
	Id      string    `json:"m.id"`
	Author  string    `json:"t.handle"`
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
	s := &Svc{
		query.NewQuery(uri),
	}
	return s
}

//
// Checks
//

func (s Svc) UserExists(handle string) bool {
	return s.Query.UserExistsByHandle(handle)
}

func (s Svc) CircleExistsInPublicDomain(circleid string) bool {
	return s.Query.CircleLinkedToPublicDomain(circleid)
}

func (s Svc) CanSeeCircle(fromPerspectiveOf string, circleid string) bool {
	if s.Query.CircleLinkedToPublicDomain(circleid) ||
		s.Query.UserPartOfCircle(fromPerspectiveOf, circleid) {
		return true
	}
	return false
}

func (s Svc) UserCanPublishTo(handle, circleid string) bool {
	return s.Query.UserPartOfCircle(handle, circleid)
}

func (s Svc) UserCanRetractPublication(handle, messageid, circleid string) bool {
	return s.Query.MessageIsPublished(handle, messageid, circleid)
}

func (s Svc) MessageExists(messageid string) bool {
	return s.Query.GetMessageById(messageid)
}

func (s Svc) HandleIsUnique(handle string) bool {
	return !s.Query.HandleExists(handle)
}

func (s Svc) EmailIsUnique(email string) bool {
	return !s.Query.EmailExists(email)
}

func (s Svc) VerifySession(sessionid string) bool {
	return s.Query.SessionBelongsToSomeUser(sessionid)
}

func (s Svc) BlockExistsFromTo(handle, target string) bool {
	return s.Query.BlockExistsFromTo(handle, target)
}

//
// Creation
//

func (s Svc) CreateNewUser(handle, email, passwordHash string) bool {
	return s.Query.CreateUser(handle, email, passwordHash)
}

func (s Svc) MakeDefaultCirclesFor(handle string) bool {
	return s.Query.CreateDefaultCirclesForUser(handle)
}

func (s Svc) NewCircle(handle, circleName string, isPublic bool) (circleid string, ok bool) {
	return s.Query.CreateCircle(handle, circleName, isPublic)
}

func (s Svc) NewMessage(handle, content string) (messageid string, ok bool) {
	return s.Query.CreateMessage(handle, content)
}

func (s Svc) PublishMessageToCircle(messageid, circleid string) bool {
	return s.Query.CreatePublishedRelation(messageid, circleid)
}

func (s Svc) JoinCircle(handle, circleid string) bool {
	// [TODO] check that `handle` is not the cheif of the circle here
	return s.Query.CreateMemberOfRelation(handle, circleid)
}

func (s Svc) JoinBroadcast(handle, target string) bool {
	return s.Query.JoinBroadcastCircleOfUser(handle, target)
}

func (s Svc) CreateBlockFromTo(handle, target string) bool {
	return s.Query.CreateBlockRelationFromTo(handle, target)
}

//
// Deletion
//

func (s Svc) FreshInitialState() {
	s.Query.DeleteAllNodesAndRelations()
	s.Query.DatabaseInit()
}

func (s Svc) KickTargetFromCircles(handle, target string) {
	s.Query.DisconnectTargetFromAllHeldCircles(handle, target)
}

func (s Svc) DeleteUser(handle string) bool {
	return s.Query.DeleteUser(handle)
}

func (s Svc) UnpublishMessageFromCircle(messageid, circleid string) bool {
	return s.Query.DeletePublishedRelation(messageid, circleid)
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

func (s Svc) GetMessagesByHandle(target string) []Message {
	messages := make([]Message, 0)
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH     (t:User)-[:WROTE]->(m:Message)
            WHERE     t.handle = {target}
            RETURN    m.id
                    , t.handle
                    , m.content
                    , m.created
            ORDER BY  m.created
        `,
		Parameters: neoism.Props{
			"target": target,
		},
		Result: &messages,
	}); err != nil {
		panicErr(err)
	}

	return messages
}

func (s Svc) GetMessageById(handle, messageid string) (message *Message, found bool) {
	messages := make([]Message, 0)
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
			MATCH   (t:User)-[:WROTE]->(m:Message)-[:PUB_TO]->(c:Circle)<-[:MEMBER_OF|CHIEF_OF]-(u:User)
			WHERE   u.handle = {handle}
            AND     m.id     = {messageid}
			RETURN  m.id
	              , t.handle
	              , m.content
	              , m.created
		`,
		Parameters: neoism.Props{
			"handle":    handle,
			"messageid": messageid,
		},
		Result: &messages,
	}); err != nil {
		panicErr(err)
	}

	if ok := len(messages) > 0; ok {
		return &messages[0], ok
	} else {
		return nil, ok
	}
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
			"sessionid": "Token " + uniuri.NewLen(uniuri.UUIDLen),
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

func (s Svc) UnsetSessionId(sessionid string) bool {
	unset := []struct {
		Handle string `json:"u.handle"`
	}{}
	if err := s.Db.Cypher(&neoism.CypherQuery{
		Statement: `
            MATCH   (u:User)<-[so:SESSION_OF]-(a:AuthToken)
            WHERE   a.sessionid = {sessionid}
            DELETE  so, a
            RETURN  u.handle
        `,
		Parameters: neoism.Props{
			"sessionid": sessionid,
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

func (s Svc) UpdateContentOfMessage(messageid, content string) bool {
	updated := []struct {
		Content string
	}{}
	cypherOrPanic(s, &neoism.CypherQuery{
		Statement: `
            MATCH  (m:Message)
            WHERE  m.id        = {messageid}
            SET    m.content   = {content}
            SET    m.lastsaved = {now}
            RETURN m.content
        `,
		Parameters: neoism.Props{
			"messageid": messageid,
			"content":   content,
			"now":       time.Now().Local(),
		},
		Result: &updated,
	})
	return len(updated) > 0
}
