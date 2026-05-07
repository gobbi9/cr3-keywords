APP := cr3
BIN_DIR := bin
CMD := ./cmd/cr3-keyword
OUT := $(BIN_DIR)/$(APP)

INSTALL_PATH ?= /usr/local/bin

VERSION_BASE := 0.0.1
GIT_HASH := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
VERSION := $(VERSION_BASE)-$(GIT_HASH)

.PHONY: help tidy build rebuild run clean install uninstall

help:
	@echo "Targets:"
	@echo "  make tidy      - run go mod tidy"
	@echo "  make build     - build $(OUT) with CGO_ENABLED=0"
	@echo "  make rebuild   - clean + build"
	@echo "  make run       - run with ARGS='...'"
	@echo "  make install   - install $(APP) to $(INSTALL_PATH)"
	@echo "  make uninstall - remove $(INSTALL_PATH)/$(APP)"
	@echo "  make clean     - remove built binaries"

tidy:
	go mod tidy

build:
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)" -o $(OUT) $(CMD)

rebuild: clean build

run: build
	./$(OUT) $(ARGS)

install: build
	install -m 0755 $(OUT) $(INSTALL_PATH)/$(APP)

uninstall:
	rm -f $(INSTALL_PATH)/$(APP)

clean:
	rm -rf $(BIN_DIR)
