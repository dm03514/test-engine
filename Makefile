PKGS = $(shell go list ./... | grep -v /vendor/)
EXECUTOR_BIN = "test-executor"
ENGINE_SERVER_BIN = "engine-server"


build-tools:
	go get -u golang.org/x/lint/golint
	go get github.com/mattn/goveralls

test-unit:
	go test $(PKGS) -v -coverprofile=coverage.out -covermode=count -tags=integration

build:
	go build -o $(EXECUTOR_BIN) ./commands/test-executor
	go build -o $(ENGINE_SERVER_BIN) ./commands/engine-server

test-functional:
	go test -tags=functional ./commands/test-executor/ -v \
		-root-test-dir=$(shell pwd)/tests/ \
		-coverprofile=coverage.functional.out -covermode=count

start-prometheus-server:
	./$(ENGINE_SERVER_BIN) -testDir=$(shell pwd)/tests -metrics=prometheus 2>&1 | jq .

test-prometheus-server:
	curl -X POST localhost:8080/execute?test=multiple_states.yml

start-test-stub-server:
	go run tests/commands/stub-server/*.go

# test-stub-server:
# $(EXECUTOR_BIN) -test $(shell pwd)/tests/gstreamer.yml
# $(EXECUTOR_BIN) -test $(shell pwd)/tests/dmonitor.yml

fmt:
	go fmt ./...

lint:
	golint $(PKGS)

.PHONY: test-unit test-functional fmt lint start-prometheus-server test-prometheus-server build-tools
