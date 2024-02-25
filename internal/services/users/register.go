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

func (s *ImplementServices) Register(c *gin.Context) {
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
		validation.Field(&user.FullName, validation.Required),
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Password, validation.Required, validation.Length(8, 64)),
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

	if profile != nil {
		s.log.Errorln("Email already registered")
		c.JSON(http.StatusBadRequest, resp.AsError().
			WithMessage("Email already registered"))
		return
	}

	hashedPassword, err := pkg.HashPassword(user.Password)
	if err != nil {
		s.log.Errorln("Unable to hash password:", err)
		c.JSON(http.StatusInternalServerError, resp.AsError().
			WithMessage("Unable to hash password"))
		return
	}

	user.Password = hashedPassword

	err = s.usersRepo.Register(c.Request.Context(), user)
	if err != nil {
		s.log.Errorln("Unable to register user:", err)
		c.JSON(http.StatusInternalServerError, resp.AsError().
			WithMessage("Unable to register user"))
		return
	}

	c.JSON(http.StatusOK, resp.WithMessage("Register user success"))
}
