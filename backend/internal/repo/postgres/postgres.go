package postgres

import (
	"context"
	"errors"
	"github.com/doug-martin/goqu/v9"
	"github.com/ilam072/wbtech-l0/backend/internal/repo"
	"github.com/ilam072/wbtech-l0/backend/internal/types/domain"
	"github.com/ilam072/wbtech-l0/backend/pkg/e"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepo struct {
	pool *pgxpool.Pool
}

func NewOrderRepo(db *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{pool: db}
}

func (r *OrderRepo) CreateOrder(
	ctx context.Context,
	order domain.Order,
	delivery domain.Delivery,
	payment domain.Payment,
	items []domain.Item,
) error {
	const op = "postgres.CreateOrder()"

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return e.Wrap(op, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	orderQuery, args, err := goqu.Insert("orders").Rows(order).ToSQL()
	if err != nil {
		return e.Wrap(op, err)
	}
	if _, err = tx.Exec(ctx, orderQuery, args...); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "orders_pkey" {
				return e.Wrap(op, repo.ErrOrderExists)
			}
		}
		return e.Wrap(op, err)
	}

	deliveryQuery, args, err := goqu.Insert("delivery").Rows(delivery).ToSQL()
	if err != nil {
		return e.Wrap(op, err)
	}
	if _, err = tx.Exec(ctx, deliveryQuery, args...); err != nil {
		return e.Wrap(op, err)
	}

	paymentQuery, args, err := goqu.Insert("payment").Rows(payment).ToSQL()
	if err != nil {
		return e.Wrap(op, err)
	}
	if _, err = tx.Exec(ctx, paymentQuery, args...); err != nil {
		return e.Wrap(op, err)
	}

	itemsQuery, args, err := goqu.Insert("items").Rows(items).ToSQL()
	if err != nil {
		return e.Wrap(op, err)
	}
	if _, err = tx.Exec(ctx, itemsQuery, args...); err != nil {
		return e.Wrap(op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return e.Wrap(op, err)
	}

	return nil
}
