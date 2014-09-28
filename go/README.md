# CherAmi API Documentation

Aside from the source, check out our evolving [API docs on the CherAmi wiki](https://github.com/rtoal/cher-ami/wiki/API-Documentation).

##To set up the Go Server on Mac:
Part 1 and 2 are a one-time step, skip to Part 3 if you have completed those before.

###Part 1: Initialize $GOPATH

1. `cd` to your home folder
2. `mkdir go`
3. `open .bash_profile`
4. Add the following to your `.bash_profile`:

        # Setting PATH for Go  
        export GOPATH=~/go
        export PATH=$PATH:$GOPATH/bin

###Part 2: Install Dependencies
1. `cd` to the `cher-ami/go/src` directory
2. Install the dependencies. Since the packages we are using are changing rapidly, please refer to the source for packages necessary. To install, run `go get <PACKAGE_URL>`.
   
###Part 3: Initialize test DB & Server
1. If you haven't install mongodb, do so. It must be running in order to access the API (`mongod`).
2. Start up the server:
    `go run main.go`

The API is now accessible on port 8228 locally.
