# The CherAmi API

This is an all-JSON API. All requests and responses should have a content-type header set to `application/json`.

All endpoints except signup (`POST /users`) and login (`POST /sessions`) require an `Authorization` header in which you pass in the token that you previously received from a successful login request. For example:

    Authorization: Token 8dsfg87ef23dkos9r9wjr32232re

If the authorization header is missing, or the token is invalid or expired, an HTTP 401 response is returned. After receiving a 401, a client should try to login (`POST /sessions`) again to obtain a new token.


# Group Users

## Signup [/users]
### Signup/create a new user [POST]
Create a user with only a handle, name, email, and password.  An initial status, reputation, and default circles will be created.  Other profile information is set in a different operation.
+ Request

        {
            "handle": "pelé",
            "name": "Edson Arantes do Nascimento",
            "email": "number10@brasil.example.com",
            "password": "Brasil Uber Alles"
        }
+ Response 201

        {
            "url": "http://cher-ami.example.com/users/206",
            "handle": "pelé",
            "name": "Edson Arantes do Nascimento",
            "email": "number10@brasil.example.com",
            "status": "new",
            "reputation": 1,
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
            "reason": "handle already used|email already used"
        }

## Login and Logout [/sessions]

### Login [POST]
If the username/password combination is valid, generate and return a token.
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
            "reason": "invalid email/password combination"
        }

### Logout [DELETE]
A token is passed in and, if it is valid, the server will invalidate it.
+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
+ Response 204
+ Response 400

        {
            "reason": "cannot invalidate token because it is missing or already invalid or expired"
        }


## User Search [/users{?circle,name,before,limit}]

### Get users [GET]
Search for, or simply fetch, some desired set of users, paginated. Users are returned by the datetime they joined CherAmi, descending.   Only a subset of user data is returned; a different endpoint is used to get the complete user data.

+ Parameters
    + circle (optional, string) ... if present, only return users from this circle
    + name (optional, string) ... return users whose handles or names begin with this value (good for autocomplete)
    + before (optional, string, `2014-01-01`) ... return only users joined before this datetime
    + limit (optional, number, `20`) ... max number of results to return, for pagination, default 20, max 100

+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
+ Response 200

        [
            {
                "url": "http://cher-ami.example.com/users/206",
                "handle": "pelé",
                "name": "Edson Arantes do Nascimento",
            },
            . . .
        ]
+ Response 400

        {
            "reason": "malformed json|illegal date|illegal limit|limit out of range"
        }
+ Response 401

        {
            "reason": "missing, illegal, or expired token"
        }
+ Response 403

        {
            "reason": "you do not own or belong to this circle"
        }

## User [/users/{handle}]

### Get user by handle [GET]
Get complete user data, include all profile information as well as blocked users and circle membership.
+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
+ Response 200

        {
            "url": "http://cher-ami.example.com/users/206",
            "avatar_url": "https://images.cher-ami.example.com/users/206",
            "handle": "pelé",
            "name": "Edson Arantes do Nascimento",
            "email": "number10@brasil.example.com",
            "status": "retired, but coaching",
            "reputation": 1435346,
            "joined": "2011-06-30",
            "circles": [
                {"name": "public", "url": "http://cher-ami.example.com/circles/207"},
                {"name": "gold", "url": "http://cher-ami.example.com/circles/208"}
            ]
        }
+ Response 401

        {
            "reason": "missing, illegal, or expired token"
        }

### Edit user [PATCH]
Change only basic user information here such as display name, email, avatar url, and status. This endpoint is not used for complex properties like the set of users that this user has blocked, or the circles in which this user participates.Certain user data such as the internal id, handle, and join date cannot be changed at all.
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
                "action": "block|unblock"
            }
+ Response 204
+ Response 400

        {
            "reason": "malformed json|missing handle|missing action|unknown action"
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

## Reputation [/users/:id/reputation]

### Get a user's reputation [GET]
+ Response 200

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
