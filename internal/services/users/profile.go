package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wildanfaz/store-management/internal/helpers"
)

func (s *ImplementServices) Profile(c *gin.Context) {
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

	profile, err := s.usersRepo.Profile(c.Request.Context(), email.(string))
	if err != nil {
		s.log.Errorln("Unable to get profile:", err)
		c.JSON(http.StatusInternalServerError, resp.AsError().
			WithMessage("Unable to get profile"))
		return
	}

	if profile == nil {
		s.log.Errorln("Profile not found")
		c.JSON(http.StatusNotFound, resp.AsError().
			WithMessage("Unable to get profile"))
		return
	}

	c.JSON(http.StatusOK, resp.WithMessage("Get profile success").
		WithData(profile))
}
