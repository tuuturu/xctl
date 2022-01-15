GOLANGCI_LINT:
	@which golangci-lint &>/dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.43.0

fmt:
	@goimports -w .
	@gofmt -w .

lint: GOLANGCI_LINT
	@golangci-lint run

test:
	@go test ./...

check: fmt lint test
