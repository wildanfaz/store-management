package migrations

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func Rollback(ctx context.Context, db *pgx.Conn) {
	for _, v := range rollbackQueries {
		_, err := db.Exec(ctx, v)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Rollback Success")
}

var rollbackQueries = []string{
	`
DROP TABLE IF EXISTS orders
`,
	`
DROP TABLE IF EXISTS products
`,
	`
DROP TABLE IF EXISTS users
`,
}
