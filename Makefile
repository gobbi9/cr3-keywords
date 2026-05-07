APP := cr3-keyword
BIN_DIR := bin
CMD := ./cmd/cr3-keyword
OUT := $(BIN_DIR)/$(APP)

.PHONY: help tidy build rebuild run clean

help:
	@echo "Targets:"
	@echo "  make tidy    - run go mod tidy"
	@echo "  make build   - build $(OUT) with CGO_ENABLED=0"
	@echo "  make rebuild - clean + build"
	@echo "  make run     - run with ARGS='...'"
	@echo "  make clean   - remove built binaries"

tidy:
	go mod tidy

build:
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 go build -o $(OUT) $(CMD)

rebuild: clean build

run: build
	./$(OUT) $(ARGS)

clean:
	rm -rf $(BIN_DIR)
