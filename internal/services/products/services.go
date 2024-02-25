package products

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/wildanfaz/store-management/internal/repositories"
)

type ImplementServices struct {
	log  *logrus.Logger
	productsRepo repositories.Products
	usersRepo  repositories.Users
}

type Services interface {
	AddProduct(c *gin.Context)
	GetProduct(c *gin.Context)
	ListProducts(c *gin.Context)
	UpdateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
}

func NewServices(log *logrus.Logger, productsRepo repositories.Products, usersRepo repositories.Users) Services {
	return &ImplementServices{
		log:  log,
		productsRepo: productsRepo,
		usersRepo:  usersRepo,
	}
}
