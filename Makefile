GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=seldondm
PATH_TO_MAIN=cmd/main.go

all: test build
build:
	$(GOBUILD) -o $(BINARY_NAME) $(PATH_TO_MAIN)
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
run: build
	./$(BINARY_NAME) model.json
