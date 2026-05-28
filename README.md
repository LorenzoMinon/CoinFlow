## Stack

- Go 1.25+
- PostgreSQL 15
- Docker Compose
- pgx v5

## Setup

```bash
docker compose up db -d
```

Create the table:

```bash
docker exec -it coinflow-db-1 psql -U admin -d coinflow
```

```sql
CREATE TABLE coin_prices (
    id          SERIAL PRIMARY KEY,
    coin_id     VARCHAR(50) NOT NULL,
    symbol      VARCHAR(10) NOT NULL,
    price_usd   DECIMAL(20,8) NOT NULL,
    change_24h  DECIMAL(10,4),
    sentiment   VARCHAR(10),
    fetched_at  TIMESTAMP DEFAULT NOW()
);
```

Set the environment variable:

```bash
export DB_URL="postgres://admin:admin@localhost:5434/coinflow?sslmode=disable"
```

Run:

```bash
go run main.go
```

## Endpoints

### POST /pipeline/run
Triggers the ETL pipeline. Fetches bitcoin and ethereum prices in parallel, calculates sentiment, and stores results in PostgreSQL.

Returns:
```json
{ "processed": 2, "failed": 0 }
```

### GET /coins
Returns all stored coin price records ordered by fetch time.

### GET /coins/latest
Returns the latest price for each coin using DISTINCT ON.

## Sentiment Logic

| Change 24h | Sentiment |
|---|---|
| > 2% | bullish |
| < -2% | bearish |
| between | neutral |