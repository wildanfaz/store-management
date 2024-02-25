package configs

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

func InitPostgreSQL() *pgx.Conn {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		for i := 1; i <= 20; i++ {
			fmt.Printf("Retrying Database Connection #%d\n", i)
			conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

			if err == nil {
				break
			}

			time.Sleep(5 * time.Second)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
	}

	return conn
}
