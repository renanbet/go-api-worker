# Sistema de pedidos assíncronos (Go + Kafka + MongoDB)

Este projeto simula a criação e o processamento assíncrono de pedidos usando:

- **API (Produtor)**: recebe `POST /orders`, cria o pedido no MongoDB com status **CREATED** e publica evento no Kafka.
- **Worker (Consumidor)**: consome o evento, atualiza status para **PROCESSING**, espera 2 segundos e finaliza como **CONCLUDED**.

## Como rodar

Subir tudo (build + start):

```bash
docker compose up -d --build
```

Ver logs (opcional):

```bash
docker compose logs -f api_service
docker compose logs -f worker_service
```

## Criar um pedido

```bash
curl -i -X POST "http://localhost:8080/orders" \
  -H "Content-Type: application/json" \
  -d '{"product":"teclado","quantity":2}'
```

Resposta esperada: **HTTP 201** com um JSON contendo `order_id`.

## Verificar no MongoDB

Entrar no container do Mongo:

```bash
docker exec -it mongodb mongosh
```

Dentro do `mongosh`:

```javascript
use orders_db
db.orders.find().sort({created_at:-1}).limit(5).pretty()
```

Você deve ver o pedido primeiro como `CREATED` e, após ~2 segundos, como `CONCLUDED` (o worker passa por `PROCESSING` antes).

## Rodar testes unitários (local)

```bash
go test ./...
```

