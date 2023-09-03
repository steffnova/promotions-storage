# Promotions Storage 

Design storage for promotions. Storage should have an HTTP endpoint to access promotion given
it's ID. Promotions are added to storage periodically using CSV file.
When designing the storage following things should be taken into consideration:
 - The .csv file could be very big (billions of entries) 
 - Every new file is immutable, that is, whole storage should be erased and rewritten again
 - Consider a way that it will be deployed into production (deployment, scaling, monitoring)
 - Operating under peak periods (millions of requests per minute)


## Quick Start

To build and test the project use provided Makefile
```
# Build server and outputs it to bin/server/server
make build
# Test all packages and report test coverage
make test
```

To build and run docker image use:
```
# Build docker file
make build-docker
# Run docker image
make docker-run
```

Once the project is built it can be run. To check available flags run:
```
## Run help to see usage
bin/server/server --help

Usage of server:
  -enable-log
        enables/disables logging
  -file-path string
        specify path to csv file from which promotions will be loaded (default "promotions.csv")
  -period duration
        period between promotion storage updates (default 1s)
```
To run the server with specific flags:
```
## with flag values
bin/server/server -enable-log true -period=10s -file-path="promotions.csv"
```



## Solution

File processing is done using go channels which allows processing of CSV file line by line and sending read data to go channel.
This ensures that file data is not read into memory all at once.

The server has been implemented with In-memory storage, for simplicy reasons. storage package offers an interface that would
allow to easily replace In-memory implementation with different type of storing mechanism (NoSQL, SQL, etc...)

Current solution is horizontally scalable. This means that increasing the number of server applications the number
of requests that can be served could be increased. the only caveat is that each server instance needs to have access to the same file.

Storage is periodically updated, consuming the same file. File content can be replaced to simulate changes during periodic update.

Server can be deployed as docker container.

Server is serving http traffic on port 8080.