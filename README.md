# golook broker

## Disclaimer 

The whole project is in an early alpha stage and not recommended for production.

## Introduction 

[![Build Status](https://travis-ci.org/ottenwbe/golook.svg?branch=development)](https://travis-ci.org/ottenwbe/golook)
[![codecov](https://codecov.io/gh/ottenwbe/golook/branch/master/graph/badge.svg)](https://codecov.io/gh/ottenwbe/golook)


If you are like me, you have multiple machines connected in your local network, i.e., laptops, servers, nas, raspberry pi etc.
Many a files on these machines are versioned in a (git) repository, tracked by a configuration management system, or are backed up. 
However, from time to time I wonder where a specific file is, e.g., 'where did I edit my new profile picture?'.
At one point, while looking for a file  , I recalled the distributed file search algorithms I studied at university, i.e., Chord and Can.

So I decided to implement a very simple distributed file search application based on broadcasts, denoted _golook_. And as a benefit I have the opportunity to learn go.

Golook is a middleware that allows users to search for files in a distributed system, i.e., your LAN.
To this end, Golook spans an overlay over all nodes in the distributed system. This allows all nodes to send messages to and receive messages from each other.
Golook acts as a client and server, a broker.
Clients report to and query from other brokers files and folders.
Servers cache the reports of multiple clients and can therefore answer queries about the location of files.
 
   

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

### Usage ###

There are basically two ways to interact with golook: the golook-cli and the golook (broker) api

### Golook CLI ###

For details see [golook cli](https://github.com/ottenwbe/golook-cli).

### Golook Broker API ###

| Path  | Method  | Purpose  |   
|---|---|---|
| /info  | "GET" | Returns information like the current version of the golook broker  |   
| /api  | "GET" | Returns all API endpoints  |
| /log  | "GET" | Returns the complete log of the golook server  |
| /v1/config  | "GET" | Returns the current configuration of the golook server |
| /v1/system  | "GET" | Returns information about the runtime environment |
| /v1/file/{file} |  "GET" |  Get all systems that host a specific file |  
| /v1/file |  "PUT" | The broker should report a file or folder |


## Development ##

Details:
* [Architecture](doc/Architecture.md)

### Prerequisites ###

For development the following dependencies are needed. 

[Ginkgo](https://onsi.github.io/ginkgo/) is required for testing:

    go get github.com/onsi/ginkgo/ginkgo
    go get github.com/onsi/gomega
    
[Govendor](https://github.com/kardianos/govendor) is used for vendoring:    
    
    go get github.com/kardianos/govendor

### Structure ###

    ├── broker              // main node for the broker's sources 
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

    make test
    
or
    
    test/unit/run_unit_tests.sh
    

#### Integration tests ####

Are executed by starting one or more peers in a docker and interacting with those peers.
 
    make integraton
    
Or you can build the docker container manually and then start the tests.

    docker build --rm=true --file=test/integration/Dockerfile --tag=golook:latest .    
    test/integration/run_integration_tests.sh
