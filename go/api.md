# The CherAmi API

This is an all-JSON API. All requests and responses should have a content-type header set to `application/json`.

All endpoints except signup (`POST /users`) and login (`POST /sessions`) require an `Authorization` header in which you pass in the token that you previously received from a successful login request. For example:

    Authorization: Token 8dsfg87ef23dkos9r9wjr32232re

If the authorization header is missing, or the token is invalid or expired, an HTTP 401 response is returned. After receiving a 401, a client should try to login (`POST /sessions`) again to obtain a new token.


# Group Users

## Signup [/users]
### Signup a new user [POST]
+ Request

        {
            "handle": "my handle",
            "email": "my email",
            "password": "my password"
        }
+ Response 201

        {
            "url": "http://cher-ami.example.com/users/206",
            "handle": "pelé",
            "name": "Edson Arantes do Nascimento",
            "circles": [
                {"name": "public", "url": "http://cher-ami.example.com/circles/207"},
                {"name": "gold", "url": "http://cher-ami.example.com/circles/208"}
            ]
        }
+ Response 400

        {
            "reason": "malformed json|handle already used|email already used|password too weak"
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
+ Response 404

        {
            "reason": "invalid email/password combination"
        }

### Logout [DELETE]
+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
+ Response 204


## Get Users [/users{?circle,before,limit}]

### Get users [GET]
Users are returned by join datetime, descending

+ Parameters
    + circle (optional, string) ... if present, only return users from this circle
    + before (optional, string, `2014-01-01`) ... return only users joined before this datetime
    + limit (optional, number, `20`) ... max number of results to return, for pagination

+ Request
    + Headers

            Authorization: Token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx


+ Response 200
+ Response 400
+ Response 401
+ Response 403

## Single User [/users/{handle}]

### Get user by handle [GET]
Get user's profile and other information
+ Response 200

### Edit user [PATCH]
+ Response 200

### Delete user [DELETE]
+ Response 204

## Blocking [/users/:id/blocked]

### Block or unblock user [PATCH]
+ Response 200

### Get the users I've blocked [GET]
+ Response 200

## Reputation [/users/:id/reputation]

### Get a user's reputation [GET]
+ Response 200

### Adjust user reputation up/down [PATCH]
+ Response 200

### Set user reputation directly [PUT]
+ Response 204

# Group Circles

## All circles [/circles]
### Create circle [POST]
+ Response 201

### Search for circles [GET]
+ Response 200

## Single Circle [/circles/{id}]
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

## All Messages [/messages]

### Create message [POST]
+ Response 201

### Get messages [GET]
+ Response 200

## Single Message [/messages/{id}]

### Get message by id [GET]
+ Response 200

### Delete message [DELETE]
+ Response 204

## All Comments [/messages/{id}/comments]

### Post comment [POST]
+ Response 201

### Get comments for message [GET]
+ Response 200

## Single Comment [/messages/{id}/comments/:id]

### Get comment by id [GET]
+ Response 200

### Delete comment [DELETE]
+ Response 204
