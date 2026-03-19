SHELL := bash


install:
	go generate
	go install -ldflags="-X main.version=dev-$(shell git log -1 --pretty=%h)" .

fmt:
	@golangci-lint fmt .
	@goimports -local $(shell go list -m) -w .
	@gofumpt -l -w .

static-check:
	golangci-lint run


.PHONY: install fmt static-check
