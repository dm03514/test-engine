PKGS = $(shell go list ./... | grep -v /vendor/)
EXECUTOR_BIN = "test-executor"

test-unit:
	go test $(PKGS) -v

build:
	go build -o $(EXECUTOR_BIN) ./commands/test-executor
	go build ./commands/engine-server

test-functional:
	./$(EXECUTOR_BIN) -test $(shell pwd)/tests/subprocess_exit_code.yml
	./$(EXECUTOR_BIN) -test $(shell pwd)/tests/subprocess_multiple_conditions.yml
	./$(EXECUTOR_BIN) -test $(shell pwd)/tests/multiple_states.yml
	./$(EXECUTOR_BIN) -test $(shell pwd)/tests/previous_state_overrides.yml

fmt:
	go fmt ./...

lint:
	golint $(PKGS)

.PHONY: test-unit test-functional fmt lint