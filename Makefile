BUILD_DIR ?= ./build
OUTPUT ?= $(BUILD_DIR)/gunny
VERSION ?= $(shell git describe --abbrev=0 --tags $(git rev-list --tags --max-count=1))
COMMIT ?= $(shell git rev-parse --short HEAD)
LDFLAGS ?= -ldflags "-X=main.version=$(VERSION) -X=main.commit=$(COMMIT)"

.PHONY: build
build:
	go build -o $(OUTPUT) $(LDFLAGS) ./cmd/gunny/

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: tidy
tidy:
	go mod tidy -v
