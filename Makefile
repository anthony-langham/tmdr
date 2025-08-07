# tmdr - too medical; didn't read
# Build configuration

# Variables
BINARY_NAME=tmdr
VERSION=v0.3
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GOPATH=$(shell go env GOPATH)
GOBIN=$(GOPATH)/bin

# Platform specific variables
PLATFORMS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64
DIST_DIR=dist

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

.PHONY: all build clean test run install uninstall dist help

## help: Show this help message
help:
	@echo 'Usage:'
	@echo '  ${YELLOW}make${NC} ${GREEN}<target>${NC}'
	@echo ''
	@echo 'Targets:'
	@grep -E '^## ' Makefile | sed 's/## /  /'

## all: Build for current platform
all: build

## build: Build binary for current platform
build:
	@echo "Building ${BINARY_NAME} ${VERSION} for current platform..."
	@go build -o ${BINARY_NAME} .
	@echo "${GREEN}✓${NC} Built ${BINARY_NAME}"

## run: Run the application
run:
	@go run .

## test: Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	@rm -f ${BINARY_NAME}
	@rm -rf ${DIST_DIR}
	@echo "${GREEN}✓${NC} Cleaned"

## install: Install binary to GOBIN
install: build
	@echo "Installing to ${GOBIN}..."
	@mkdir -p ${GOBIN}
	@cp ${BINARY_NAME} ${GOBIN}/
	@echo "${GREEN}✓${NC} Installed to ${GOBIN}/${BINARY_NAME}"

## uninstall: Remove binary from GOBIN
uninstall:
	@echo "Uninstalling from ${GOBIN}..."
	@rm -f ${GOBIN}/${BINARY_NAME}
	@echo "${GREEN}✓${NC} Uninstalled"

## dist: Build for all platforms
dist: clean
	@echo "Building for all platforms..."
	@mkdir -p ${DIST_DIR}
	@for platform in $(PLATFORMS); do \
		GOOS=$$(echo $$platform | cut -d'/' -f1); \
		GOARCH=$$(echo $$platform | cut -d'/' -f2); \
		output_name="${DIST_DIR}/${BINARY_NAME}-${VERSION}-$$GOOS-$$GOARCH"; \
		if [ "$$GOOS" = "windows" ]; then \
			output_name="$$output_name.exe"; \
		fi; \
		echo "  Building $$GOOS/$$GOARCH..."; \
		GOOS=$$GOOS GOARCH=$$GOARCH go build -o $$output_name . || exit 1; \
		echo "  ${GREEN}✓${NC} $$output_name"; \
	done
	@echo "${GREEN}✓${NC} All platforms built"

## release: Create release archives
release: dist
	@echo "Creating release archives..."
	@cd ${DIST_DIR} && \
	for file in ${BINARY_NAME}-*; do \
		if [[ "$$file" == *.exe ]]; then \
			zip "$${file%.exe}.zip" "$$file" && echo "  ${GREEN}✓${NC} $${file%.exe}.zip"; \
		else \
			tar czf "$$file.tar.gz" "$$file" && echo "  ${GREEN}✓${NC} $$file.tar.gz"; \
		fi; \
	done
	@echo "${GREEN}✓${NC} Release archives created"

# Default target
.DEFAULT_GOAL := help