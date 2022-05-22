start:
	go run main.go

lint:
	golangci-lint run --timeout=5m

.PHONY: test

test:
	go clean --testcache
	go test ./...
