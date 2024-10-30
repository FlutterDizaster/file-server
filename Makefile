BINARY_NAME=file-server
BINARY_PATH=./build

CMD_DIR=./cmd

.PHONY: all
all: build

# tests
.PHONY: test
test:
	go test ./... -v

# code generation
.PHONY: generate
generate:
	go generate ./...

# Building
.PHONY: build
build: generate
	go build -ldflags="-s -w" -o $(BINARY_PATH)/$(BINARY_NAME) $(CMD_DIR)

