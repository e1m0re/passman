lint:
	golangci-lint run ./...

format-code:
	goimports --local "github.com/e1m0re/passman" -w .

test:
	go test -race -covermode=atomic -coverprofile=coverage.out ./...

build-server:
	go build -o bin/server cmd/server/*go

statickcheck:
	go run cmd/staticlint/main.go ./...

generate:
	go generate ./...