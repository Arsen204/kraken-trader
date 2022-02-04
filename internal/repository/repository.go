package repository

import (
	"course-project/internal/domain"

	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	CreateOrder(ctx context.Context, order domain.Order) error
}

type repo struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &repo{
		pool: pool,
	}
}

const createOrderQuery = `INSERT INTO orders(order_id, order_type, symbol, side, quantity, filled, order_time, reduced_only)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

func (r *repo) CreateOrder(ctx context.Context, order domain.Order) error {
	_, err := r.pool.Exec(ctx, createOrderQuery,
		order.OrderID,
		order.Type,
		order.Symbol,
		order.Side,
		order.Quantity,
		order.Filled,
		order.Timestamp,
		order.ReducedOnly,
	)
	if err != nil {
		return err
	}

	return nil
}
