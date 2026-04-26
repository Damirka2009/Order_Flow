package repository

import (
	"context"
	"master/internal/domain"

	"github.com/jackc/pgx/v5"
)

type OrdersRepository struct {
	db *pgx.Conn
}

func New(db *pgx.Conn) *OrdersRepository {
	return &OrdersRepository{
		db: db,
	}
}

func (r *OrdersRepository) Save(ctx context.Context, order *domain.Order) error {
	query := `
		INSERT INTO orders (id, item, category, currency, price, quantity)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query,
		order.Id,
		order.Item,
		order.Category,
		order.Currency,
		order.Price,
		order.Quantity,
	)

	return err
}

func (r *OrdersRepository) Get(ctx context.Context, id string) (*domain.Order, bool) {
	query := `
		SELECT id, item, category, currency, price, quantity
		FROM orders
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var order domain.Order
	err := row.Scan(
		&order.Id,
		&order.Item,
		&order.Category,
		&order.Currency,
		&order.Price,
		&order.Quantity,
	)

	if err != nil {
		return nil, false
	}

	return &order, true
}

func (r *OrdersRepository) Update(ctx context.Context, order *domain.Order) error {
	query := `
		UPDATE orders
		SET item = $1,
			category = $2,
			currency = $3,
			price = $4,
			quantity = $5
		WHERE id = $6
	`

	_, err := r.db.Exec(ctx, query,
		order.Item,
		order.Category,
		order.Currency,
		order.Price,
		order.Quantity,
		order.Id,
	)

	return err
}

func (r *OrdersRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM orders WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *OrdersRepository) List(ctx context.Context) []*domain.Order {
	query := `
		SELECT id, item, category, currency, price, quantity
		FROM orders
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var orders []*domain.Order

	for rows.Next() {
		var o domain.Order
		err := rows.Scan(
			&o.Id,
			&o.Item,
			&o.Category,
			&o.Currency,
			&o.Price,
			&o.Quantity,
		)
		if err != nil {
			return nil
		}
		orders = append(orders, &o)
	}

	return orders
}
