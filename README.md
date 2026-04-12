# Async Order System (Go + Kafka + RabbitMQ + MongoDB)

This project simulates the creation and asynchronous processing of orders using an event-driven architecture with two message brokers.

## Services

- **API**: receives `POST /orders`, creates the order in MongoDB with status **CREATED** and publishes an event to Kafka.
- **Worker**: consumes the Kafka event, updates the order status to **PROCESSING**, waits 2 seconds and marks it as **CONCLUDED**.
- **Email Worker**: consumes the same Kafka topic (`order_events`) with a separate consumer group, enqueues an email job to RabbitMQ, then processes the queue by simulating an email send (2-second fake processing).

## Architecture

```
POST /orders
    │
    ▼
[API Service] ──► Kafka: order_events
                        │
          ┌─────────────┴──────────────┐
          ▼                            ▼
  [Worker Service]            [Email Worker Service]
  group: orders-worker        group: email-worker
  CREATED → PROCESSING             │
         → CONCLUDED               ▼
                          RabbitMQ: email_queue
                                   │
                                   ▼
                          fake email send (2s)
```

## Running

Build and start all services:

```bash
docker compose up -d --build
```

View logs:

```bash
docker compose logs -f api_service
docker compose logs -f worker_service
docker compose logs -f email_service
```

## Creating an order

```bash
curl -i -X POST "http://localhost:8080/orders" \
  -H "Content-Type: application/json" \
  -d '{"product":"keyboard","quantity":2}'
```

Expected response: **HTTP 201** with a JSON body containing `order_id`.

## Checking the order in MongoDB

Enter the MongoDB container:

```bash
docker exec -it mongodb mongosh
```

Inside `mongosh`:

```javascript
use orders_db
db.orders.find().sort({created_at:-1}).limit(5).pretty()
```

You should see the order first as `CREATED`, then as `CONCLUDED` after ~2 seconds (the worker transitions through `PROCESSING` in between).

## Monitoring

| Tool | URL | Purpose |
|---|---|---|
| Kafdrop | http://localhost:9000 | Kafka topics and messages |
| RabbitMQ Management | http://localhost:15672 | Queues and message rates (guest/guest) |

## Running unit tests

```bash
go test ./...
```
