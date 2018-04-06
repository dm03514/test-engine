PKGS = $(shell go list ./... | grep -v /vendor/)
EXECUTOR_BIN = "test-executor"
ENGINE_SERVER_BIN = "engine-server"


build-tools:
	go get -u golang.org/x/lint/golint
	go get github.com/mattn/goveralls

test-unit:
	go test $(PKGS) -v -coverprofile=coverage.out -covermode=count

build:
	go build -o $(EXECUTOR_BIN) ./commands/test-executor
	go build -o $(ENGINE_SERVER_BIN) ./commands/engine-server

test-functional:
	go test -tags=functional ./commands/test-executor/ -v \
		-root-test-dir=$(shell pwd)/tests/ \
		-coverprofile=coverage.functional.out -covermode=count
	# ./$(EXECUTOR_BIN) -test $(shell pwd)/tests/subprocess_exit_code.yml
	# ./$(EXECUTOR_BIN) -test $(shell pwd)/tests/subprocess_multiple_conditions.yml
	# ./$(EXECUTOR_BIN) -test $(shell pwd)/tests/multiple_states.yml
	#./$(EXECUTOR_BIN) -test $(shell pwd)/tests/previous_state_overrides.yml

start-prometheus-server:
	./$(ENGINE_SERVER_BIN) -testDir=$(shell pwd)/tests -metrics=prometheus 2>&1 | jq .

test-prometheus-server:
	curl -X POST localhost:8080/execute?test=multiple_states.yml

fmt:
	go fmt ./...

lint:
	golint $(PKGS)

.PHONY: test-unit test-functional fmt lint start-prometheus-server test-prometheus-server build-tools
