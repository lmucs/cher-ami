package main

import (
	a "./api"
	routes "./routes"
	"fmt"
	"github.com/jmcvetta/neoism"
	"log"
	"net/http"
	"os"
)

func main() {
	args := os.Args[1:]
	port := "8228"

	uri := args[0]
	neo4jdb, err := neoism.Connect(uri)
	if err != nil {
		log.Fatal(err)
	}

	api := a.NewApi(neo4jdb)
	handler, err := routes.MakeHandler(*api)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/api/", http.StripPrefix("/api", &handler))

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("../../web/src/"))))

	fmt.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
	//log.Fatal(http.ListenAndServe(":"+port, &handler))
}
