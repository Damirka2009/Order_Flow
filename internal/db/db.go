package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func NewDB(ctx context.Context) (*pgx.Conn, error) {
	connStr := "postgres://admin:secret@localhost:5432/orders_db"
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to PostgreSQL")
	return conn, nil
}
