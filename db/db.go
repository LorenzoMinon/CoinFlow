package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func Connect() (*pgx.Conn, error) {
	db_url := os.Getenv("DB_URL")

	conn, err := pgx.Connect(context.Background(), db_url)

	return conn, err
}
