GOPATH := $(shell go env GOPATH)
GOBIN  ?= $(GOPATH)/bin

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
	@go test ./... | tee >($(RICHGO) testfilter)

check: fmt lint test
