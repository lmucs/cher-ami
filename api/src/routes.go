//
// Template serving
//
package main

import (
    "flag"
    "html/template"
    "net/http"
    "fmt"
    "log"
    "path/filepath"
)

var (
    TMPL_DIR = filepath.Join("..", "..", "web", "html")
    signup   = "signup.html"
)

var (
    CSS_DIR   = filepath.Join("..", "..", "web", "css")
    signupcss = "signup.css"
)

var templates = template.Must(template.ParseFiles(
    filepath.Join(TMPL_DIR, signup),
    filepath.Join(CSS_DIR, signupcss),
))

func renderTemplate(w http.ResponseWriter, tmpl string) {
    err := templates.ExecuteTemplate(w, tmpl, nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, signup)
}

// func cssHandler(w http.ResponseWriter, r *http.Request) {
//     css := r.URL.Path[len("/css/"):]
// }

func main() {
    port := "8000"
    
    flag.Parse()
    http.HandleFunc("/signup", signupHandler)
    http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("../../web/css"))))

    fmt.Printf("Listening on port %s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

