package service

import (
	"../../types"
	"./query"
	"time"
)

//
// Constants
//

const (
	// Reserved Circles
	GOLD           = "Gold"
	BROADCAST      = "Broadcast"
	CHERAMI_PREFIX = "http://"
	DOMAIN         = "cherami.io"
	CHERAMI_URL    = CHERAMI_PREFIX + DOMAIN
	API_URL        = CHERAMI_URL + "/api"
)

//
// Service Types
//

type Svc struct {
	Query *query.Query
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

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
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

func (s Svc) VerifyAuthToken(token string) bool {
	return s.Query.AuthTokenBelongsToSomeUser(token)
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

func (s Svc) NewCircle(handle, circleName string, isPublic bool,
) (circleid string, ok bool) {
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

func (s Svc) SearchForUsers(circle, nameprefix string, skip, limit int, sort string,
) (results string, count int) {
	return s.Query.SearchForUsers(circle, nameprefix, skip, limit, sort)
}

func (s Svc) SearchCircles(user string, before time.Time, limit int) (results []types.CircleResponse, count int) {
	circles := s.Query.SearchCircles(user, before, limit)
	formatted := make([]types.CircleResponse, len(circles))
	for i, c := range circles {
		var visibility string
		if c.Private != nil {
			visibility = "private"
		} else {
			visibility = "public"
		}
		formatted[i] = types.CircleResponse{
			Name:        c.Name,
			Url:         API_URL + "/circles/" + c.Id,
			Description: c.Description,
			Owner:       c.Owner,
			Visibility:  visibility,
			Members:     API_URL + "/circles/" + c.Id + "/members",
			Created:     c.Created,
		}
	}
	return formatted, len(formatted)
}

func (s Svc) GetJoinedCircles(user string, skip, limit int) (results []types.CircleResponse, count int) {
	return s.Query.GetJoinedCirclesByHandle(user, skip, limit)
}

func (s Svc) GetPasswordHash(handle string) (passwordHash []byte, ok bool) {
	return s.Query.GetPasswordHash(handle)
}

func (s Svc) GetCircleId(handle, circleName string) (circleid string) {
	return s.Query.GetCircleIdByName(handle, circleName)
}

func (s Svc) GetMessagesByHandle(target string) []query.Message {
	return s.Query.GetAllMessagesByHandle(target)
}

func (s Svc) GetVisibleMessageById(handle, messageid string,
) (message *query.Message, ok bool) {
	return s.Query.GetVisibleMessageById(handle, messageid)
}

func (s Svc) GetHandleFromAuthorization(token string) (handle string, ok bool) {
	return s.Query.DeriveHandleFromAuthToken(token)
}

func (s Svc) GetSelf(handle string) (result types.OwnUserView, ok bool) {
	user := s.Query.GetVisibleUserByHandle(handle)
	return types.OwnUserView{
		Handle: user.Handle,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Gender: user.Gender,
		Birthday: user.Birthday,
		Bio: user.Bio,
		Interests: user.Interests,
		Languages: user.Languages,
		Location: user.Location,
		Circles: s.Query.GetJoinedCircles(handle, 0, 100),
		Blocked: s.Query.GetBlockedUsers(handle),
	}
}

func (s Svc) GetVisibleUser(handle string) (result types.UserView, ok bool) {
	user := s.Query.GetVisibleUserByHandle(handle)
	circles := s.Query.GetPublicCirclesByHandle(handle)
	formatted := make([]types.CircleResponse, len(circles))
	for i, c := range circles {
		var visibility string
		if c.Private != nil {
			visibility = "private"
		} else {
			visibility = "public"
		}
		formatted[i] = types.CircleResponse{
			Name:        c.Name,
			Url:         API_URL + "/circles/" + c.Id,
			Description: c.Description,
			Owner:       c.Owner,
			Visibility:  visibility,
			Members:     API_URL + "/circles/" + c.Id + "/members",
			Created:     c.Created,
		}
	}
	user.Circles = formatted
	return user, true
}

//
// Node Attributes
//

// Creates a new AuthToken node that points to a particular user
// returning the value of the token created
func (s Svc) SetGetNewAuthToken(handle string) (string, bool) {
	return s.Query.SetGetNewAuthTokenForUser(handle)
}

func (s Svc) SetNewPassword(handle, newPasswordHash string) bool {
	return s.Query.UpdatePassword(handle, newPasswordHash)
}

func (s Svc) DestroyAuthToken(token string) bool {
	return s.Query.DestroyAuthToken(token)
}

func (s Svc) SetGetName(handle, newName string) (string, bool) {
	return s.Query.SetGetUserName(handle, newName)
}

func (s Svc) UpdateContentOfMessage(messageid, content string) bool {
	return s.Query.UpdateMessageContent(messageid, content)
}

func (s Svc) UpdateUserAttribute(handle, resource, content string) bool {
	return s.Query.UpdateUserAttribute(handle, resource, content)
}
