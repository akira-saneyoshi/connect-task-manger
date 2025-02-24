# Go-gRPC
tidy:
	go mod tidy

run:
	go run cmd/server/main.go

build:
	go build -o connect-task-manger cmd/server/main.go

migrate:
	go run cmd/migration/main.go

sqlc:
	sqlc generate

buf:
	buf generate