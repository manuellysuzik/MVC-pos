export GOPROXY=https://proxy.golang.org,direct
export GONOSUMDB=*
export GONOPROXY=

.PHONY: run test build tidy

run:
	go run main.go

build:
	go build -o bin/ecommerce-api ./...

test:
	go test ./...

tidy:
	go mod tidy
