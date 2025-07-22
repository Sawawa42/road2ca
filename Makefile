all:
	go run ./cmd/main.go

seed:
	go run ./cmd/seed/main.go

fmt:
	go fmt ./...

up:
	docker compose up -d

down:
	docker compose down
