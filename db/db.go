package db

import (
	"context"
	"fmt"
	"os"

	"github.com/LorenzoMinon/coinflow/pipeline"
	"github.com/jackc/pgx/v5"
)

func Connect() (*pgx.Conn, error) {
	db_url := os.Getenv("DB_URL")

	conn, err := pgx.Connect(context.Background(), db_url)

	return conn, err
}

func InsertCoinRecord(conn *pgx.Conn, record pipeline.CoinRecord) error {
	_, err := conn.Exec(context.Background(), "INSERT INTO coin_prices (coin_id, symbol, price_usd, change_24h, sentiment) VALUES ($1, $2, $3, $4, $5)", record.CoinID, record.Symbol, record.PriceUSD, record.Change24h, record.Sentiment)
	if err != nil {
		fmt.Println("Error while inserting coin data!", err)
	}
	return err
}
