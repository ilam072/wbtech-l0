package postgres

import (
	"context"
	"errors"
	"github.com/doug-martin/goqu/v9"
	"github.com/ilam072/wbtech-l0/backend/internal/repo"
	"github.com/ilam072/wbtech-l0/backend/internal/types/domain"
	"github.com/ilam072/wbtech-l0/backend/pkg/e"
	"github.com/jackc/pgx/v5"
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

func (r *OrderRepo) GetOrder(ctx context.Context, ID string) (domain.FullOrder, error) {
	const op = "postgres.GetOrder()"

	var (
		order    = domain.Order{}
		delivery = domain.Delivery{}
		payment  = domain.Payment{}
		items    []domain.Item
	)

	sql, args, err := goqu.From("orders").
		Join(goqu.T("delivery"), goqu.On(goqu.Ex{"orders.id": goqu.I("delivery.order_id")})).
		Join(goqu.T("payment"), goqu.On(goqu.Ex{"orders.id": goqu.I("payment.order_id")})).
		Select(
			// orders
			goqu.I("orders.id"),
			goqu.I("orders.track_number"),
			goqu.I("orders.entry"),
			goqu.I("orders.locale"),
			goqu.I("orders.internal_signature"),
			goqu.I("orders.customer_id"),
			goqu.I("orders.delivery_service"),
			goqu.I("orders.shardkey"),
			goqu.I("orders.sm_id"),
			goqu.I("orders.date_created"),
			goqu.I("orders.oof_shard"),

			// delivery
			goqu.I("delivery.id"),
			goqu.I("delivery.order_id"),
			goqu.I("delivery.name"),
			goqu.I("delivery.phone"),
			goqu.I("delivery.zip"),
			goqu.I("delivery.city"),
			goqu.I("delivery.address"),
			goqu.I("delivery.region"),
			goqu.I("delivery.email"),

			// payment
			goqu.I("payment.transaction"),
			goqu.I("payment.order_id"),
			goqu.I("payment.request_id"),
			goqu.I("payment.currency"),
			goqu.I("payment.provider"),
			goqu.I("payment.amount"),
			goqu.I("payment.payment_dt"),
			goqu.I("payment.bank"),
			goqu.I("payment.delivery_cost"),
			goqu.I("payment.goods_total"),
			goqu.I("payment.custom_fee"),
		).
		Where(goqu.Ex{"orders.id": ID}).
		ToSQL()
	if err != nil {
		return domain.FullOrder{}, e.Wrap(op, err)
	}

	row := r.pool.QueryRow(ctx, sql, args...)
	if err != nil {

		return domain.FullOrder{}, e.Wrap(op, err)
	}

	if err := row.Scan(
		&order.ID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,

		&delivery.ID,
		&delivery.OrderID,
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email,

		&payment.Transaction,
		&payment.OrderID,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDt,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.FullOrder{}, e.Wrap(op, repo.ErrOrderNotFound)
		}
		return domain.FullOrder{}, e.Wrap(op, err)
	}

	sql, args, err = goqu.From("items").
		Where(goqu.Ex{"order_id": ID}).ToSQL()
	if err != nil {
		return domain.FullOrder{}, e.Wrap(op, err)
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return domain.FullOrder{}, e.Wrap(op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var itm domain.Item

		err = rows.Scan(
			&itm.ChrtID,
			&itm.OrderID,
			&itm.TrackNumber,
			&itm.Price,
			&itm.Rid,
			&itm.Name,
			&itm.Sale,
			&itm.Size,
			&itm.TotalPrice,
			&itm.NmID,
			&itm.Brand,
			&itm.Status,
		)
		if err != nil {
			return domain.FullOrder{}, e.Wrap(op, err)
		}

		items = append(items, itm)
	}

	return domain.FullOrder{
		Order:    order,
		Delivery: delivery,
		Payment:  payment,
		Items:    items,
	}, nil
}
