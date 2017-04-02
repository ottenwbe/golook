# golook

# The whole project is WIP. More info is coming soon #

[![Build Status](https://travis-ci.org/ottenwbe/golook.svg?branch=development)](https://travis-ci.org/ottenwbe/golook)
[![codecov](https://codecov.io/gh/ottenwbe/golook/branch/master/graph/badge.svg)](https://codecov.io/gh/ottenwbe/golook)

Golook is a cli tool that allows users to search for files in a distributed system.
To this end, the application can act as a client and server.
Clients report and query a server for files on a specific systems.
Servers cache the reports of multiple clients and can therefore answer queries about the location of files.
Client and server communicate via an Rest API. 

## Install ##

1. Ensure that [Go](https://golang.org/doc/install) is installed and [GOPATH](https://golang.org/doc/code.html) is set. 
Moreover `$GOPATH/bin` is in your `PATH`.

1. Go build the executable in `${GOPATH}/bin`:
    
    ```bash    
    go get github.com/ottenwbe/golook
    ```
1. Now you can execute the application by typing: 

    ```bash    
    golook --help
    ```

## Development ##

For development the following dependencies are needed. 

[Ginkgo](https://onsi.github.io/ginkgo/) is required for testing:

    go get github.com/onsi/ginkgo/ginkgo
    go get github.com/onsi/gomega
    
[Govendor](https://github.com/kardianos/govendor) is used for vendoring:    
    
    go get github.com/kardianos/govendor

### Tests ###

#### Unit tests ####

The execution of unit tests is simplified by a script:

    ./testAllPackages.sh

#### Integration tests ####

TODO: Will be executed by executing client and server in Docker containers.

## Server API ##

| Path  | Method  | Purpose  |   
|---|---|---|
| /  | "GET" | Returns the current version of the server  |   
| /files/{file} |  "GET" |  Get all systems that host a specific file |  
| /systems/{system}/files  |  "GET" | Get all files of a system | 
| /systems/{system}/files/{file}  |  "POST" | Report a distinct file |
| /systems/{system}/files  |  "PUT" | Replace all reported files |
| /systems/{system} |  "GET" | Get details about a system  |
| /systems |  "PUT" |  Report a new system |
| /systems/{system} |  "DELETE" |  Delete a system |
 
 
