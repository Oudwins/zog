test-watch:
	@gotestsum --format=short-verbose --watch

test:
	@gotestsum --format=short-verbose

lint:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run