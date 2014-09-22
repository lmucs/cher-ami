package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jmcvetta/neoism"
	//"github.com/gorilla/schema"
	"fmt"
	"log"
	"net/http"
	"src/api"
	"time"
)

func main() {
	port := "8228"
	handler := rest.ResourceHandler{
		EnableRelaxedContentType: true,
	}

	db, err := neoism.Connect("http://localhost:7474/db/data")
	if err != nil {
		log.Fatal(err)
	}

	api := Api{db}

	err = handler.SetRoutes(
		&rest.Route{"POST", "/signup", api.Signup},
		&rest.Route{"POST", "/login", api.Login},
		&rest.Route{"GET", "/users", api.GetUser},
		&rest.Route{"DELETE", "/users", api.DeleteUser},
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
