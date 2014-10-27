# The CherAmi Server

Here is a server for CherAmi, exposing a REST API.  It's written in Go.

To learn about the API, look at the source code, and our always-in-progress [API docs on the CherAmi wiki](https://github.com/rtoal/cher-ami/wiki/API-Documentation).

##Setting Up the Server

Parts 1 through 3 are done only once, so skip to Part 4 if you have completed those before.

###Part 1: Put a secret config file in your project root

1. Ask @jadengore for the config file will all the secret urls and passwords
2. `cd` to the project root
3. Store the config file with name `config.cfg`.


###Part 1: Initialize $GOPATH

1. `cd` to your home folder
2. `mkdir go`
3. Add the following to your `.bash_profile`:

        # Setting PATH for Go
        export GOPATH=~/go
        export PATH=$PATH:$GOPATH/bin

###Part 2: Install or Update Dependencies

1. `cd` to the project root directory
2. `make install-deps` if you are building for the first time, or `make update-deps` if you need to refresh them.

###Part 3: Initialize test DB & Server

1. `cd` to the project root directory
2. `make start`

The API is now accessible on port 8228 locally.
