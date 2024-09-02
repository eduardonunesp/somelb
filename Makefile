# Makefile for building and running the Go project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=server

# Build the project
build:
	$(GOBUILD) -o $(BINARY_NAME)

# Run the project
run: build
	./$(BINARY_NAME)

# Clean the project
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Test the project
test:
	$(GOTEST) -v ./...

# Install dependencies
deps:
	$(GOGET) -v ./...

.PHONY: build clean run test deps