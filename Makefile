test-watch:
	@gotestsum --format=short-verbose --watch

test:
	@gotestsum --format=short-verbose

test-cover:
	@CGO_ENABLED=1 gotestsum -- -race -covermode=atomic -coverprofile="profile.cov" ./... && go tool cover -func=profile.cov | grep -v "100.0%"

lint:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run