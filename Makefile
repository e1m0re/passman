lint:
	golangci-lint run ./...

format-code:
	goimports --local "github.com/e1m0re/passman" -w .

test:
	go test -race -covermode=atomic -coverprofile=coverage.out ./...

build-server:
	go build -o bin/server cmd/server/*go

build-client:
	rm ./bin/* ;\
	GOOS=windows GOARCH=386 go build -ldflags "-X github.com/e1m0re/passman/internal/client/app.BuildVersion=0.1.0 -X github.com/e1m0re/passman/internal/client/app.BuildDate=03.09.2024"  -o bin/client_win_x86.exe cmd/client/main.go ;\
	GOOS=windows GOARCH=amd64 go build -ldflags "-X github.com/e1m0re/passman/internal/client/app.BuildVersion=0.1.0 -X github.com/e1m0re/passman/internal/client/app.BuildDate=03.09.2024"  -o bin/client_win_x64.exe cmd/client/main.go ;\
	GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/e1m0re/passman/internal/client/app.BuildVersion=0.1.0 -X github.com/e1m0re/passman/internal/client/app.BuildDate=03.09.2024"  -o bin/client_linux_x64 cmd/client/main.go ;\
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X github.com/e1m0re/passman/internal/client/app.BuildVersion=0.1.0 -X github.com/e1m0re/passman/internal/client/app.BuildDate=03.09.2024"  -o bin/client_mac_x64 cmd/client/main.go ;\
	go build -ldflags "-X github.com/e1m0re/passman/internal/client/app.BuildVersion=0.1.0 -X github.com/e1m0re/passman/internal/client/app.BuildDate=03.09.2024" -o bin/client_mac_arm cmd/client/main.go ;\
	cp config/passman.yml bin/

statickcheck:
	go run cmd/staticlint/main.go ./...

generate:
	go generate ./...

migrates:
	goose -dir internal/server/db/migrations postgres "postgresql://passman:passman@127.0.0.1:5432/passman?sslmode=disable" up

grpc-gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/*.proto
