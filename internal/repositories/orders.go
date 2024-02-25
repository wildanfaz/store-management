package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/wildanfaz/store-management/internal/constants"
	"github.com/wildanfaz/store-management/internal/helpers"
	"github.com/wildanfaz/store-management/internal/models"
)

type ImplementOrders struct {
	dbPostgreSQL *pgx.Conn
}

type Orders interface {
	AddOrder(ctx context.Context, payload models.Order) error
	ListOrders(ctx context.Context, payload models.Order) (*models.Orders, error)
	UpdateOrderStatus(ctx context.Context, payload models.Order) error
}

func NewOrdersRepository(dbPostgreSQL *pgx.Conn) Orders {
	return &ImplementProducts{
		dbPostgreSQL: dbPostgreSQL,
	}
}

func (r *ImplementProducts) AddOrder(ctx context.Context, payload models.Order) error {
	var query = `
	INSERT INTO orders
	(product_id, user_id, quantity)
	VALUES
	($1, $2, $3)
	`

	_, err := r.dbPostgreSQL.Exec(ctx, query,
		payload.ProductId, payload.UserId, payload.Quantity,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *ImplementProducts) ListOrders(ctx context.Context, payload models.Order) (*models.Orders, error) {
	var orders models.Orders

	query, values := helpers.PostgreSQLQueryListOrders(payload)

	rows, err := r.dbPostgreSQL.Query(ctx, query, values...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var order models.OrderResponse

		err := rows.Scan(
			&order.Order.Id, &order.Order.Status, &order.Order.Quantity, &order.Order.CreatedAt, &order.Order.UpdatedAt,
			&order.Product.Id, &order.Product.Name, &order.Product.Description,
			&order.Product.Price, &order.Product.Quantity, &order.Product.CreatedAt, &order.Product.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		order.ToLocal()

		orders = append(orders, order)
	}

	return &orders, nil
}

func (r *ImplementProducts) UpdateOrderStatus(ctx context.Context, payload models.Order) error {
	var queries = []string{
		`
	SELECT status
	FROM orders
	WHERE id = $1
		`,
		`
	SELECT product_id
	FROM orders
	WHERE id = $1
		`,
		`
	SELECT quantity
	FROM orders
	WHERE id = $1 AND user_id = $2
	`,
		`
	UPDATE products
	SET quantity = quantity - $1
	WHERE id = $2 AND user_id = $3
	`,
		`
	UPDATE orders
	SET status = 'selesai'
	WHERE id = $1
	`,
	}

	tx, err := r.dbPostgreSQL.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return err
	}

	for i, v := range queries {
		if i == 0 {
			err := tx.QueryRow(ctx, v, payload.Id).Scan(&payload.Status)
			if err != nil {
				tx.Rollback(ctx)
				return err
			}

			if payload.Status == constants.Selesai {
				tx.Rollback(ctx)
				return errors.New("Order already done")
			}
		} else if i == 1 {
			err := tx.QueryRow(ctx, v, payload.Id).Scan(&payload.ProductId)
			if err != nil {
				tx.Rollback(ctx)
				return err
			}
		} else if i == 2 {
			err := tx.QueryRow(ctx, v, payload.Id, payload.UserId).Scan(&payload.Quantity)
			if err != nil {
				tx.Rollback(ctx)
				return err
			}
		} else if i == 3 {
			_, err := tx.Exec(ctx, v, payload.Quantity, payload.ProductId, payload.UserId)
			if err != nil {
				tx.Rollback(ctx)
				return err
			}
		} else {
			_, err := tx.Exec(ctx, v, payload.Id)
			if err != nil {
				tx.Rollback(ctx)
				return err
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return nil
}
