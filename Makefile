all:
	go run ./cmd/main.go

fmt:
	go fmt ./...

up:
	docker compose up -d

down:
	docker compose down
