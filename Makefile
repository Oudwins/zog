test-watch:
	watch -n 1 go test -v ./...

test:
	@go test -v ./...

lint:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run