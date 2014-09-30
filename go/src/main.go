package main

import (
	cheramiapi "./api"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jmcvetta/neoism"
	"log"
	"net/http"
	"os"
)

func main() {
	args := os.Args[1:]
	port := "8228"
	handler := rest.ResourceHandler{
		EnableRelaxedContentType: true,
	}

	uri := args[0]
	neo4jdb, err := neoism.Connect(uri)
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
		&rest.Route{"GET", "/messages/:handle", api.GetMessagesByHandle},
		&rest.Route{"POST", "/messages", api.NewMessage},
		&rest.Route{"DELETE", "/messages", api.DeleteMessage},
		&rest.Route{"POST", "/publish", api.PublishMessage},
	)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/api/", http.StripPrefix("/api", &handler))

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("../../web/src/"))))

	fmt.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
	//log.Fatal(http.ListenAndServe(":"+port, &handler))
}
