package main

import (
	"context"

	_ "github.com/joho/godotenv/autoload"
	"github.com/wildanfaz/store-management/cmd"
)

func main() {
	cmd.InitCmd(context.Background())
}
