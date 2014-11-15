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
	config, err := goconfig.ReadConfigFile("../config.cfg")
	port, err := config.GetString("default", "server-port")
	var uri string
	if len(os.Args) > 1 && os.Args[1] == "local" {
		uri, err = config.GetString("local-test", "url")
	} else {
		uri, err = config.GetString("gen-test", "url")
	}
	api := a.NewApi(uri)
	handler, err := routes.MakeHandler(*api, false)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/api/", http.StripPrefix("/api", &handler))

	http.Handle("/api/docs/", http.StripPrefix("/api/docs", http.FileServer(http.Dir("docs/"))))

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("../web/src/"))))

	fmt.Printf("The CherAmi server is listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}
