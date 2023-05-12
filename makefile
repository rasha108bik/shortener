LOCAL_BIN=$(CURDIR)/bin

include bin-deps.mk

default: help

.PHONY: help
help: ## help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: app-run
app-run: ## run app
	go run $(PWD)/cmd/shortener/main.go

.PHONY: unit-test 
unit-test: ## unit-test 
	go test -count=1 -v ./...

# .PHONY: mockgen-install
# mockgen-install: ## mockgen-install
# 	GOBIN=$(LOCAL_BIN) go install github.com/golang/mock/mockgen@v1.6.0

.PHONY: pprof-diff
	go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof

.PHONY: go-generate-all
go-generate-all: ## go-generate-all
	PATH=$(LOCAL_BIN):$(PATH) go generate ./...
