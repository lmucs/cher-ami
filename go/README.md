# The CherAmi Server

Here is a server for CherAmi, exposing a REST API.  It's written in Go.

To learn about the API, look at the source code, and our always-in-progress [API docs on the CherAmi wiki](https://github.com/rtoal/cher-ami/wiki/API-Documentation).

##Setting Up the Server

Parts 1 and 2 are done only once, so skip to Part 3 if you have completed those before.

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
2. Install the dependencies. Since the packages we are using are changing frequently, please refer to the source for packages necessary. To install, run `go get <PACKAGE_URL>`.
   
###Part 3: Initialize test DB & Server

1. In the `cher-ami` directory:
    `make start`

The API is now accessible on port 8228 locally.
