# Health dot JSON #
## API Health

# Key

    {
        "No work done at all":               0
      , "Skeletoned":                        1
      , "Weak":                              2
      , "Partial functionality implemented": 4
      , "Full coverage":                     8
    }

## Endpoints

```json
{
    "method":         "POST"
  , "uri":            "/signup"
  , "desc":           "Sign up new user"
  , "testing":        8
  , "implementation": 8
  , "note":           ""
}
```

```json
{
    "method":         "POST"
  , "uri":            "/changepassword"
  , "desc":           "Change password"
  , "testing":        8
  , "implementation": 4
  , "note":           "route needs change"
}
```

```json
{
    "method":         "POST"
  , "uri":            "/sessions"
  , "desc":           "Login (new Auth token)"
  , "testing":        8
  , "implementation": 8
  , "note":           ""
}
```

```json
{
    "method":         "DELETE"
  , "uri":            "/sessions"
  , "desc":           "Logout (delete Auth token)"
  , "testing":        4
  , "implementation": 8
  , "note":           ""
}
```

```json
{
    "method":         "GET"
  , "uri":            "/users/:handle"
  , "desc":           "Get single user"
  , "testing":        1
  , "implementation": 2
  , "note":           "route doesn't follow specs, takes in query params"
}
```

```json
{
    "method": "GET"
  , "uri":    "/users"
  , "desc":   "Search For Users"
  , "testing":        2
  , "implementation": 4
  , "note":           "needs a thorough testing to 8-8"
}
```

```json
{
    "method":         "DELETE"
  , "uri":            "/users/:handle"
  , "desc":           "Delete User"
  , "testing":        4
  , "implementation": 4
  , "note":           "needs GET /user for complete testing"
}
```

```json
{
    "method":         "GET"
  , "uri":            "/messages"
  , "desc":           "Get Authored Messages"
  , "testing":        0
  , "implementation": 1
  , "note":           "artifact implementation."
}
```

```json
{
    "method":         "GET"
  , "uri":            "/messages:author"
  , "desc":           "Get Messages By Handle"
  , "testing":        0
  , "implementation": 1
  , "note":           "artifact implementation."
}
```

```json
{
    "method":         "POST"
  , "uri":            "/messages"
  , "desc":           "Create New Message"
  , "testing":        0
  , "implementation": 4
  , "note":           "Not tested, but code implemented."
}
```

```json
{
    "method":         "DELETE"
  , "uri":            "/messages"
  , "desc":           "Delete Message"
  , "testing":        0
  , "implementation": 1
  , "note":           "artifact implementation."
}
```

```json
{
    "method":         "POST"
  , "uri":            "/publish"
  , "desc":           "Publish Message"
  , "testing":        0
  , "implementation": 4
  , "note":           "Needs testing. Change to PATCH /users/:user, no verbs"
}
```

```json
{
    "method":         "POST"
  , "uri":            "/joindefault"
  , "desc":           "Join Default"
  , "testing":        8
  , "implementation": 4
  , "note":           "URI will be changed. No verbs."
}
```

```json
{
    "method":         "POST"
  , "uri":            "/join"
  , "desc":           "Join Default"
  , "testing":        8
  , "implementation": 8
  , "note":           "5 tests"
}
```

```json
{
    "method":         "POST"
  , "uri":            "/block"
  , "desc":           "Block User"
  , "testing":        8
  , "implementation": 8
  , "note":           ""
}
```

```json
{
    "method":         "POST"
  , "uri":            "/circles"
  , "desc":           "New Circle"
  , "testing":        8
  , "implementation": 8
  , "note":           ""
}
```
