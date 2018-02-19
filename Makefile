PKGS = $(shell go list ./... | grep -v /vendor/)


test-unit:
	go test $(PKGS) -v

fmt:
	go fmt ./...


.PHONY: test-unit