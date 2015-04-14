DEPS = $(shell go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
PACKAGES = $(shell go list ./...)

all: deps
	@echo "--> Install and Build"
	@mkdir -p bin/
	@go install
	@go build -o bin/kapal

deps:
	@echo "--> Installing build dependencies"
	@go get -d -v ./...
	@echo $(DEPS) | xargs -n1 go get -d

clean:
	@echo "--> Cleaning"
	@rm -rf bin/
	@go clean

.PHONY: all deps
