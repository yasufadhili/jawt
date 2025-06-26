

.PHONY: all clean generate-parser generate build \
        build-linux-amd64 build-linux-arm64 \
        build-windows-amd64 build-windows-arm64 \
        build-macos-amd64 build-macos-arm64 \
        test run help



BUILD_DIR := build
BIN_DIR := $(BUILD_DIR)/bin



MAIN_SRC := main.go



LINUX_AMD64_BIN := $(BIN_DIR)/jawt-linux-amd64
LINUX_ARM64_BIN := $(BIN_DIR)/jawt-linux-arm64
WINDOWS_AMD64_BIN := $(BIN_DIR)/jawt-windows-amd64.exe
WINDOWS_ARM64_BIN := $(BIN_DIR)/jawt-windows-arm64.exe
MACOS_AMD64_BIN := $(BIN_DIR)/jawt-macos-amd64
MACOS_ARM64_BIN := $(BIN_DIR)/jawt-macos-arm64


.DEFAULT_GOAL := help


all: clean generate build


help:
	@echo "jawt Makefile - targets:"
	@echo "  make build                - Build the project for the host platform"
	@echo "  make build-all            - Build for all major platforms/architectures"
	@echo "  make test                 - Run the project tests"
	@echo "  make clean                - Remove build and generated files"
	@echo "  make generate             - Generate code (compiler/parser)"
	@echo "  make run                  - Run the built binary (host platform)"
	@echo "  make help                 - Show this help message"


clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@rm -rf internal/compiler/parser/*.go
	@rm -rf internal/compiler/parser/generated/
	@find internal/compiler/parser -name "*.tokens" -delete
	@find internal/compiler/parser -name "*.interp" -delete
	@echo "Clean complete"


generate-parser:
	@echo "Generating compiler parser..."
	@cd internal/compiler/parser && ./generate.sh
	@echo "Compiler parser generation complete"


generate: generate-parser


build:
	@echo "Building project for host platform..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/jawt $(MAIN_SRC)
	@echo "Build complete. Binary available at: $(BIN_DIR)/jawt"


build-linux-amd64:
	@echo "Building Linux amd64..."
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(LINUX_AMD64_BIN) $(MAIN_SRC)


build-linux-arm64:
	@echo "Building Linux arm64..."
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=arm64 go build -o $(LINUX_ARM64_BIN) $(MAIN_SRC)


build-windows-amd64:
	@echo "Building Windows amd64..."
	@mkdir -p $(BIN_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(WINDOWS_AMD64_BIN) $(MAIN_SRC)


build-windows-arm64:
	@echo "Building Windows arm64..."
	@mkdir -p $(BIN_DIR)
	GOOS=windows GOARCH=arm64 go build -o $(WINDOWS_ARM64_BIN) $(MAIN_SRC)


build-macos-amd64:
	@echo "Building macOS amd64..."
	@mkdir -p $(BIN_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(MACOS_AMD64_BIN) $(MAIN_SRC)


build-macos-arm64:
	@echo "Building macOS arm64..."
	@mkdir -p $(BIN_DIR)
	GOOS=darwin GOARCH=arm64 go build -o $(MACOS_ARM64_BIN) $(MAIN_SRC)


build-all: build-linux-amd64 build-linux-arm64 build-windows-amd64 build-windows-arm64 build-macos-amd64 build-macos-arm64
	@echo "All cross-platform builds complete. Binaries are in $(BIN_DIR)."


test:
	@echo "Running tests (with race detector and verbose output)..."
	@go test -race -v ./...
	@echo "Testing complete."


run: build
	@echo "Running jawt..."
	@$(BIN_DIR)/jawt
