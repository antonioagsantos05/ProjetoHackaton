.PHONY: run test build docker-up docker-down fmt vet coverage

run:
	go run cmd/api/main.go

worker:
	go run cmd/worker/main.go

test:
	go test ./... -v

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	go tool cover -func=coverage.out

build:
	go build -o bin/api cmd/api/main.go
	go build -o bin/worker cmd/worker/main.go

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

fmt:
	go fmt ./...

vet:
	go vet ./...
