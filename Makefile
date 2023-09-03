# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
GOCLEAN = $(GOCMD) clean
GOGET = $(GOCMD) get
BINARY_PATH = bin
BINARY_NAME = ${BINARY_PATH}/server/server
APP_NAME = cmd/server/main.go

# Main build target
build:
	CGO_ENABLED=0 $(GOBUILD) -o $(BINARY_NAME) ${APP_NAME}

# Build docker image
build-docker:
	docker build -t storage .

# Run docker image
docker-run:
	docker run -it --entrypoint ./server storage -enable-log=true -period=10s

# Test target
test:
	$(GOTEST) -v ./... -cover

# Clean target
clean:
	rm -rf ./$(BINARY_PATH)

# Default target
default: build

.PHONY: build test clean