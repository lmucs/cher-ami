# CherAmi API Documentation

A more developed API will be posted here in the future.

Aside from the source, check out our evolving [API docs on the CherAmi wiki](https://github.com/rtoal/cher-ami/wiki/API-Documentation).

To set up the Go Server Locally:
    1. `cd` to your home folder
    2. `mkdir go`
    3. `open .bash_profile`
    4. add the following to your `.bash_profile` file:
    ```# Setting PATH for Go
    export GOPATH=~/go
    export PATH=$PATH:$GOPATH/bin```
5. `cd` to the `cher-ami/api/src` directory
6. Install the dependencies (this is the current list):
    ```go get github.com/ant0ine/go-json-rest/rest
    go get gopkg.in/mgo.v2```
7. Build the Go main file:
    `go build main.go`
8. If you haven't install mongodb, do so. It must be running in order to access the API.
9. Start up the server:
    `go run main.go`

YOU CAN NOW ACCESS THE API ON PORT 8228!
