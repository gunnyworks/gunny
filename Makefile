BUILD_DIR ?= ./build
OUTPUT ?= $(BUILD_DIR)/gunny
VERSION ?= $(shell git describe --tags $(git rev-list --tags --max-count=1))
BUILD ?= $(shell git rev-parse --short HEAD)
LDFLAGS ?= -ldflags "-X=main.gunnyVersion=$(VERSION) -X=main.gunnyBuild=$(BUILD)"

.PHONY: build
build:
	go build -o $(OUTPUT) $(LDFLAGS) ./cmd/gunny/

.PHONY: test
test:
	go test -v ./...

.PHONY: tidy
tidy:
	go mod tidy -v
