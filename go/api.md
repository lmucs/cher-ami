# The CherAmi API

This is an all-JSON API. All requests and responses should have a content-type header set to `application/json`.

All endpoints except signup (`POST /users`) and login (`POST /sessions`) require an `Authorization` header in which you pass in the token that you previously received from a successful login request. For example:

    Authorization: Token 8dsfg87ef23dkos9r9wjr32232re

If the authorization header is missing, or the token is invalid or expired, an HTTP 401 response is returned. After receiving a 401, a client should try to login (`POST /sessions`) again to obtain a new token.


# Group Users

## Signup [/users]
### Signup/create a new user [POST]
Create a user given only a handle, name, email, and password.  The service will create an initial status, reputation, and default circles, as well as record the creation datetime.  All other profile information is set in a different operation.
+ Request

        {
            "handle": "pelé",
            "name": "Edson Arantes do Nascimento",
            "email": "number10@brasil.example.com",
            "password": "Brasil Uber Alles"
        }
+ Response 201

        {
            "url": "http://cher-ami.example.com/users/pelé",
            "handle": "pelé",
            "name": "Edson Arantes do Nascimento",
            "email": "number10@brasil.example.com",
            "status": "new",
            "reputation": 1,
            "joined": "2011-10-20",
            "circles": [
                {"name": "public", "url": "http://cher-ami.example.com/circles/207"},
                {"name": "gold", "url": "http://cher-ami.example.com/circles/208"}
            ]
        }
+ Response 400

        {
            "reason": "malformed json"
        }
+ Response 403

        {
            "reason": "password too weak"
        }
+ Response 409

        {
            "reason": ("handle already used"|"email already used")
        }

## Login and Logout [/sessions]

### Login [POST]
If the given username-password combination is valid, generate and return a token.
+ Request

        {
            "handle": "a string",
            "password": "a string"
        }
+ Response 201

        {
           "token": "hu876xvyft3ufib230ffn0spdfmwefna"
        }
+ Response 400

        {
            "reason": "malformed json"
        }
+ Response 403

        {
            "reason": "invalid handle-password combination"
        }

### Logout [DELETE]
The token is passed in a header (not as a parameter in the URL) and, if it is valid, the server will invalidate it.
+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
+ Response 204
+ Response 403

        {
            "reason": "cannot invalidate token because it is missing or already invalid or expired"
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

        sort results by name ascending, reputation descending, or join datetime descending (newest users first)
        + Values
            + `name`
            + `reputation`
            + `joined`

+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
+ Response 200

        [
            {
                "url": "http://cher-ami.example.com/users/pelé",
                "handle": "pelé",
                "name": "Edson Arantes do Nascimento",
                "reputation": 303,
                "joined": "2011-10-20"
            },
            . . .
        ]
+ Response 400

        {
            "reason": ("malformed json"|"missing sort"|"no such sort"|"malformed skip"|"malformed limit")
        }
+ Response 401

        {
            "reason": "missing, illegal, or expired token"
        }
+ Response 403

        {
            "reason": ("you do not own or belong to this circle"|"skip out of range"|"limit out of range")
        }

## User [/users/{handle}]

### Get user by handle [GET]
Get _complete_ user data, including all profile information as well as blocked users and circle membership.
+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
+ Response 200

        {
            "url": "http://cher-ami.example.com/users/pelé",
            "avatar_url": "https://images.cher-ami.example.com/users/pelé",
            "handle": "pelé",
            "name": "Edson Arantes do Nascimento",
            "email": "number10@brasil.example.com",
            "status": "retired, but coaching",
            "reputation": 1435346,
            "joined": "2011-06-30",
            "circles": [
                {"name": "public", "url": "http://cher-ami.example.com/circles/207"},
                {"name": "gold", "url": "http://cher-ami.example.com/circles/208"}
                {"name": "coaches", "url": "http://cher-ami.example.com/circles/5922"}
            ]
        }
+ Response 401

        {
            "reason": "missing, illegal, or expired token"
        }
+ Response 404

        {
            "reason": "no such user"
        }


### Edit user [PATCH]
Change only basic user information here such as display name, email, and status. Use a different endpoint for complex properties like the set of users that this user has blocked, or the circles in which this user participates. Also use different endpoints to adjust reputation and to upload a new avatar picture. Note that certain user data, such as the internal id, handle, and join date, cannot be changed at all.
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
            "reason": "missing, illegal, or expired token"
        }
+ Response 403

        {
            "reason": "you can only edit yourself unless you are an admin"
        }

### Delete user [DELETE]
+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
+ Response 204
+ Response 401

        {
            "reason": "missing, illegal, or expired token"
        }
+ Response 403

        {
            "reason": "you can only delete yourself unless you are an admin"
        }

## Blocking [/users/{handle}/blocked]

### Block or unblock user [PATCH]
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

### Get blocked users [GET]
+ Response 200

## Reputation [/users/{id}/reputation]

### Adjust user reputation +/- [PATCH]
+ Response 200

### Set user reputation directly [PUT]
+ Response 204

# Group Circles

## Circle Creation [/circles]
### Create circle [POST]
+ Response 201

## Circle Search [/circles{?user,before,limit}]
### Search for circles [GET]
+ Response 200

## Circle [/circles/{id}]
### Get circle by id [GET]
+ Response 200

### Edit circle info [PATCH]
+ Response 200

## Circle Members [/circles/{id}/members]

### Get circle members [GET]
+ Response 200

### Add/remove circle members [PATCH]
+ Response 200


# Group Messages

## Message Creation [/messages]

### Create message [POST]
+ Response 201

## Message Search [/messages{?circle,before,limit}]

### Get messages [GET]
+ Response 200

## Message [/messages/{id}]

### Get message by id [GET]
+ Response 200

### Delete message [DELETE]
+ Response 204

## Comment Creation [/messages/{id}/comments]

### Post comment [POST]
+ Response 201

## Comment Search [/messages/{id}/comments{?before,limit}]

### Get comments for message [GET]
+ Response 200
+ Response 400
+ Response 401
+ Response 404

## Comment [/messages/{id}/comments/{id}]

### Get comment by id [GET]
+ Response 200
+ Response 401
+ Response 403
+ Response 404

### Delete comment [DELETE]
+ Response 204
+ Response 401
+ Response 403
