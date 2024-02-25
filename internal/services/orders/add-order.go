package orders

import (
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/wildanfaz/store-management/internal/helpers"
	"github.com/wildanfaz/store-management/internal/models"
)

func (s *ImplementServices) AddOrder(c *gin.Context) {
	var (
		resp          = helpers.NewResponse()
		order         models.Order
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
			WithMessage("Something went wrong"))
		return
	}

	err = c.ShouldBindJSON(&order)
	if err != nil {
		s.log.Errorln("Unable to binding JSON:", err)
		c.JSON(http.StatusBadRequest, resp.AsError().
			WithMessage("Unable to binding JSON"))
		return
	}

	err = validation.ValidateStruct(&order,
		validation.Field(&order.ProductId, validation.Required, validation.Min(1)),
		validation.Field(&order.Quantity, validation.Required, validation.Min(1)),
	)
	if err != nil {
		s.log.Errorln("Validate order got error:", err)
		c.JSON(http.StatusBadRequest, resp.AsError().
			WithMessage(err.Error()))
		return
	}

	order.UserId = profile.Id

	err = s.ordersRepo.AddOrder(c.Request.Context(), order)
	if err != nil {
		s.log.Errorln("Unable to add order:", err)
		c.JSON(http.StatusInternalServerError, resp.AsError().
			WithMessage("Unable to add order"))
		return
	}

	c.JSON(http.StatusOK, resp.WithMessage("Add order success"))
}
