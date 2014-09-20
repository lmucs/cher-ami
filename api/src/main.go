package main

import (
    "github.com/ant0ine/go-json-rest/rest"
    "github.com/jmcvetta/neoism"
    //"github.com/gorilla/schema"
    "log"
    "net/http"
    "fmt"
    "time"
)

func panicErr(err error) {
    if err != nil {
        panic(err)
    }
}

func httpError(w rest.ResponseWriter, err error) {
    if err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func main() {
    port    := "8228"
    handler := rest.ResourceHandler{
        EnableRelaxedContentType: true,
    }

    db, err := neoism.Connect("http://localhost:7474/db/data")
    if  err != nil {
        log.Fatal(err)
    }

    api      := Api{db}

    err = handler.SetRoutes(
        &rest.Route{"POST",   "/signup", api.Signup},
        // &rest.Route{"POST",   "/login", api.Login},
        &rest.Route{"GET",    "/users", api.GetUser},
        // &rest.Route{"DELETE", "/users/:id", api.DeleteUser},
        // &rest.Route{"GET",  "/message", GetAllMessages},
        // &rest.Route{"POST",   "/messages", api.CreateMessage},
        // &rest.Route{"GET",    "/messages/:id", api.GetMessage},
        // &rest.Route{"DELETE", "/messages/:id", api.DeleteMessage},
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Listening on port %s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, &handler))
}

//
// Application Types
//

type Api struct {
    db   *neoism.Database
}

//
// Data types
// All data types are stored in mongodb,
// this gives them an '_id' identifier
//

/*type Message struct {
    Id         bson.ObjectId
    Owner      bson.ObjectId
    Created    time.Time
    Content    string
    ResponseTo bson.ObjectId
    RepostOf   bson.ObjectId
    Circles    []bson.ObjectId
}

type Circle struct {
    Owner      bson.ObjectId
    Members    []bson.ObjectId
    Name       string
}*/

type UserProposal struct {
    Handle          string
    Email           string
    Password        string
    ConfirmPassword string
}

type UserView struct {
    Handle       string
    Joined       time.Time
}

type UserSignIn struct {
    Handle     string
    Password   string
}

//
// API
//

/*
 * Expects a json POST with "Username", "Password", "ConfirmPassword"
 */
func (a Api) Signup(w rest.ResponseWriter, r *rest.Request) {
    proposal := UserProposal{}
    err      := r.DecodeJsonPayload(&proposal)
    if err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    stmt := `CREATE (user:User { handle:{handle}, password:{password}, joined: {joined} })
             RETURN user`

    res := []struct {
        User   neoism.Node
    }{}

    params := neoism.Props{
        "handle": proposal.Handle,
        "password": proposal.Password,
        "joined": time.Now().Local(),
    }

    cq := neoism.CypherQuery {
        Statement: stmt,
        Parameters: params,
        Result: &res,
    }

    err = a.db.Cypher(&cq)
    panicErr(err)

    // check results
    if len(res) != 1 {
        panic(fmt.Sprintf("Incorrect results len in query1()\n\tgot %d, expected 1\n", len(res)))
    }

    n := res[0].User // Only one row of data returned
    fmt.Println("createNode()", n.Data)

    // ensure unique handle
    // count, err := a.db.C("users").Find(bson.M{ "handle": proposal.Handle }).Count()
    // if count > 0 {
    //     rest.Error(w, proposal.Handle+" is already taken", 400)
    //     return
    // }
    // if err != nil {
    //     rest.Error(w, err.Error(), http.StatusInternalServerError)
    //     return
    // }

    // password checks
    // if proposal.Password != proposal.ConfirmPassword {
    //     rest.Error(w, "Passwords do not match", 400)
    //     return
    // }
}

/* Transition to PostgreSQL in progress */

// func (a Api) Login(w rest.ResponseWriter, r *rest.Request) {
//     credentials := UserSignIn{}

//     err := r.DecodeJsonPayload(&credentials)
//     if err != nil {
//         rest.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     result := User{}
//     err = a.db.C("users").
//             Find(bson.M{"handle": credentials.Handle, "password": credentials.Password}).
//             One(&result)
//     if err != nil {
//         rest.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
// }

func (a Api) GetUser(w rest.ResponseWriter, r *rest.Request) {
    querymap   := r.URL.Query()

    // Get by handle
    if handle, ok := querymap["handle"]; ok {

        stmt := `MATCH (user:User)
                 WHERE user.handle = {handle}
                 RETURN user`
        params := neoism.Props{
            "handle": handle[0],
        }
        res := []struct{
            User   neoism.Node
        }{}

        cq := neoism.CypherQuery{
            Statement:   stmt,
            Parameters:  params,
            Result:      &res,
        }

        err := a.db.Cypher(&cq)
        panicErr(err)

        u := res[0].User.Data

        w.WriteJson(u)
        return
    }

    // // Get by id
    // if id, ok  := querymap["id"]; ok {
    //     found  := User{}
    //     err    := a.db.C("users").
    //                 Find(bson.M{ "id": bson.ObjectIdHex(id[0]) }).
    //                 One(&found)
    //     if err != nil {
    //         rest.Error(w, err.Error(), http.StatusInternalServerError)
    //         return
    //     }
    //     w.WriteJson(found)
    //     return
    // }


    // // All users
    // var users []interface{}

    // a.db.C("users").
    //     Find(bson.M{}).
    //     Select(bson.M{ "handle":1 }).
    //     All(&users)

}

// func (a Api) DeleteUser(w rest.ResponseWriter, r *rest.Request) {
//     bid := bson.ObjectIdHex(r.PathParam("id"))

//     err := a.db.C("users").Remove(bson.M{"id": bid})
//     if err != nil {
//         rest.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
// }

// func (a Api) CreateMessage(w rest.ResponseWriter, r *rest.Request) {
//     message := Message{
//         bson.NewObjectId(),
//         bson.NewObjectId(),     // owner ID
//         time.Now().Local(),
//         "",                     // content
//         NIL_ID,
//         NIL_ID,
//         []bson.ObjectId{},
//     }

//     payload := Message{}
//     err     := r.DecodeJsonPayload(&payload)
//     if err != nil {
//         rest.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
//     message.Content = payload.Content

//     if message.Content == "" {
//         rest.Error(w, "please enter some content for your message", 400)
//         return
//     }

//     err = a.db.C("messages").Insert(message)
//     if err != nil {
//         log.Fatal("Can't insert document: %v\n", err)
//     }
// }

// func (a Api) GetMessage(w rest.ResponseWriter, r *rest.Request) {
//     bid     := bson.ObjectIdHex(r.PathParam("id"))
//     message := Message{}
//     err     := a.db.C("messages").Find(bson.M{"id": bid}).One(&message)
//     if err != nil {
//         rest.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
//     w.WriteJson(message)
// }

// func (a Api) DeleteMessage(w rest.ResponseWriter, r *rest.Request) {
//     bid := bson.ObjectIdHex(r.PathParam("id"))
//     err := a.db.C("messages").Remove(bson.M{"id": bid})
//     if err != nil {
//         rest.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
// }
