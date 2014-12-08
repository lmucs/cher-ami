package service

import (
	"../../types"
	"./query"
	"fmt"
	"time"
)

//
// Constants
//

const (
	// Reserved Circles
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

func MakeCircleUrl(circleid string) string {
	return API_URL + "/circles/" + circleid
}

func MakeCircleMembersUrl(circleid string) string {
	return MakeCircleUrl(circleid) + "/members"
}

func MakeMessageUrl(messageid string) string {
	return API_URL + "/messages/" + messageid
}

func AddMessageUrl(m *types.MessageView) {
	m.Url = MakeMessageUrl(m.Id)
}

func AddMessageUrlToArray(messages []types.PublishedMessageView) {
	for _, m := range messages {
		m.Url = MakeMessageUrl(m.Id)
	}
}

func formatCircleView(c query.RawCircleView) types.CircleResponse {
	var visibility string
	if c.Public == nil {
		visibility = "private"
	} else {
		visibility = "public"
	}
	formatted := types.CircleResponse{
		Name:        c.Name,
		Url:         MakeCircleUrl(c.Id),
		Description: c.Description,
		Owner:       c.Owner,
		Visibility:  visibility,
		Members:     MakeCircleMembersUrl(c.Id),
		Created:     c.Created,
	}
	return formatted
}

//
// Checks
//

func (s Svc) UserExists(handle string) bool {
	return s.Query.UserExistsByHandle(handle)
}

// Returned whether the target user exists and has not blocked handle
func (s Svc) UserExistsAndNoBlocking(handle, target string) bool {
	return s.Query.UserExistsByHandle(target) &&
		s.Query.NoBlockingRelationshipBetween(handle, target)
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

func (s Svc) NewCircle(handle, circleName string, isPublic bool) (types.CircleResponse, bool) {
	view, ok := s.Query.CreateCircle(handle, circleName, isPublic)
	return formatCircleView(view), ok
}

func (s Svc) NewMessage(handle, content string) (message types.MessageView, ok bool) {
	m, ok := s.Query.CreateMessage(handle, content)
	if ok {
		AddMessageUrl(&m)
		return m, ok
	} else {
		return types.MessageView{}, ok
	}
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
		formatted[i] = formatCircleView(c)
	}
	return formatted, len(formatted)
}

func (s Svc) CirclesUserIsPartOf(user string, before time.Time, limit int) (results []types.CircleResponse, count int) {
	circles, _ := s.Query.GetJoinedCirclesByHandle(user, before, limit)
	formatted := make([]types.CircleResponse, len(circles))
	for i, c := range circles {
		formatted[i] = formatCircleView(c)
	}
	return formatted, len(formatted)
}

func (s Svc) GetPasswordHash(handle string) (passwordHash []byte, ok bool) {
	return s.Query.GetPasswordHash(handle)
}

func (s Svc) GetCircleId(handle, circleName string) (circleid string) {
	return s.Query.GetCircleIdByName(handle, circleName)
}

func (s Svc) GetPublicMessagesByHandle(self, target string) ([]types.PublishedMessageView, bool) {
	if ok := s.UserExistsAndNoBlocking(target, self); ok {
		messages := s.Query.GetPublicPublishedMessagesByAuthor(target)
		AddMessageUrlToArray(messages)
		return messages, ok
	} else {
		return []types.PublishedMessageView{}, ok
	}
}

func (s Svc) GetMessagesByTargetInCircle(self, target, circleid string) ([]types.PublishedMessageView, bool) {
	if ok := s.CanSeeCircle(self, circleid) && s.UserExistsAndNoBlocking(target, self); ok {
		messages := s.Query.GetMessagesByHandleInCircle(target, circleid)
		AddMessageUrlToArray(messages)
		return messages, ok
	} else {
		return []types.PublishedMessageView{}, ok
	}
}

func (s Svc) GetMessagesInCircle(self, circleid string) ([]types.PublishedMessageView, bool) {
	if ok := s.CanSeeCircle(self, circleid); ok {
		messages := s.Query.GetMessageFeedOfCircle(circleid)
		AddMessageUrlToArray(messages)
		return messages, ok
	} else {
		return []types.PublishedMessageView{}, ok
	}
}

// Should only be used on the logged-in user
// retrieves the personalized feed of the user
func (s Svc) GetMessageFeedOfSelf(handle string) ([]types.PublishedMessageView, bool) {
	if ok := s.UserExists(handle); ok {
		messages := s.Query.GetMessageFeedOfHandle(handle)
		AddMessageUrlToArray(messages)
		return messages, ok
	} else {
		return []types.PublishedMessageView{}, ok
	}
}

func (s Svc) GetVisibleMessageById(handle, messageid string) (message types.MessageView, ok bool) {
	if message, ok := s.Query.GetVisibleMessageById(handle, messageid); ok {
		AddMessageUrl(&message)
		return message, ok
	} else {
		return types.MessageView{}, ok
	}

}

func (s Svc) GetHandleFromAuthorization(token string) (handle string, ok bool) {
	return s.Query.DeriveHandleFromAuthToken(token)
}

func (s Svc) GetVisibleUser(handle, target string) (result types.UserView, ok bool) {
	if user, ok := s.Query.GetVisibleUserByHandle(handle, target); !ok {
		return types.UserView{}, ok
	} else {
		// var blockedUsers types.UserView
		if handle == target {
			if blocked, count := s.Query.GetBlockedUsers(handle); count > 0 {
				user.Blocked = blocked
			}
		}
		fmt.Printf("%+v", user)
		circles, _ := s.Query.GetPublicCirclesByHandle(handle)
		formatted := make([]types.CircleResponse, len(circles))
		for i, c := range circles {
			formatted[i] = formatCircleView(c)
		}
		user.Circles = formatted
		return user, true
	}
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
