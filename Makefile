# SPX Graph Visualizer Makefile

.PHONY: build run test clean install

# Vars
BINARY_NAME=spx-graph
BINARY_DIR=bin
CMD_DIR=./cmd

# Build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BINARY_DIR)
	@go build -o $(BINARY_DIR)/$(BINARY_NAME) $(CMD_DIR)

# Release

release:
	@echo "Build release..."
	@GOOS=linux GOARCH=amd64 go build -o $(BINARY_DIR)/$(BINARY_LINUX) $(CMD_DIR)
	@GOOS=darwin GOARCH=arm64 go build -o $(BINARY_DIR)/$(BINARY_MACOS) $(CMD_DIR)

# Tests
test:
	go test ./...

# Очистка
clean:
	rm -f $(BINARY_NAME)
	go clean
