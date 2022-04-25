GOPATH := $(shell go env GOPATH)
GOBIN  ?= $(GOPATH)/bin
BUILD_DIR := ./build
INSTALL_DIR := ~/.local/bin

GOLANGCILINT := $(GOBIN)/golangci-lint
$(GOLANGCILINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.43.0

RICHGO := $(GOBIN)/richgo
$(RICHGO):
	@go install github.com/kyoh86/richgo@v0.3.6

fmt:
	@goimports -w .
	@gofmt -w .

lint: $(GOLANGCILINT)
	@golangci-lint run

test: $(RICHGO)
	@$(RICHGO) test ./...

check: fmt lint test

build:
	mkdir -p ${BUILD_DIR}
	go build -o ${BUILD_DIR}/xctl

clean:
	rm -r ${BUILD_DIR}

install:
	mkdir -p ${INSTALL_DIR}
	cp ${BUILD_DIR}/xctl ${INSTALL_DIR}/xctl

uninstall:
	rm ${INSTALL_DIR}/xctl
