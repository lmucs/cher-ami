# The CherAmi API

This is the API for the CherAmi social network.

This is an all-JSON API. All requests and responses, except those transferring audio, image, or video content, must have a content-type header set to `application/json`.

All endpoints except signup (`POST /users`) and login (`POST /sessions`) require an `Authorization` header in which you pass in the token that you previously received from a successful login request. For example:

    Authorization: Token 8dsfg87ef23dkos9r9wjr32232re

If the authorization header is missing, or the token is invalid or expired, an HTTP 401 response is returned. After receiving a 401, a client should try to login (`POST /sessions`) again to obtain a new token.

The API supports discovery of further endpoints, linking objects with absolute URIs.

# Group Users



## Signup [/users]



### Signup/create a new user [POST]
Create a user given only a handle, email, and password. The service will create an initial status, stir, and default circles, as well as record the creation timestamp. All other profile information is set using different operations.
+ Request

        {
            "handle": "pelé",
            "email": "number10@brasil.example.com",
            "password": "Brasil Uber Alles",
            "confirmpassword": "Brasil Uber Alles"
        }
+ Response 201

        {
            "url": "https://cher-ami.example.com/users/pelé",
            "handle": "pelé",
            "email": "number10@brasil.example.com",
            "stir": 10,
            "status": "Hi, I'm new here!",
            "joined": "2011-10-20T08:15Z",
            "circles": [
                {"name": "public", "url": "https://cher-ami.example.com/circles/207"},
                {"name": "gold", "url": "https://cher-ami.example.com/circles/208"}
            ]

        }
+ Response 400

        {
            "reason": ("malformed json"
                      |"Handle is a required field for signup"
                      |"Email is a required field for signup")
        }
+ Response 403

        {
            "reason": ("Passwords do not match"
                      |"Passwords must be at least 8 characters long")
        }
+ Response 409

        {
            "reason": "Sorry, handle or email is already taken"
        }


## Login and Logout [/sessions]


### Login [POST]
If the given username-password combination is valid, generate and return a token.
+ Request

        {
            "handle": "pelé",
            "password": "a string"
        }
+ Response 201

        {
           "handle": "pelé",
           "token": "Token hu876xvyft3ufib230ffn0spdfmwefna"
        }
+ Response 400

        {
            "reason": "malformed json"
        }
+ Response 403

        {
            "reason": "Invalid username or password"
        }



### Logout [DELETE]
The token is passed in a header (not as a parameter in the URL) and, if it is valid, the server will invalidate it.
+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
+ Response 204
+ Response 403

        {
            "reason": "Cannot invalidate token because it is missing"
        }



## User Search [/users{?circle,nameprefix,skip,limit,sort}]



### Get users [GET]
Fetch a desired set of users. You may filter by circle or leading characters of a name. You _must_ specify a sort order. The results _will_ be paginated since there is a potential for returning millions of users. Only a subset of user data is returned; however, the url to get the complete data _is_ returned, in good HATEOAS-style.

+ Parameters
    + circle (optional, string, `393`) ... only return users from the circle with this id
    + nameprefix (optional, string, `sta`) ... only return users whose names begin with this value (good for autocomplete)
    + skip (optional, number, `0`) ... number of results to skip, for pagination, default 0, min 0
    + limit (optional, number, `20`) ... max number of results to return, for pagination, default 20, min 1, max 100
    + sort (required, string, `joined`)

        sort results by name ascending, stir descending, or join datetime descending (newest users first)
        + Values
            + `name`
            + `stir`
            + `joined`

+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

+ Response 200

        [
            {
                "url": "https://cher-ami.example.com/users/pelé",
                "handle": "pelé",
                "name": "Edson Arantes do Nascimento",
                "stir": 303,
                "joined": "2011-10-20T08:15Z"
            },
            . . .
        ]
+ Response 400

        {
            "reason": ("malformed json"
                      |"missing sort"
                      |"no such sort"
                      |"malformed skip"
                      |"malformed limit")
        }
+ Response 401

        {
            "response": "Failed to authenticate user request",
            "reason": "Missing, illegal or expired token",
        }
+ Response 403

        {
            "reason": ("you do not own or belong to this circle"
                      |"skip out of range"
                      |"limit out of range")
        }



## User [/users/{handle}]



### Get user by handle [GET]
Get _complete_ user data, including all profile information as well as blocked users and circle membership.
+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
+ Response 200

        {
            "url": "https://cher-ami.example.com/users/pelé",
            "avatar_url": "https://images.cher-ami.example.com/users/pelé",
            "handle": "pelé",
            "name": "Edson Arantes do Nascimento",
            "email": "number10@brasil.example.com",
            "status": "retired, but coaching",
            "stir": 1435346,
            "joined": "2011-10-20T08:15Z",
            "circles": [
                {"name": "public", "url": "https://cher-ami.example.com/circles/207"},
                {"name": "gold", "url": "https://cher-ami.example.com/circles/208"}
                {"name": "coaches", "url": "https://cher-ami.example.com/circles/5922"}
            ]
        }
+ Response 401

        {
            "response": "Failed to authenticate user request",
            "reason":   "Missing, illegal or expired token",
        }

+ Response 404

        {
            "reason": "no such user"
        }



### Edit user [PATCH]
Change only basic user information here such as display name, email, and status. Use a different endpoint for complex properties like the set of users that this user has blocked, or the circles in which this user participates. Also use different endpoints to adjust stir and to upload a new avatar picture. Note that certain user data, such as the internal id, handle, and join date, cannot be changed at all.
+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    + Body

            {
                "name": "New name (optional)",
                "email": "New email (optional)",
                "status": "New status (optional)",
            }
+ Response 204
+ Response 401

        {
            "response": "Failed to authenticate user request",
            "reason":   "Missing, illegal or expired token",
        }

+ Response 403

        {
            "reason": "you can only edit yourself unless you are an admin"
        }

### Delete user [DELETE]
+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    + Body

            {
                "password": "Brasil Uber Alles"
            }
+ Response 204
+ Response 401

        {
            "response": "Failed to authenticate user request",
            "reason":   "Missing, illegal or expired token",
        }

        or

        {
            "reason": "Invalid password, please try again"
        }

+ Response 403

        {
            "reason": "you can only delete yourself unless you are an admin"
        }



## Blocking [/users/{handle}/blocked]



### Block or unblock user [PATCH]
If user A blocks user B, then B is removed from all of A's circles, public and private.  As long as B is blocked by A, B will not be allowed to join any of A's circles.
+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    + Body

            {
                "handle": "pelé",
                "action": ("block"|"unblock")
            }
+ Response 204
+ Response 400

        {
            "reason": ("malformed json"|"missing handle"|"missing action"|"unknown action")
        }
+ Response 401

        {
            "reason": "missing, illegal, or expired token"
        }
+ Response 403

        {
            "reason": "you can not block yourself"
        }
+ Response 404

        {
            "reason": "no such user"
        }



## Viewing Blocked Users [/users/{handle}/blocked{?skip,limit}]



### Get blocked users [GET]
Fetch the list of blocked users for the given user, paginated. The blocked users will always be returned in alphabetical order by handle.

+ Parameters
    + skip (optional, number, `10`) ... number of results to skip, default is 0, min 0
    + limit (optional, number, `20`) ... max number of results to return, for pagination, default 20, min 1, max 100

+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
+ Response 200

        [
            {
                "url": "https://cher-ami.example.com/users/liane",
                "handle": "pelé",
                "name": "Liane Cartman",
                "stir": 303,
                "joined": "2011-10-20T08:15Z"
            },
            . . .
        ]
+ Response 400

        {
            "reason": ("malformed json"|"malformed skip"|"malformed limit")
        }
+ Response 401

        {
            "reason": "missing, illegal, or expired token"
        }
+ Response 403

        {
            "reason": ("limit out of range"|"you not allowed to see this user's blocked list")
        }
+ Response 404

        {
            "reason": "no such user"
        }




## Avatar [/users/{handle}/avatar]



### Upload avatar [PUT]
+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
            Content-type: image/png
    + Body

            ... image content ...
+ Response 204
+ Response 400

        {
            "reason": "bad media"
        }
+ Response 401

        {
            "reason": "missing, illegal, or expired token"
        }
+ Response 403

        {
            "reason": "you can not set others' avatars unless you are an admin"
        }
+ Response 404

        {
            "reason": "no such user"
        }



# Group Circles



## Circle Creation [/circles]



### Create circle [POST]
Create a circle given only a name and visibility setting, setting the owner to the currently logged-in user. Members will be added to the circle using a different endpoint (whose url is part of the returned resource).
+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    + Body

            {
                "circleName": "bffs",
                "public": true|false
            }
+ Response 201

        {
            "name": "bffs",
            "url": "https://cher-ami.example.com/circles/2997",
            "owner": "wendy",
            "visibility": "private",
            "members": "https://cher-ami.example.com/circles/2997/members",
            "creation": "2011-10-20T14:22:09Z",
            "stir": 0
        }
+ Response 400

        {
            "reason": "malformed json"
        }
+ Response 401

        {
            "reason": "missing, illegal, or expired token"
        }
+ Response 403

        {
            "reason": "bffs is a reserved circle name"
        }
+ Response 409

        {
            "reason": "name already used"
        }



## Circle Search [/circles{?user,before,limit}]



### Search for circles [GET]
Fetch the circles a user is part of, if the optional user parameter is absent this will return the circles that the authenticated user is apart of. The results will be paginated. Circles are returned in order of descending creation date.  We may add custom sorting capability in the future.

+ Parameters
    + user (optional, string, `alice`) ... only return circles owned by this user
    + before (optional, string, `2015-02-28`) ... only return circles created before this date (YYYY-MM-DD)
    + limit (optional, number, `20`) ... max number of results to return, for pagination, default 20, min 1, max 100

+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
+ Response 200

        [
            {
                "name": "bffs",
                "url": "https://cher-ami.example.com/circles/2997",
                "description": "All my closest friends",
                "owner": "wendy",
                "visibility": "private",
                "members": "https://cher-ami.example.com/circles/2997/members",
                "stir": 75,
                "created": "2011-10-20T14:22:09Z"
            },
            . . .
        ]
+ Response 400

        {
            "reason": ("malformed json"|"malformed before"|"malformed limit")
        }
+ Response 401

        {
            "reason": "missing, illegal, or expired token"
        }
+ Response 403

        {
            "reason": "limit out of range"
        }



## Circle [/circles/{id}]



### Get circle by id [GET]
Get complete circle data for the circle with the given id.
+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
+ Response 200

        {
            "name": "bffs",
            "url": "https://cher-ami.example.com/circles/2997",
            "description": "All my closest friends",
            "owner": "wendy",
            "visibility": "private",
            "members": "https://cher-ami.example.com/circles/2997/members",
            "stir": 80,
            "creation": "2011-10-20T14:22:09Z"
        }
+ Response 401

        {
            "reason": "missing, illegal, or expired token"
        }
+ Response 404

        {
            "reason": "no such circle"
        }



### Edit circle info [PATCH]
Edits only the name and description of the circle. Members are managed elsewhere. You cannot ever change the owner or creation time.
+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    + Body

            {
                "name": "New name (optional)",
                "description": "New description (optional)",
            }
+ Response 204
+ Response 401

        {
            "reason": "missing, illegal, or expired token"
        }
+ Response 403

        {
            "reason": "you can only edit circles you own unless you are an admin"
        }



## Get Circle Members [/circles/{id}/members{?skip,limit}]



### Get circle members [GET]
Fetch the list of members of this circle, paginated. The members will always be returned in alphabetical order by handle.

+ Parameters
    + skip (optional, number, `10`) ... number of results to skip, default is 0, min 0
    + limit (optional, number, `20`) ... max number of results to return, for pagination, default 20, min 1, max 100

+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
+ Response 200

        [
            {
                "url": "https://cher-ami.example.com/users/towelie",
                "handle": "towelie",
                "name": "Smart Towel RG-400",
                "stir": 420,
                "joined": "2011-04-20T16:20Z"
            },
            . . .
        ]
+ Response 400

        {
            "reason": ("malformed json"|"malformed skip"|"malformed limit")
        }
+ Response 401

        {
            "reason": "missing, illegal, or expired auth token"
        }
+ Response 403

        {
            "reason": "limit out of range"
        }
+ Response 404

        {
            "reason": "no such circle that you can see"
        }



## Manage Circle Members [/circles/{id}/members]



### Add/remove circle members [PATCH]
If a circle is public, all user can let themselves in, unless blocked by the circle owner. If private, only the owner can add. To remove a user, the requestor must be that very user, the circle owner, or an admin.
+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    + Body

            {
                "handle": "sharon",
                "action": ("add"|"remove")
            }
+ Response 204
+ Response 400

        {
            "reason": ("malformed json"|"missing handle"|"missing action"|"unknown action")
        }
+ Response 401

        {
            "reason": "missing, illegal, or expired auth token"
        }
+ Response 403

        {
            "reason": ("no such circle visible to you"|"only owner can add others to private circles"|"blocked by circle owner"|"not allowed to remove")
        }



# Group Messages



## Message Creation [/messages]



### Create message [POST]
Creates a message for a given circle, with the given content. Optionally, the creator can set a minimum stir threshold for viewing (defaults to 0, in a way). Server sets the id, creation timestamp, and author.
+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    + Body

            {
                "circle": 488,
                "min_stir": 500,
                "content": "There are no such things as stupid questions, only stupid people"
            }
+ Response 201

        {
            "content": "There are no such things as stupid questions, only stupid people",
            "url": "https://cher-ami.example.com/messages/98",
            "author": "garrison",
            "min_stir": 500,
            "creation": "2011-10-20T14:22:09Z"
        }
+ Response 400

        {
            "reason": ("malformed json"|"missing circle"|"missing content")
        }
+ Response 403

        {
            "reason": "no such circle or you are not allowed to post to it"
        }
+ Response 413

        {
            "reason": "content too large",
            "max_comment_size_in_bytes": 2048
        }



## Message Search [/messages{?circle,name,before,limit}]

### Get messages [GET]
Fetch the messages for the given circle or user, paginated. The messages will always be returned in order of descending creation date. Only messages corresponding to circles that the current user is allowed to see , and that the user has enough stir to see, will be returned.

+ Parameters
    + circle (optional, string, `284`) ... only return messages from this circle, required if name not supplied
    + name (optional, string, `liane`) ... only return messages authored by this user, required if circle not supplied
    + before (optional, string, `2015-02-28T22:11:07Z`) ... only return messages created before this datetime
    + limit (optional, number, `20`) ... max number of results to return, for pagination, default 20, min 1, max 100

+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
+ Response 200

        [
            {
                "url": "https://cher-ami.example.com/messages/802",
                "author": "stan",
                "content": "I'm not getting on this bus",
                "date": "2012-10-20T14:22:09Z"
            },
            . . .
        ]
+ Response 400

        {
            "reason": ("malformed json"|"malformed before"|"malformed limit"|"circle or user required"),
            "details": ("before must be a datetime in ISO 8601 format"|"limit must be a positive integer")
        }
+ Response 401

        {
            "reason": "missing, illegal, or expired auth token"
        }
+ Response 403

        {
            "reason": "limit out of range",
            "details": "limit must be between 1 and 100, inclusive"
        }
+ Response 404

        {
            "reason": ("no such circle that you can see"|"no such user")
        }



## Message [/messages/{id}]

### Get message by id [GET]
Get the message with the given id.
+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
+ Response 200

        {
            "url": "https://cher-ami.example.com/messages/802",
            "author": "stan",
            "content": "I'm not getting on this bus",
            "date": "2012-10-20T14:22:09Z"
        }
+ Response 401

        {
            "reason": "missing, illegal, or expired auth token"
        }
+ Response 403

        {
            "reason": "you don't have enough stir to see this message",
            "stir_required": 250
        }
+ Response 404

        {
            "reason": "no such message in any circle you can see"
        }

### Delete message [DELETE]
+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
+ Response 204
+ Response 401

        {
            "reason": "missing, illegal, or expired auth token"
        }
+ Response 404

        {
            "reason": "you are not the author of any such message"
        }



## Edit Message [/messages/{id}{&circles,content}]

### Patch a message by id [PATCH]
Edit an existing message by id. Is used to change properties like the message's content but will be more commonly used to publish a message. Patching a message will return a 200 and the newly-updated field values as well as the `published` field. At least one parameter must be supplied for a successful PATCH request.

+ Parameters
    + circles (optional, []string, `["circleid_001", "circleid_075"]`) ... Target circle(s) that the message should be posted to
    + content (optional, string, `some new content`) ... Set the new content of this message
+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

    + Body

            {
                "circles": ["circleid_001", "circleid_075"],
                "content": "Hello world ... again"
            }
+ Response 200

        {
            "updated": "https://cher-ami.example.com/messages/802",
            "circles": [
                {"circleid": "circleid_001", "published_at": "2012-10-20T14:22:09Z"},
                {"circleid": "circleid_075", "published_at": "2012-10-20T14:22:11Z"}
            ],
            "dateModified": "2012-10-20T14:22:09Z",
            "published": true
        }
+ Response 400

        {
            "response": "Failed to patch message",
            "reason": ("No field to patch specified"|"Some specified circle did not exist, or could not be published to")
        }
+ Response 401

        {
            "reason": "missing, illegal, or expired auth token"
        }
+ Response 404

        {
            "reason": "you are not the author of any such message"
        }

## Comment Creation [/messages/{id}/comments]



### Post comment [POST]
Post a comment to the given message. Comments are text-only. The server sets the timestamp and sets the author to the currently logged in user.
+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    + Body

            {
                "content": "You killed Kenny! You BASTARD!"
            }
+ Response 201

        {
            "content": "You killed Kenny! You BASTARD!",
            "url": "https://cher-ami.example.com/messages/983/comments/22",
            "author": "kyle",
            "creation": "2015-10-20T14:22:09Z"
        }
+ Response 400

        {
            "reason": ("malformed json"|"missing content")
        }
+ Response 403

        {
            "reason": "no such message or you are not allowed to comment on it"
        }
+ Response 413

        {
            "reason": "content too large",
            "max_comment_size_in_bytes": 512
        }



## Comment Search [/messages/{id}/comments{?before,limit}]



### Get comments for message [GET]
Fetch the comments for the given message, paginated. The comments will always be returned in order of descending creation date. If the current user is not an admin, the message must be in a public circle or a private circle to which the current user belongs AND the current user must not blocked by the message author or circle owner.

+ Parameters
    + before (optional, string, `2015-02-28T22:11:07Z`) ... only return comments created before this datetime
    + limit (optional, number, `20`) ... max number of results to return, for pagination, default 20, min 1, max 100

+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
+ Response 200

        [
            {
                "content": "You killed Kenny! You BASTARD!",
                "url": "https://cher-ami.example.com/messages/983/comments/22",
                "author": "kyle",
                "creation": "2015-10-20T14:22:09Z"
            },
            . . .
        ]
+ Response 400

        {
            "reason": ("malformed json"|"malformed before"|"malformed limit"),
            "details": ("before must be a datetime in ISO 8601 format"|"limit must be a positive integer")
        }
+ Response 401

        {
            "reason": "missing, illegal, or expired auth token"
        }
+ Response 403

        {
            "reason": "limit out of range",
            "details": "limit must be between 1 and 100, inclusive"
        }
+ Response 404

        {
            "reason": "message not found in any circle you can see",
            "message_id": 20,
        }



## Comment [/messages/{id}/comments/{id}]



### Get comment by id [GET]
Get the comment. If the current user is not an admin, the comment must be on a message of a public circle or a private circle to which the current user belongs AND the current user must not be blocked by either the message author, the comment author, or the circle owner.
+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
+ Response 200

        {
            "content": "You killed Kenny! You BASTARD!",
            "url": "https://cher-ami.example.com/messages/983/comments/22",
            "author": "kyle",
            "creation": "2015-10-20T14:22:09Z"
        }
+ Response 401

        {
            "reason": "missing, illegal, or expired auth token"
        }
+ Response 404

        {
            "reason": "comment not found in any circle you can see",
            "message_id": 20,
            "comment_id": 7
        }



### Delete comment [DELETE]
Permanently delete the comment. Current user must be the comment author or an admin.
+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
+ Response 204
+ Response 401

        {
            "reason": "missing, illegal, or expired auth token"
        }
+ Response 404

        {
            "reason": "comment not found among those you have authored",
            "message_id": 20,
            "comment_id": 7
        }
