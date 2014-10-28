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

    {
        "method":         "POST"
      , "uri":            "/signup"
      , "desc":           "Sign up new user"
      , "testing":        8
      , "implementation": 8
      , "note":           ""
    }

    {
        "method":         "POST"
      , "uri":            "/changepassword"
      , "desc":           "Change password"
      , "testing":        8
      , "implementation": 4
      , "note":           "route needs change"
    }

    {
        "method":         "POST"
      , "uri":            "/sessions"
      , "desc":           "Login (new Auth token)"
      , "testing":        8
      , "implementation": 8
      , "note":           ""
    }

    {
        "method":         "DELETE"
      , "uri":            "/sessions"
      , "desc":           "Logout (delete Auth token)"
      , "testing":        4
      , "implementation": 8
      , "note":           ""
    }

    {
        "method":         "GET"
      , "uri":            "/users/:handle"
      , "desc":           "Get single user"
      , "testing":        1
      , "implementation": 2
      , "note":           "route doesn't follow specs, takes in query params"
    }


    {
        "method": "GET"
      , "uri":    "/users"
      , "desc":   "Search For Users"
      , "testing":        2
      , "implementation": 4
      , "note":           "needs a thorough testing to 8-8"
    }


    {
        "method":         "DELETE"
      , "uri":            "/users/:handle"
      , "desc":           "Delete User"
      , "testing":        4
      , "implementation": 4
      , "note":           "needs GET /user for complete testing"
    }


    {
        "method":         "GET"
      , "uri":            "/messages"
      , "desc":           "Get Authored Messages"
      , "testing":        0
      , "implementation": 1
      , "note":           "artifact implementation."
    }

    {
        "method":         "GET"
      , "uri":            "/messages:author"
      , "desc":           "Get Messages By Handle"
      , "testing":        0
      , "implementation": 1
      , "note":           "artifact implementation."
    }

    {
        "method":         "POST"
      , "uri":            "/messages"
      , "desc":           "Create New Message"
      , "testing":        0
      , "implementation": 4
      , "note":           "Not tested, but code implemented."
    }

    {
        "method":         "DELETE"
      , "uri":            "/messages"
      , "desc":           "Delete Message"
      , "testing":        0
      , "implementation": 1
      , "note":           "artifact implementation."
    }

    {
        "method":         "POST"
      , "uri":            "/publish"
      , "desc":           "Publish Message"
      , "testing":        0
      , "implementation": 4
      , "note":           "Needs testing."
    }

    {
        "method":         "POST"
      , "uri":            "/joindefault"
      , "desc":           "Join Default"
      , "testing":        8
      , "implementation": 4
      , "note":           "URI will be changed"
    }

    {
        "method":         "POST"
      , "uri":            "/join"
      , "desc":           "Join Default"
      , "testing":        8
      , "implementation": 8
      , "note":           "5 tests"
    }

    {
        "method":         "POST"
      , "uri":            "/block"
      , "desc":           "Block User"
      , "testing":        8
      , "implementation": 8
      , "note":           ""
    }

    {
        "method":         "POST"
      , "uri":            "/circles"
      , "desc":           "New Circle"
      , "testing":        8
      , "implementation": 8
      , "note":           ""
    }


