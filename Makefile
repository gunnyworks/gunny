BUILD_DIR ?= ./build
OUTPUT ?= $(BUILD_DIR)/gunny

.PHONY: build
build:
	go build -o $(OUTPUT) ./cmd/gunny/

.PHONY: test
test:
	go test -v ./...

.PHONY: tidy
tidy:
	go mod tidy -v
