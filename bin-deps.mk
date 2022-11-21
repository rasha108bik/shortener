## mockery
MOCKERY_BIN=$(LOCAL_BIN)/mockery
$(MOCKERY_BIN):
	GOBIN=$(GOBIN)/bin go install github.com/vektra/mockery/v2@latest

## golangci-lint
GOLANGCI_BIN=$(LOCAL_BIN)/golangci-lint:
$(GOLANGCI_BIN):
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1