install:
	go mod tidy
	go mod download
	go mod vendor

start:
	go run main.go start

migrate:
	go run main.go migrate

rollback:
	go run main.go rollback