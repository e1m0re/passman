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

migrates:
	goose -dir server/db/migrations postgres "postgresql://passman:passman@127.0.0.1:5432/passman?sslmode=disable" up

grpc-gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/test.proto
