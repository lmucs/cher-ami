package main

import (
    "github.com/ant0ine/go-json-rest/rest"
    "log"
    "net/http"
    "gopkg.in/mgo.v2"
)

func main() {

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
        //&rest.Route{"GET", "/posts", GetAllPosts},
        &rest.Route{"POST", "/posts", api.CreatePost},
        &rest.Route{"GET", "/posts/:id", api.GetPost},
        //&rest.Route{"DELETE", "/posts/:id", DeletePost},
    )
    if err != nil {
        log.Fatal(err)
    }

    log.Fatal(http.ListenAndServe(":8228", &handler))
}

type Api struct {
    session *mgo.Session
    db      *mgo.Database // main db
}

type Post struct {
    Name string
    Content string
}

func (a Api) CreatePost(w rest.ResponseWriter, r *rest.Request) {
    post := Post{"A pre-created post", "Some pre-created body text"}
    err := a.db.C("posts").Insert(post)
    if err != nil {
        log.Fatal("Can't insert document: %v\n", err)
    }
}

func (a Api) GetPost(w rest.ResponseWriter, r *rest.Request) {
    id := r.PathParam("id")
    post := Post{"Zane", "This is just a sample post with id " + id}
    w.WriteJson(post)
}
