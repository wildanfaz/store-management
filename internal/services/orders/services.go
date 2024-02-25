package orders

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/wildanfaz/store-management/internal/repositories"
)

type ImplementServices struct {
	log        *logrus.Logger
	usersRepo  repositories.Users
	ordersRepo repositories.Orders
}

type Services interface {
	AddOrder(c *gin.Context)
	ListOrders(c *gin.Context)
	UpdateOrderStatus(c *gin.Context)
}

func NewServices(log *logrus.Logger, usersRepo repositories.Users, ordersRepo repositories.Orders) Services {
	return &ImplementServices{
		log:        log,
		usersRepo:  usersRepo,
		ordersRepo: ordersRepo,
	}
}
