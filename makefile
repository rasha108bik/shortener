default: help

LOCAL_BIN=$(CURDIR)/bin
TIMEOUT = 30s

include bin-deps.mk

.PHONY: help
help: ## help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: tools
tools: ## instal binary
	cd tools && go mod tidy && go generate -tags tools

.PHONY: app-run
app-run: ## run app
	go run -ldflags "-X main.buildVersion=v1.0.0 -X 'main.buildDate=$(shell date)'" $(PWD)/cmd/shortener/main.go

.PHONY: unit-test 
unit-test: ## unit-test 
	go test -cover -race -timeout $(TIMEOUT) ./... | column -t | sort -r

# .PHONY: mockgen-install
# mockgen-install: ## mockgen-install
# 	GOBIN=$(LOCAL_BIN) go install github.com/golang/mock/mockgen@v1.6.0

.PHONY: pprof-diff
	go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof

.PHONY: go-generate-all
go-generate-all: ## go-generate-all
	PATH=$(LOCAL_BIN):$(PATH) go generate ./...

# if ypu want to open docs for api, you need to open: http://localhost:8080/pkg/?m=all 
.PHONY: godoc-play 
godoc-play: ## godoc play
	godoc -play -http=:8080

.PHONY: statichcheck 
statichcheck: ## statichcheck
	go run $(PWD)/cmd/staticlint/main.go

.PHONY: install-githooks
install-githooks:
	cp ./githooks/* .git/hooks
