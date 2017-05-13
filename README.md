# golook

# The whole project is WIP---Work in Progress. More info is coming soon 

[![Build Status](https://travis-ci.org/ottenwbe/golook.svg?branch=development)](https://travis-ci.org/ottenwbe/golook)
[![codecov](https://codecov.io/gh/ottenwbe/golook/branch/master/graph/badge.svg)](https://codecov.io/gh/ottenwbe/golook)

Golook is a middleware that allows users to search for files in a distributed system, i.e., your LAN.
To this end, the application can act as a client and server, a servant.
Clients report to and query from a server files and folders.
Servers cache the reports of multiple clients and can therefore answer queries about the location of files.
 

## Background ##

If you are like me, you have multiple machines connected in your local network, i.e., laptops, servers, nas, raspberry pi etc.
Many of the files on these machines are versioned in a (git) repository, tracked by a configuration management system, or are backed up. 
However, from time to time I wonder where a specific file is, e.g., 'where did I download the latest Linux Image?'.
At this point, I recalled the distributed file search algorithms, i.e., Chord, Tapestry, and Can.

So I decided to implement a simple distributed file search. And as a benefit I have the opportunity to learn go.   

## Install ##

1. Ensure that [Go](https://golang.org/doc/install) is installed and [GOPATH](https://golang.org/doc/code.html) is set. 
Moreover `$GOPATH/bin` is in your `PATH`.

1. Go get and build the executable in `${GOPATH}/bin`:
    
    ```bash    
    go get github.com/ottenwbe/golook
    ```
1. Now you can execute the application by typing: 

    ```bash    
    golook --help
    ```

### Interactions ###

There are basically two ways to interact with golook: the golook-cli and the golook (broker) api

### Golook CLI ###

For details see [golook cli](https://github.com/ottenwbe/golook-cli).

### Golook Broker API ###

| Path  | Method  | Purpose  |   
|---|---|---|
| /info  | "GET" | Returns information like the current version of the golook broker  |   
| /file/{file} |  "GET" |  Get all systems that host a specific file |  
| /file |  "PUT" | The broker should report a file |
| /folder |  "PUT" | The broker should report all files in a folder  |


## Development ##

For development the following dependencies are needed. 

[Ginkgo](https://onsi.github.io/ginkgo/) is required for testing:

    go get github.com/onsi/ginkgo/ginkgo
    go get github.com/onsi/gomega
    
[Govendor](https://github.com/kardianos/govendor) is used for vendoring:    
    
    go get github.com/kardianos/govendor

### Structure ###

    ├── broker              // main source node for the broker
    │   ├── api             // http api of the golook broker (e.g., towards the golook cli)
    │   ├── communication   // communication substrate (e.g., json rpc communication)
    │   ├── models          // common models for all layers
    │   ├── repository      // repository for caching and persisting data
    │   ├── routing         // routing layer
    │   ├── runtime         // runtime components
    │   │   ├── cmd         // commands and configuration
    │   │   └── core        // runtime components
    │   └── service         // service layer
    ├── client              // source of the client for the http api
    ├── doc                 // main documenation node    
    ├── test
    │   ├── integration     // scripts and docker files for (integration) testing
    │   └── unit            // scripts to perorm unit tests
    ├── utils               // source code of utils for all packages
    └── vendor              // vendored packages

### Tests ###

#### Unit tests ####

The execution of unit tests is simplified by a script:

    test/unit/run_unit_tests.sh
    
or
    
    make test

#### Integration tests ####

Are executed by starting one or more peers in a docker and interacting with those peers. 

    test/integration/run_integration_tests.sh
