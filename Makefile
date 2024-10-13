# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get

# Build target
build:
	$(GOBUILD) -o resume-generator

# Clean target
clean:
	$(GOCLEAN)
	rm -f resume-generator

# Test target
test:
	$(GOTEST) -v ./...

# Get dependencies target
deps:
	$(GOGET) -v ./...

.PHONY: build clean test deps
