package db

import (
	"context"
	"database/sql"
	"fmt"
)

type DbConn interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row

	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row

	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Close() error
}

// TxRunner — подмножество методов, которые реализуют и *sql.DB, и *sql.Tx.
// Репозитории используют его, чтобы работать как внутри транзакции, так и без неё.
type TxRunner interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func InitDb(url string) (DbConn, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}
	
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db ping error: %v", err)
	}

	return db, nil
}
