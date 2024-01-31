.DEFAULT_GOAL := help

BIN = $(CURDIR)/bin
VERSION ?= $(shell git describe --tags --always --dirty --match='v*' 2> /dev/null || echo v0)
SOURCES = $(shell find . -type f -name *.go)
TARGET = $(BIN)/instacart-export

$(BIN):
	@mkdir -p $@

.PHONY:phony

push: tidy no-dirty ## push changes to github with checks
	git push

no-dirty:
	git diff --exit-code

tidy: phony ## verify sources
	go fmt ./...
	go mod tidy -v
	go mod verify
	go vet ./...

$(TARGET): $(SOURCES) | $(BIN)
	go build \
		-tags release \
		-ldflags '-X main.Version=$(VERSION)' \
		-o $(BIN) ./...

build: $(TARGET) ## build the binary

run: $(TARGET) ## run the binary
	./$(TARGET)

clean: phony
	rm -rf $(BIN) data/*

GREEN  := $(shell tput -Txterm setaf 2)
RESET  := $(shell tput -Txterm sgr0)

version: phony ## print the version
	@echo $(VERSION)

help: phony ## print this help message
	@awk -F ':|##' '/^[^\t].+?:.*?##/ { printf "${GREEN}%-20s${RESET}%s\n", $$1, $$NF }' $(MAKEFILE_LIST)
