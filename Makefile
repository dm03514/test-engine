PKGS = $(shell go list ./... | grep -v /vendor/)
BIN = "test-engine"


test-unit:
	go test $(PKGS) -v

build:
	go build -o $(BIN) .

test-functional:
	./$(BIN) -test $(shell pwd)/tests/subprocess_exit_code.yml
	./$(BIN) -test $(shell pwd)/tests/subprocess_multiple_conditions.yml


fmt:
	go fmt ./...


.PHONY: test-unit