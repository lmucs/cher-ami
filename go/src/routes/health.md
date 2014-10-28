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
      , "desc":           "DeleteUser"
      , "testing":        4
      , "implementation": 4
      , "note":           "needs GET /user for complete testing"
    }


    {
        "method":         "GET"
      , "uri":            "/messages"
      , "desc":           "GetAuthoredMessages"
      , "testing":        0
      , "implementation": 1
      , "note":           "artifact implementation. "
    }






"GET", "/messages/:author", a.GetMessagesByHandle},
"POST", "/messages", a.NewMessage},
"DELETE", "/messages", a.DeleteMessage},
"POST", "/publish", a.PublishMessage},
"POST", "/joindefault", a.JoinDefault},
"POST", "/join", a.Join},
"POST", "/block", a.BlockUser},
"POST", "/circles", a.NewCircle},