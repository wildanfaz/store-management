package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wildanfaz/store-management/internal/helpers"
)

func (s *ImplementServices) Logout(c *gin.Context) {
	var (
		resp          = helpers.NewResponse()
		email, exists = c.Get("email")
	)

	if !exists {
		s.log.Errorln("Unable to get email from context")
		c.JSON(http.StatusInternalServerError, resp.AsError().
			WithMessage("Something went wrong"))
		return
	}

	err := s.usersRepo.UpdateIsLogin(c.Request.Context(), email.(string), false)
	if err != nil {
		s.log.Errorln("Unable to update is_login:", err)
		c.JSON(http.StatusInternalServerError, resp.AsError().
			WithMessage("Unable to logout"))
		return
	}

	c.JSON(http.StatusOK, resp.WithMessage("Logout success"))
}
