LOCAL_BIN=$(CURDIR)/bin

.PHONY: lint
lint: $(LOCAL_BIN) ## lint
	$(LOCAL_BIN) run ./...

.PHONY: app-run
app-run: ## run app
	go run cmd/shortener/main.go

.PHONY: unit-test 
unit-test: ## unit-test 
	go test -count=1 -v ./...

.PHONY: help
help: ## help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
