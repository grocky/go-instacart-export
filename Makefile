.DEFAULT_GOAL := help

BIN = $(CURDIR)/bin
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || echo v0)

$(BIN):
	@mkdir -p $@

.PHONY:phony

fmt: phony ## format the codes
	@go fmt ./...

lint: phony fmt ## lint the codes
	@golint ./...

vet: phony lint ## vet the codes
	@go vet ./...

build: phony vet | $(BIN) ## build the binary
	@go build \
		-tags release \
		-ldflags '-X main.Version=$(VERSION)' \
		-o $(BIN)/ ./...

run: phony vet ## run the binary
	@go run main.go

clean: phony
	rm -rf $(BIN)

GREEN  := $(shell tput -Txterm setaf 2)
RESET  := $(shell tput -Txterm sgr0)

help: phony ## print this help message
	@awk -F ':|##' '/^[^\t].+?:.*?##/ { printf "${GREEN}%-20s${RESET}%s\n", $$1, $$NF }' $(MAKEFILE_LIST)
