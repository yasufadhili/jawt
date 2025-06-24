
.PHONY: all clean generate-pc generate-cc generate build


# Project directories
BUILD_DIR := build
BIN_DIR := $(BUILD_DIR)/bin


all: clean generate build


clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@rm -rf internal/page_compiler/parser/*.go internal/component_compiler/parser/*.go
	@find internal/page_compiler/parser internal/component_compiler/parser -name "*.tokens" -delete
	@find internal/page_compiler/parser internal/page_compiler/parser -name "*.interp" -delete
	@echo "Clean complete"


generate-pc:
	@echo "Generating Page Compiler parser..."
	@cd internal/page_compiler/parser && ./generate.sh
	@echo "Page Compiler parser generation complete"


generate-cc:
	@echo "Generating Component Compiler parser..."
	@cd internal/component_compiler/parser && ./generate.sh
	@echo "Component Compiler parser generation complete"


generate: generate-pc generate-cc


build:
	@echo "Building project..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/jawt main.go
	@echo "Build complete. Binary available at: $(BIN_DIR)/jawt"


run: build
	@$(BIN_DIR)/jawt




