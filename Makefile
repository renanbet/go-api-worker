run:
	docker compose up -d --build

down:
	docker compose down

logs:
	docker compose logs -f api_service
	docker compose logs -f worker_service

test:
	go test ./...

clean:
	docker compose down -v
	docker volume rm mongo_data

lint:
	go vet ./...
	golangci-lint run ./...