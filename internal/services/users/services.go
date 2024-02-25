package users

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/wildanfaz/store-management/configs"
	"github.com/wildanfaz/store-management/internal/repositories"
)

type ImplementServices struct {
	log       *logrus.Logger
	usersRepo repositories.Users
	conf      *configs.Config
}

type Services interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Profile(c *gin.Context)
	ResetPassword(c *gin.Context)
	Logout(c *gin.Context)
}

func NewServices(log *logrus.Logger, usersRepo repositories.Users, conf *configs.Config) Services {
	return &ImplementServices{
		log:       log,
		usersRepo: usersRepo,
		conf:      conf,
	}
}
