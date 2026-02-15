package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"tgfin/internal/domain/usecase"
)

type exec struct {
	q querier
}

type querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

func NewExecFromPool(pool *pgxpool.Pool) user.Exec {
	return &exec{q: pool}
}

func NewExecFromTx(tx pgx.Tx) user.Exec {
	return &exec{q: tx}
}

func (e *exec) QueryRow(ctx context.Context, sql string, args ...any) user.Row {
	return e.q.QueryRow(ctx, sql, args...)
}

func (e *exec) Query(ctx context.Context, sql string, args ...any) (user.Rows, error) {
	return e.q.Query(ctx, sql, args...)
}

func (e *exec) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := e.q.Exec(ctx, sql, args...)
	return err
}

func (e *exec) ExecWithResult(ctx context.Context, sql string, args ...any) (user.CommandTag, error) {
	return e.q.Exec(ctx, sql, args...)
}
