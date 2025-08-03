all:
	go run ./cmd/road2ca-api/main.go

seed:
	go run ./cmd/seed/main.go

fmt:
	go fmt ./...

up:
	docker compose up -d

down:
	docker compose down

rmvolumes:
	@volumes=$$(docker volume ls -qf dangling=true); \
    if [ -n "$$volumes" ]; then \
        docker volume rm $$volumes; \
    else \
        echo "No dangling volumes to remove."; \
    fi
