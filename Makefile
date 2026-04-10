up:
	docker compose up -d --build --scale worker_service=3

run:
	docker compose up -d --build --scale worker_service=3

down:
	docker compose down

logs-api:
	docker compose logs -f api_service

logs-worker:
	docker compose logs -f worker_service

test:
	go test ./...

clean:
	docker compose down -v
	docker volume rm mongo_data

lint:
	go vet ./...
	golangci-lint run ./...

build-api:
	docker compose build api_service --no-cache

build-worker:
	docker compose build worker_service --no-cache

restart-worker:
	docker compose restart worker_service