# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_x64=cpd_api_amd64
BINARY_RPI=cpd_api_arm

all: test build
build:
	$(GOBUILD) -o $(BINARY_x64) -v
# test:
# 	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_x64)
	rm -f $(BINARY_RPI)
# deps:
#	$(GOGET) url

# Cross compilation
linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v
rpi:
	GOOS=linux GOARCH=arm GOARM=5 $(GOBUILD) -o $(BINARY_NAME) -v