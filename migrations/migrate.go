package migrations

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func Migrate(ctx context.Context, db *pgx.Conn) {
	for _, v := range migrateQueries {
		_, err := db.Exec(ctx, v)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Migration Success")
}

var migrateQueries = []string{
	`
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql
`,
	`
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL NOT NULL,
    full_name varchar(255) NOT NULL,
    email varchar(255) NOT NULL UNIQUE,
    password varchar(255) NOT NULL,
	is_login boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
)
`,
	`
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (email)
`,
	`
DROP TRIGGER IF EXISTS users_updated_at_trigger ON users
`,
	`
CREATE TRIGGER users_updated_at_trigger
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at()
`,
	`
DO $$ 
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_enum') THEN
		CREATE TYPE status_enum AS ENUM ('diproses', 'selesai');
	END IF;
END $$;
`,
	`
CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL NOT NULL,
    name varchar(255) NOT NULL,
    description text,
    price int NOT NULL,
    quantity int NOT NULL,
	user_id bigint NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
	FOREIGN KEY (user_id) REFERENCES users(id)
)
`,
	`
CREATE INDEX IF NOT EXISTS idx_products_name ON products (name)
`,
	`
DROP TRIGGER IF EXISTS products_updated_at_trigger ON products
`,
	`
CREATE TRIGGER products_updated_at_trigger
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION update_updated_at()
`,
	`
CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL NOT NULL,
    product_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    quantity int NOT NULL,
    status status_enum NOT NULL DEFAULT 'diproses',
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
	FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
)
`,
	`
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders (user_id)
`,
	`
DROP TRIGGER IF EXISTS orders_updated_at_trigger ON orders
`,
	`
CREATE TRIGGER orders_updated_at_trigger
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION update_updated_at()
`,
}
