MOCKGEN_BIN=$(LOCAL_BIN)/mockgen
$(MOCKGEN_BIN):
	GOBIN=$(LOCAL_BIN) go install github.com/golang/mock/mockgen@v1.6.0