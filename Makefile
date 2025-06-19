
.PHONY: all clean generate-pc generate-cc generate build


# Project directories
BUILD_DIR := build
BIN_DIR := $(BUILD_DIR)/bin


all: clean generate build


clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@rm -rf internal/pc/parser/*.go internal/cc/parser/*.go
	@find internal/pc/parser internal/cc/parser -name "*.tokens" -delete
	@find internal/pc/parser internal/cc/parser -name "*.interp" -delete
	@echo "Clean complete"


generate-pc:
	@echo "Generating Page Compiler parser..."
	@cd internal/pc/parser && ./generate.sh
	@echo "Page Compiler parser generation complete"


generate-cc:
	@echo "Generating Component Compiler parser..."
	@cd internal/cc/parser && ./generate.sh
	@echo "Component Compiler parser generation complete"


generate: generate-pc generate-cc


build:
	@echo "Building project..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/jawt main.go
	@echo "Build complete. Binary available at: $(BIN_DIR)/jawt"


run: build
	@$(BIN_DIR)/jawt




