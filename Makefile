
.PHONY: all clean generate-parser generate build


# Project directories
BUILD_DIR := build
BIN_DIR := $(BUILD_DIR)/bin


all: clean generate build


clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@rm -rf internal/compiler/parser/*.go
	@rm -rf internal/compiler/parser/generated/
	@find internal/compiler/parser -name "*.tokens" -delete
	@find internal/compiler/parser -name "*.interp" -delete
	@echo "Clean complete"


generate-parser:
	@echo "Generating Compiler parser..."
	@cd internal/compiler/parser && ./generate.sh
	@echo "Compiler parser generation complete"


generate: generate-parser


build:
	@echo "Building project..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/jawt main.go
	@echo "Build complete. Binary available at: $(BIN_DIR)/jawt"


run: build
	@$(BIN_DIR)/jawt




