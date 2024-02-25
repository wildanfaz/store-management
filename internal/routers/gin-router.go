package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/wildanfaz/store-management/configs"
	"github.com/wildanfaz/store-management/internal/middlewares"
	"github.com/wildanfaz/store-management/internal/pkg"
	"github.com/wildanfaz/store-management/internal/repositories"
	"github.com/wildanfaz/store-management/internal/services/orders"
	"github.com/wildanfaz/store-management/internal/services/products"
	"github.com/wildanfaz/store-management/internal/services/users"
)

func InitGinRouter() {
	// configs
	conf := configs.InitConfig()
	dbPostgreSql := configs.InitPostgreSQL()

	// pkg
	log := pkg.InitLogger()

	// repositories
	productsRepo := repositories.NewProductsRepository(dbPostgreSql)
	usersRepo := repositories.NewUsersRepository(dbPostgreSql)
	ordersRepo := repositories.NewOrdersRepository(dbPostgreSql)

	// services
	products := products.NewServices(log, productsRepo, usersRepo)
	users := users.NewServices(log, usersRepo, conf)
	orders := orders.NewServices(log, usersRepo, ordersRepo)

	r := gin.Default()
	apiV1 := r.Group("/api/v1")

	apiV1.POST("/users/register", users.Register)
	apiV1.POST("/users/login", users.Login)

	apiV1.Use(middlewares.Auth(log, conf.JWTSecret, usersRepo))

	apiV1.GET("/users/profile", users.Profile)
	apiV1.PUT("/users/reset-password", users.ResetPassword)
	apiV1.POST("/users/logout", users.Logout)

	apiV1.POST("/products", products.AddProduct)
	apiV1.GET("/products/:id", products.GetProduct)
	apiV1.GET("/products", products.ListProducts)
	apiV1.PUT("/products/:id", products.UpdateProduct)
	apiV1.DELETE("/products/:id", products.DeleteProduct)

	apiV1.POST("/orders", orders.AddOrder)
	apiV1.GET("/orders", orders.ListOrders)
	apiV1.PUT("/orders/:id", orders.UpdateOrderStatus)

	r.Run() // listen and serve on 0.0.0.0:8080
}
