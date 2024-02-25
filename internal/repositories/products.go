package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/wildanfaz/store-management/internal/helpers"
	"github.com/wildanfaz/store-management/internal/models"
)

type ImplementProducts struct {
	dbPostgreSQL *pgx.Conn
}

type Products interface {
	AddProduct(ctx context.Context, payload models.Product) error
	GetProduct(ctx context.Context, id int) (*models.Product, error)
	ListProducts(ctx context.Context, payload models.Product) ([]models.Product, error)
	UpdateProduct(ctx context.Context, id int, payload models.Product) error
	DeleteProduct(ctx context.Context, id int) error
}

func NewProductsRepository(dbPostgreSQL *pgx.Conn) Products {
	return &ImplementProducts{
		dbPostgreSQL: dbPostgreSQL,
	}
}

func (r *ImplementProducts) AddProduct(ctx context.Context, payload models.Product) error {
	var query = `
	INSERT INTO products
	(name, description, price, quantity, user_id)
	VALUES
	($1, $2, $3, $4, $5)
	`

	_, err := r.dbPostgreSQL.Exec(ctx, query,
		payload.Name, payload.Description,
		payload.Price, payload.Quantity,
		payload.UserId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *ImplementProducts) GetProduct(ctx context.Context, id int) (*models.Product, error) {
	var product models.Product

	var query = `
	SELECT id, name, description, price, quantity, created_at, updated_at
	FROM products
	WHERE id = $1
	`

	err := r.dbPostgreSQL.QueryRow(ctx, query, id).Scan(
		&product.Id, &product.Name, &product.Description,
		&product.Price, &product.Quantity, &product.CreatedAt,
		&product.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	product.ToLocal()

	return &product, nil
}

func (r *ImplementProducts) ListProducts(ctx context.Context, payload models.Product) ([]models.Product, error) {
	var products []models.Product

	query, values := helpers.PostgreSQLQueryListProducts(payload)

	rows, err := r.dbPostgreSQL.Query(ctx, query, values...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var product models.Product

		err := rows.Scan(
			&product.Id, &product.Name, &product.Description,
			&product.Price, &product.Quantity, &product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		product.ToLocal()

		products = append(products, product)
	}

	return products, nil
}

func (r *ImplementProducts) UpdateProduct(ctx context.Context, id int, payload models.Product) error {
	query, values := helpers.PostgreSQLQueryUpdateProduct(id, payload)

	_, err := r.dbPostgreSQL.Exec(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (r *ImplementProducts) DeleteProduct(ctx context.Context, id int) error {
	_, err := r.dbPostgreSQL.Exec(ctx, "DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
