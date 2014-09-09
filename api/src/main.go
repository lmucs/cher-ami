package main

import (
    "github.com/ant0ine/go-json-rest/rest"
    "log"
    "net/http"
)

func main() {

    handler := rest.ResourceHandler{
        EnableRelaxedContentType: true,
    }
    err := handler.SetRoutes(
        //&rest.Route{"GET", "/posts", GetAllPosts},
        //&rest.Route{"POST", "/posts", CreatePost},
        &rest.Route{"GET", "/posts/:id", GetPost},
        //&rest.Route{"DELETE", "/posts/:id", DeletePost},
    )
    if err != nil {
        log.Fatal(err)
    }
    log.Fatal(http.ListenAndServe(":8228", &handler))
}

type Post struct {
    Name string
    Content string
}

func GetPost(w rest.ResponseWriter, r *rest.Request) {
    id := r.PathParam("id")
    post := Post{"Zane", "This is just a sample post with id " + id}
    w.WriteJson(post)
}
