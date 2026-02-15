package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"tgfin/internal/domain/usecase"
)

type Client struct {
	Pool *pgxpool.Pool
}

func (c *Client) WithinTx(ctx context.Context, fn func(ctx context.Context, exec user.Exec) error) error {
	tx, err := c.Pool.BeginTx(ctx, pgx.TxOptions{})

	if err != nil{
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	exec := NewExecFromTx(tx)

	if err := fn(ctx, exec); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
