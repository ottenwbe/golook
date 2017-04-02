# golook

# The whole project is WIP. More info is coming soon #

[![Build Status](https://travis-ci.org/ottenwbe/golook.svg?branch=development)](https://travis-ci.org/ottenwbe/golook)
[![codecov](https://codecov.io/gh/ottenwbe/golook/branch/master/graph/badge.svg)](https://codecov.io/gh/ottenwbe/golook)

Search for files in a distributed system via Rest API. 

## Build ##

### Prerequisites for both, server and client ###

1. [Go](https://golang.org/doc/install) is installed and [GOPATH](https://golang.org/doc/code.html) is set.

1. Get dependencies
    ```bash
    go get ./.. 
    ```

### Build the server ###

1. Go build the executable
    ```bash    
    cd server
    go build ./.. -o bin/server
    ```

### Build the client ###

TBD

## Development ##

For development the dependencies are needed. 

Ginkgo for testing:

    go get github.com/onsi/ginkgo/ginkgo
    go get github.com/onsi/gomega
    
Govendor for vendoring:    
    
    go get github.com/kardianos/govendor

## Tests ##

### Unit tests ###
```bash
./testAllPackages.sh
```

## API ##

| Path  | Method  |   |   |   |
|---|---|---|---|---|
| /  | "GET" |   |   |   |
| /files/{file} |  "GET" |   |   |   |
| /systems/{system}/files  |  "GET" |   |   |   |
| /systems/{system}/files/{file}  |  "POST" |   |   |   |
| /systems/{system}/files  |  "PUT" |   |   |   |
| /systems/{system} |  "GET" |   |   |   |
| /systems |  "PUT" |   |   |   |
| /systems/{system} |  "DELETE" |   |   |   |
 
 
