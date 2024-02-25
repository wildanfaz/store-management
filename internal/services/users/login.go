package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/wildanfaz/store-management/internal/helpers"
	"github.com/wildanfaz/store-management/internal/models"
	"github.com/wildanfaz/store-management/internal/pkg"
)

func (s *ImplementServices) Login(c *gin.Context) {
	var (
		resp = helpers.NewResponse()
		user models.User
	)

	err := c.ShouldBindJSON(&user)
	if err != nil {
		s.log.Errorln("Unable to binding JSON:", err)
		c.JSON(http.StatusBadRequest, resp.AsError().
			WithMessage("Unable to binding JSON"))
		return
	}

	err = validation.ValidateStruct(&user,
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Password, validation.Required),
	)
	if err != nil {
		s.log.Errorln("Validate user got error:", err)
		c.JSON(http.StatusBadRequest, resp.AsError().
			WithMessage(err.Error()))
		return
	}

	profile, err := s.usersRepo.Profile(c.Request.Context(), user.Email)
	if err != nil {
		s.log.Errorln("Unable to get profile:", err)
		c.JSON(http.StatusInternalServerError, resp.AsError().
			WithMessage("Something went wrong"))
		return
	}

	if profile == nil {
		s.log.Errorln("Email or password is incorrect")
		c.JSON(http.StatusBadRequest, resp.AsError().
			WithMessage("Email or password is incorrect"))
		return
	}

	err = pkg.ComparePassword(user.Password, profile.Password)
	if err != nil {
		s.log.Errorln("Compare password got error:", err)
		c.JSON(http.StatusBadRequest, resp.AsError().
			WithMessage("Email or password is incorrect"))
		return
	}

	token, err := pkg.GenerateToken(&pkg.NewClaims{
		Email: profile.Email,
	}, s.conf.JWTSecret)
	if err != nil {
		s.log.Errorln("Generate token got error:", err)
		c.JSON(http.StatusInternalServerError, resp.AsError().
			WithMessage("Something went wrong"))
		return
	}

	err = s.usersRepo.UpdateIsLogin(c.Request.Context(), profile.Email, true)
	if err != nil {
		s.log.Errorln("Unable to update is_login:", err)
		c.JSON(http.StatusInternalServerError, resp.AsError().
			WithMessage("Something went wrong"))
		return
	}

	c.JSON(http.StatusOK, resp.WithMessage("Login success").
		WithData(map[string]interface{}{
			"token": token,
		}))
}
