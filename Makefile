PKGS = $(shell go list ./... | grep -v /vendor/)
EXECUTOR_BIN = "test-executor"
ENGINE_SERVER_BIN = "engine-server"

test-unit:
	go test $(PKGS) -v

build:
	go build -o $(EXECUTOR_BIN) ./commands/test-executor
	go build -o $(ENGINE_SERVER_BIN) ./commands/engine-server

test-functional:
	./$(EXECUTOR_BIN) -test $(shell pwd)/tests/subprocess_exit_code.yml
	./$(EXECUTOR_BIN) -test $(shell pwd)/tests/subprocess_multiple_conditions.yml
	./$(EXECUTOR_BIN) -test $(shell pwd)/tests/multiple_states.yml
	./$(EXECUTOR_BIN) -test $(shell pwd)/tests/previous_state_overrides.yml

start-prometheus-server:
	./$(ENGINE_SERVER_BIN) -testDir=$(shell pwd)/tests -metrics=prometheus 2>&1 | jq .

test-prometheus-server:
	curl -X POST localhost:8080/execute?test=multiple_states.yml

fmt:
	go fmt ./...

lint:
	golint $(PKGS)

.PHONY: test-unit test-functional fmt lint start-prometheus-server test-prometheus-server