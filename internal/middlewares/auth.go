package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/wildanfaz/store-management/internal/helpers"
	"github.com/wildanfaz/store-management/internal/pkg"
	"github.com/wildanfaz/store-management/internal/repositories"
)

func Auth(log *logrus.Logger, jwtSecret []byte, usersRepo repositories.Users) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			resp        = helpers.NewResponse()
			bearerToken = c.GetHeader("Authorization")
		)

		if len(strings.Split(bearerToken, " ")) < 2 {
			log.Errorln("Bearer token not found")
			c.JSON(http.StatusUnauthorized, resp.AsError().
				WithMessage("Unauthorized"))
			c.Abort()
			return
		}

		token := strings.Split(bearerToken, " ")[1]

		claims, err := pkg.ValidateToken(token, jwtSecret)
		if err != nil {
			log.Errorln("Validate token got error:", err)
			c.JSON(http.StatusUnauthorized, resp.AsError().
				WithMessage("Unauthorized"))
			c.Abort()
			return
		}

		isLogin, err := usersRepo.IsLogin(c.Request.Context(), claims.Email)
		if err != nil {
			log.Errorln("Check is login got error:", err)
			c.JSON(http.StatusUnauthorized, resp.AsError().
				WithMessage("Unauthorized"))
			c.Abort()
			return
		}

		if !isLogin {
			log.Errorln("User is not login")
			c.JSON(http.StatusUnauthorized, resp.AsError().
				WithMessage("Unauthorized"))
			c.Abort()
			return
		}

		c.Set("email", claims.Email)

		c.Next()
	}
}
