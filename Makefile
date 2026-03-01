# Project name
DATE ?= $(shell date +%FT%T%z)
PROJECT_NAME ?= go-discordbot
VERSION ?= 1.0.0
TAG ?=

COMMIT := $(shell git rev-parse --short HEAD)
LDFLAGS = -ldflags "-s -w -X=main.VERSION=$(VERSION) -X=main.DATE=$(DATE) -X=main.COMMIT=$(COMMIT)"
MAIN_FILE := main.go
BUILD_PREFIX := build
BUILD_APP_PREFIX := $(BUILD_PREFIX)/app

IMAGE_NAME = $(PROJECT_NAME):$(VERSION)

# ifndef BUILD_PLATFORM
# BUILD_PLATFORM := linux/arm64
# endif

# Default target
.PHONY: all
all: clean build-arm-app build-arm64-app build-linux-app build-mac-app build-win-app

# Build Go application
.PHONY: build
build:
ifdef BUILD_PLATFORM
	@echo "BUILD_PREFIX is set to $(BUILD_PREFIX)"
	@GOOS=$(word 1, $(subst /, ,$(BUILD_PLATFORM))) && \
	GOARCH=$(word 2, $(subst /, ,$(BUILD_PLATFORM))) && \
	echo "Building for $$GOOS/$$GOARCH" && \
	CGO_ENABLED=0 GOOS=$$GOOS GOARCH=$$GOARCH go build $(LDFLAGS) -o $(PROJECT_NAME).app $(MAIN_FILE)
else
	go build $(LDFLAGS) -o $(PROJECT_NAME).app $(MAIN_FILE)
endif

# Build for ARM
.PHONY: build-arm-app
build-arm-app:
	mkdir -p $(BUILD_APP_PREFIX)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build $(LDFLAGS) -o $(BUILD_APP_PREFIX)/$(PROJECT_NAME)-arm.app $(MAIN_FILE)

# Build for ARM64
.PHONY: build-arm64-app
build-arm64-app:
	mkdir -p $(BUILD_APP_PREFIX)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_APP_PREFIX)/$(PROJECT_NAME)-arm64.app $(MAIN_FILE)

# Build for Linux
.PHONY: build-linux-app
build-linux-app:
	mkdir -p $(BUILD_APP_PREFIX)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_APP_PREFIX)/$(PROJECT_NAME)-linux.app $(MAIN_FILE)

# Build for macOS
.PHONY: build-mac-app
build-mac-app:
	mkdir -p $(BUILD_APP_PREFIX)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_APP_PREFIX)/$(PROJECT_NAME)-mac.app $(MAIN_FILE)

# Build for Windows
.PHONY: build-win-app
build-win-app:
	mkdir -p $(BUILD_APP_PREFIX)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_APP_PREFIX)/$(PROJECT_NAME)-win.app $(MAIN_FILE)

# Run Go application
.PHONY: run
run:
	go run $(LDFLAGS) $(MAIN_FILE)

# Clean build files
.PHONY: clean
clean:
	@echo "Cleaning build files..."
	go clean
	rm -rf $(BUILD_PREFIX)/*

# Build Docker image
.PHONY: docker
docker:
	docker rmi -f $(IMAGE_NAME)
ifdef BUILD_PLATFORM
	docker build --no-cache \
	--build-arg BUILD_PLATFORM=$(BUILD_PLATFORM) \
	--build-arg BUILD_DATE=$(DATE) \
	--platform $(BUILD_PLATFORM) \
	-t $(IMAGE_NAME) .
else
	docker build --no-cache --build-arg BUILD_DATE=$(DATE) -t $(IMAGE_NAME) .
endif

# Build Docker image and save as tar file
.PHONY: docker-tar
docker-tar:
ifdef BUILD_PLATFORM
	make docker BUILD_PLATFORM=$(BUILD_PLATFORM)
else
	make docker
endif
	mkdir -p $(BUILD_PREFIX)
	docker save -o $(BUILD_PREFIX)/$(PROJECT_NAME)-docker.tar $(IMAGE_NAME)

# Clean Docker images and containers
.PHONY: docker-clean
docker-clean:
	for id in $(docker images --filter=reference="$(PROJECT_NAME)*" -q); do \
	  docker ps -a --filter "ancestor=$$id" -q | xargs -r docker rm -f; \
	done
	docker rmi -f $(docker images --filter=reference="$(PROJECT_NAME)*" -q)

# Build for ARM
.PHONY: arm
arm:
	mkdir -p $(BUILD_APP_PREFIX)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build $(LDFLAGS) -o $(BUILD_APP_PREFIX)/$(PROJECT_NAME)-arm.app $(MAIN_FILE)
