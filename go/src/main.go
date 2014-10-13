package main

import (
	a "./api"
	routes "./routes"
	"fmt"
	"github.com/jadengore/goconfig"
	"log"
	"net/http"
	"os"
)

func main() {
	c, err := goconfig.ReadConfigFile("../../config.cfg")
	var uri string
	if len(os.Args) > 1 {
		if os.Args[1] == "local" {
			fmt.Println("Local session requested.")
			uri, err = c.GetString("local-test", "url")
		}
	} else {
		uri, err = c.GetString("gen-test", "url")
	}
	api := a.NewApi(uri)
	handler, err := routes.MakeHandler(*api)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/api/", http.StripPrefix("/api", &handler))

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("../../web/src/"))))

	httpPort := "8228"
	fmt.Printf("Listening on port %s\n", httpPort)
	log.Fatal(http.ListenAndServe(":"+httpPort, nil))
}
