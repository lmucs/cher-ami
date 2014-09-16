package main

import (
    "github.com/ant0ine/go-json-rest/rest"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "log"
    "net/http"
    "fmt"
    "time"
)

func main() {
    port := "8228"
    handler := rest.ResourceHandler{
        EnableRelaxedContentType: true,
    }

    session, err := mgo.Dial("mongodb://localhost")
    if err != nil {
        log.Fatal(err)
    }

    database := session.DB("cher-ami")
    api := Api{session, database}

    err = handler.SetRoutes(
        &rest.Route{"POST", "/signup", api.Signup},
        &rest.Route{"POST", "/login", api.Login},
        &rest.Route{"GET", "/users/delete/:id", api.DeleteUser},
        //&rest.Route{"GET", "/message", GetAllMessages},
        &rest.Route{"POST", "/messages", api.CreateMessage},
        &rest.Route{"GET", "/messages/:id", api.GetMessage},
        //&rest.Route{"DELETE", "/message/:id", DeleteMessage},
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
    session *mgo.Session
    db      *mgo.Database // main db
}

//
// Data types
// All data types are stored in mongodb,
// this gives them an '_id' identifier
//

type (
    UserId    string
    MessageId string
    CircleId  string
)

type Message struct {
    Owner      UserId
    Created    time.Time
    Content    string
    ResponseTo MessageId    // "" if not response
    RepostOf   MessageId    // "" if not repost
    Circles    []CircleId
}

type UserProposal struct {
    Handle          string
    Password        string
    ConfirmPassword string
}

type User struct {
    Handle       string
    Password     string
    Joined       time.Time
    Follows      []UserId
    BlockedUsers []UserId
}

type UserSignIn struct {
    Handle       string
    Password     string
}

//
// API
//

func (a Api) Signup(w rest.ResponseWriter, r *rest.Request) {
    proposal := UserProposal{}

    // expects a json POST with "Username", "Password", "ConfirmPassword"
    err := r.DecodeJsonPayload(&proposal)
    if err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // ensure unique handle
    count, err := a.db.C("users").Find(bson.M{ "handle": proposal.Handle }).Count()
    if count > 0 {
        rest.Error(w, proposal.Handle+" is already taken", 400)
        return
    }
    if err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // password checks
    if proposal.Password != proposal.ConfirmPassword {
        rest.Error(w, "Passwords do not match", 400)
        return
    }

    user := User{
        proposal.Handle,
        proposal.Password,  // plaintext for now
        time.Now().Local(),
        []UserId{},
        []UserId{},
    }
    err = a.db.C("users").Insert(user)
    if err != nil {
        log.Fatal("Can't insert user: %v\n", err)
    }
}

func (a Api) Login(w rest.ResponseWriter, r *rest.Request) {
    credentials := UserSignIn{}

    err := r.DecodeJsonPayload(&credentials)
    if err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    result := User{}
    err = a.db.C("users").Find(bson.M{"handle": credentials.Handle, "password": credentials.Password}).One(&result)
    if err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func (a Api) DeleteUser(w rest.ResponseWriter, r *rest.Request) {
    bid := bson.ObjectIdHex(r.PathParam("id"))

    err := a.db.C("users").Remove(bson.M{"_id": bid})
    if err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func (a Api) CreateMessage(w rest.ResponseWriter, r *rest.Request) {
    message := Message{}

    // expects a json POST with Message properties, case-sensitive
    // use custom strings as placeholders for testing. This will be
    // remidied as the user api grows.
    err := r.DecodeJsonPayload(&message)
    if err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if message.Content == "" {
        rest.Error(w, "please enter some content for your message", 400)
        return
    }
    message.Created = time.Now().Local()
    err = a.db.C("messages").Insert(message)
    if err != nil {
        log.Fatal("Can't insert document: %v\n", err)
    }
}

func (a Api) GetMessage(w rest.ResponseWriter, r *rest.Request) {
    var id UserId = UserId(r.PathParam("id"))
    // sample
    message := Message{
        id,
        time.Now().Local(),
        "This is a sample message, ayeee",
        "",
        "",
        []CircleId{"c_777", "c_w0qweq45", "c_888282"},
    }
    w.WriteJson(message)
}
