all: build

build:
	go vet ./...
	go test -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	rm -f coverage.out
	go install golang.org/x/lint/golint@latest
	golint -set_exit_status=true ./...
	go build ./...
	go mod tidy

build_fast:
	go test -failfast -short -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	rm -f coverage.out
	go build ./...

clean:
	go clean