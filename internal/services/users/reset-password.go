package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/wildanfaz/store-management/internal/helpers"
	"github.com/wildanfaz/store-management/internal/models"
	"github.com/wildanfaz/store-management/internal/pkg"
)

func (s *ImplementServices) ResetPassword(c *gin.Context) {
	var (
		resp          = helpers.NewResponse()
		payload       models.ResetPassword
		email, exists = c.Get("email")
	)

	if !exists {
		s.log.Errorln("Unable to get email from context")
		c.JSON(http.StatusInternalServerError, resp.AsError().
			WithMessage("Something went wrong"))
		return
	}

	err := c.ShouldBindJSON(&payload)
	if err != nil {
		s.log.Errorln("Unable to binding JSON:", err)
		c.JSON(http.StatusBadRequest, resp.AsError().
			WithMessage("Unable to binding JSON"))
		return
	}

	err = validation.ValidateStruct(&payload,
		validation.Field(&payload.OldPassword, validation.Required),
		validation.Field(&payload.NewPassword, validation.Required, validation.Length(8, 64)),
	)
	if err != nil {
		s.log.Errorln("Validate user got error:", err)
		c.JSON(http.StatusBadRequest, resp.AsError().
			WithMessage(err.Error()))
		return
	}

	profile, err := s.usersRepo.Profile(c.Request.Context(), email.(string))
	if err != nil {
		s.log.Errorln("Unable to get profile:", err)
		c.JSON(http.StatusInternalServerError, resp.AsError().
			WithMessage("Something went wrong"))
		return
	}

	err = pkg.ComparePassword(payload.OldPassword, profile.Password)
	if err != nil {
		s.log.Errorln("Compare password got error:", err)
		c.JSON(http.StatusBadRequest, resp.AsError().
			WithMessage("Invalid old password"))
		return
	}

	hashedPassword, err := pkg.HashPassword(payload.NewPassword)
	if err != nil {
		s.log.Errorln("Unable to hash password:", err)
		c.JSON(http.StatusInternalServerError, resp.AsError().
			WithMessage("Unable to hash password"))
		return
	}

	err = s.usersRepo.ResetPassword(c.Request.Context(), email.(string), hashedPassword)
	if err != nil {
		s.log.Errorln("Unable to reset password:", err)
		c.JSON(http.StatusInternalServerError, resp.AsError().
			WithMessage("Unable to reset password"))
		return
	}

	c.JSON(http.StatusOK, resp.WithMessage("Reset password success"))
}
