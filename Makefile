.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -v -cover ./... -coverprofile=coverage.out

.PHONY: coverage
coverage:
	go tool cover -html=coverage.out
