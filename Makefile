APP := cr3
BIN_DIR := bin
CMD := ./cmd/cr3-keyword
OUT := $(BIN_DIR)/$(APP)

INSTALL_PATH ?= /usr/local/bin

GIT_HASH := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
VERSION ?= $(shell git describe --tags --always --dirty)
CGO_ENABLED ?= 0

.PHONY: help tidy build rebuild run clean install uninstall release-tag

help:
	@echo "Targets:"
	@echo "  make tidy                         - run go mod tidy"
	@echo "  make build [VERSION=x.y.z]        - build $(OUT) (override OUT/GOOS/GOARCH/CGO_ENABLED as needed)"
	@echo "  make rebuild [VERSION=x.y.z]      - clean + build"
	@echo "  make run ARGS='...'               - run with ARGS"
	@echo "  make install [VERSION=x.y.z]      - install $(APP) to $(INSTALL_PATH)"
	@echo "  make uninstall                    - remove $(INSTALL_PATH)/$(APP)"
	@echo "  make release-tag VERSION=x.y.z    - create and push tag vX.Y.Z from current commit"
	@echo "  make clean                        - remove built binaries"

tidy:
	go mod tidy

build:
	mkdir -p $(dir $(OUT))
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)" -o $(OUT) $(CMD)

rebuild: clean build

run: build
	./$(OUT) $(ARGS)

install: build
	install -m 0755 $(OUT) $(INSTALL_PATH)/$(APP)

uninstall:
	rm -f $(INSTALL_PATH)/$(APP)

release-tag:
	@if [ -z "$(VERSION)" ] || [ "$(VERSION)" = "dev" ]; then \
		echo "Usage: make release-tag VERSION=x.y.z"; \
		exit 1; \
	fi
	@echo "$(VERSION)" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+$$' || (echo "VERSION must look like x.y.z" && exit 1)
	@git diff --quiet || (echo "Working tree is dirty. Commit or stash changes first." && exit 1)
	@git diff --cached --quiet || (echo "Staged but uncommitted changes found. Commit first." && exit 1)
	@git rev-parse --verify "v$(VERSION)" >/dev/null 2>&1 && (echo "Tag v$(VERSION) already exists." && exit 1) || true
	@git tag -a "v$(VERSION)" -m "Release v$(VERSION)"
	@git push origin "v$(VERSION)"
	@echo "Pushed tag v$(VERSION). GitHub Actions will publish release binaries."

clean:
	rm -rf $(BIN_DIR)
