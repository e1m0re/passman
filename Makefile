lint:
	golangci-lint run ./...

format-code:
	goimports --local "github.com/e1m0re/passman" -w .

test:
	cd server;\
	go test -race -covermode=atomic -coverprofile=coverage.out ./...

build-server:
	go build -o bin/server cmd/server/*go

build-client:
	GOOS=windows GOARCH=386 go build -o bin/client_win_x86.exe cmd/client/main.go ;\
	GOOS=windows GOARCH=amd64 go build -o bin/client_win_x64.exe cmd/client/main.go ;\
	GOOS=linux GOARCH=amd64 go build -o bin/client_linux_x64 cmd/client/main.go ;\
	GOOS=darwin GOARCH=amd64 go build -o bin/client_mac_x64 cmd/client/main.go ;\
	go build -o bin/client_mac_arm cmd/client/main.go ;\

statickcheck:
	go run cmd/staticlint/main.go ./...

generate:
	go generate ./...

migrates:
	goose -dir internal/db/migrations postgres "postgresql://passman:passman@127.0.0.1:5432/passman?sslmode=disable" up

grpc-gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/*.proto
