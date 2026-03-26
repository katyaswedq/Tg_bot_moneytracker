package pg

import (
    "context"
    "fmt"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(databaseURL string) (*pgxpool.Pool, error) {
    config, err := pgxpool.ParseConfig(databaseURL)
    if err != nil {
        return nil, fmt.Errorf("parse config: %w", err)
    }

    config.MaxConns = 25
    config.MinConns = 5
    config.MaxConnLifetime = time.Hour
    config.MaxConnIdleTime = 30 * time.Minute
    config.HealthCheckPeriod = time.Minute

    pool, err := pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        return nil, fmt.Errorf("create pool: %w", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := pool.Ping(ctx); err != nil {
        return nil, fmt.Errorf("ping failed: %w", err)
    }

    return pool, nil
}