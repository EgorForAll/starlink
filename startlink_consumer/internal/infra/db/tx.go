package db

import (
	"context"
	"database/sql"
	"fmt"
)

type txKey struct{}

type TxManager struct {
	db DbConn
}

func NewTxManager(db DbConn) *TxManager {
	return &TxManager{db: db}
}

func (m *TxManager) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	ctx = context.WithValue(ctx, txKey{}, tx)

	if err := fn(ctx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// RunnerFromCtx возвращает *sql.Tx из контекста, если он есть, иначе fallback.
func RunnerFromCtx(ctx context.Context, fallback TxRunner) TxRunner {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return fallback
}
