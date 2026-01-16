SHELL := bash


.PHONY: install
install:
	go generate
	go install

.PHONY: fmt
fmt:
	@goimports -local $(shell go list -m) -w .
	@gofumpt -l -w .

