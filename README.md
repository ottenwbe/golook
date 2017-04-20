# golook

# The whole project is WIP---Work in Progress. More info is coming soon 

[![Build Status](https://travis-ci.org/ottenwbe/golook.svg?branch=development)](https://travis-ci.org/ottenwbe/golook)
[![codecov](https://codecov.io/gh/ottenwbe/golook/branch/master/graph/badge.svg)](https://codecov.io/gh/ottenwbe/golook)

Golook is a broker backend that allows users to search for files in a distributed system, i.e., your LAN.
To this end, the application can act as a client and server.
Clients report and query a server for files on a specific systems.
Servers cache the reports of multiple clients and can therefore answer queries about the location of files.
Client and server communicate via an Rest API. 

## Background ##

If you are like me, you have multiple machines connected in your local network, i.e., laptops, servers, nas, raspi etc.
Many of the files on these machines are versioned in a (git) repository, tracked by a configuration management system, or are backed up. 
However, from time to time I wonder where a specific file is, i.e., where I downloaded the latest Linux Image.
At this point, I recalled the distributed file search algorithms, i.e., Chord, Tapestry, and Can.

So I decided to implement a simple distributed file search. And as a benefit I have the opportunity to learn go.   

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

### Structure ###
    .
    ├── doc                 // main documenation node    
    ├── broker              // main source node
    │   ├── cmd             // commands and configuration    
    │   ├── api             // http api of the golook broker (e.g., towards the golook cli)
    │   ├── communication   // communication substrate (e.g., json rpc communication)
    │   ├── routing         // routing layer     
    │   ├── management      // management layer 
    │   ├── repository      // repository 
    │   ├── runtime         // runtime components
    │   └── models          // common models for all layers
    └── test                // docker files for (integration) testing
        └── integration

### Tests ###

#### Unit tests ####

The execution of unit tests is simplified by a script:

    ./testAllPackages.sh

#### Integration tests ####

TODO: Will be implemented by executing client and server in Docker containers.

## Broker API ##

| Path  | Method  | Purpose  |   
|---|---|---|
| /info  | "GET" | Returns the current version of the server  |   
| /file/{file} |  "GET" |  Get all systems that host a specific file |  
| /file |  "PUT" | The broker should report a file |
| /folder |  "PUT" | The broker should report all files in a folder  |
