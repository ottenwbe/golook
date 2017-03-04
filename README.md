# golook

# The whole project is WIP. More info is coming soon #

[![Build Status](https://travis-ci.org/ottenwbe/golook.svg?branch=development)](https://travis-ci.org/ottenwbe/golook)
[![codecov](https://codecov.io/gh/ottenwbe/golook/branch/master/graph/badge.svg)](https://codecov.io/gh/ottenwbe/golook)

Search for files in a distributed system via Rest API.

## Build ##

### Prerequisites for both, server and client ###

1. Assumption: go is installed and GOPATH is set.

1. Get dependencies
    ```sh
    go get ./.. 
    ```

### Build the client ###

1. Go build the executable
    ```sh
    
    cd server
    go build ./.. -o bin/server
    ```

### Build the server ###

TBD

### Execute tests ###
```sh
./testAllPackages.sh
```

