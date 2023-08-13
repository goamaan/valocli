BINARY_NAME=valocli
BUILD_DIR=bin

all: build run

build:
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) *.go

run:
	@$(BUILD_DIR)/$(BINARY_NAME)

clean:
	@rm -rf $(BUILD_DIR)

.PHONY: all build run clean