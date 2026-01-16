SHELL := bash


.PHONY: install
install:
	go generate
	go install
