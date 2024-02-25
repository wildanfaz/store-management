package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/wildanfaz/store-management/configs"
	"github.com/wildanfaz/store-management/internal/routers"
	"github.com/wildanfaz/store-management/migrations"
)

func InitCmd(ctx context.Context) {
	var rootCmd = &cobra.Command{
		Short: "Store Management",
	}

	rootCmd.AddCommand(startGinApp)
	rootCmd.AddCommand(migrate)
	rootCmd.AddCommand(rollback)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		panic(err)
	}
}

var startGinApp = &cobra.Command{
	Use:   "start",
	Short: "Start the application",
	Run: func(cmd *cobra.Command, args []string) {
		routers.InitGinRouter()
	},
}

var migrate = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate database",
	Run: func(cmd *cobra.Command, args []string) {
		dbPostgreSql := configs.InitPostgreSQL()

		migrations.Migrate(context.Background(), dbPostgreSql)
	},
}

var rollback = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback database",
	Run: func(cmd *cobra.Command, args []string) {
		dbPostgreSql := configs.InitPostgreSQL()

		migrations.Rollback(context.Background(), dbPostgreSql)
	},
}
