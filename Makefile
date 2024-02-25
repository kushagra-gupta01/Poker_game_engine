build:
	@go build -o bin/poker

run: build
	@./bin/poker

test:
	@go test -v ./...