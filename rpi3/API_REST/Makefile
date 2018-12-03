# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_x64=reservas_api_amd64
BINARY_RPI=reservas_api_arm

all: build
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
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_x64) -v
amd64: # Same as Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_x64) -v
rpi:
	GOOS=linux GOARCH=arm GOARM=5 $(GOBUILD) -o $(BINARY_RPI) -v
arm: # Same as rpi
	GOOS=linux GOARCH=arm GOARM=5 $(GOBUILD) -o $(BINARY_RPI) -v