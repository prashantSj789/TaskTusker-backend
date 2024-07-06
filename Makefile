build:
	@go build -o bin/go-jira

run: build
	@./bin/go-jira

test:
	@go test -v ./...
