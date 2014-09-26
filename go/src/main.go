package main

import (
	cheramiapi "./api"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jmcvetta/neoism"
	"log"
	"net/http"
)

func main() {
	port := "8228"
	handler := rest.ResourceHandler{
		EnableRelaxedContentType: true,
	}

	neo4jdb, err := neoism.Connect("http://localhost:7474/db/data")
	if err != nil {
		log.Fatal(err)
	}

	api := &cheramiapi.Api{neo4jdb}

	err = handler.SetRoutes(
		&rest.Route{"POST", "/signup", api.Signup},
		&rest.Route{"POST", "/login", api.Login},
		&rest.Route{"POST", "/logout", api.Logout},
		&rest.Route{"GET", "/users", api.GetUser},
		&rest.Route{"DELETE", "/users", api.DeleteUser},
		&rest.Route{"GET", "/messages", api.GetAuthoredMessages},
		&rest.Route{"GET", "/messages/:handle", GetMessagesByHandle},
		&rest.Route{"POST", "/messages", api.NewMessage},
		&rest.Route{"DELETE", "/messages", api.DeleteMessage},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, &handler))
}
