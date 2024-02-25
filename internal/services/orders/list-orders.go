package orders

import (
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/wildanfaz/store-management/internal/helpers"
	"github.com/wildanfaz/store-management/internal/models"
)

func (s *ImplementServices) ListOrders(c *gin.Context) {
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

	err = c.ShouldBindQuery(&order)
	if err != nil {
		s.log.Errorln("Unable to binding query:", err)
		c.JSON(http.StatusBadRequest, resp.AsError().
			WithMessage("Unable to binding query"))
		return
	}

	err = validation.ValidateStruct(&order,
		validation.Field(&order.ProductId, validation.Min(1)),
	)
	if err != nil {
		s.log.Errorln("Validate order got error:", err)
		c.JSON(http.StatusBadRequest, resp.AsError().
			WithMessage(err.Error()))
		return
	}

	order.UserId = profile.Id

	orders, err := s.ordersRepo.ListOrders(c.Request.Context(), order)
	if err != nil {
		s.log.Errorln("List orders got error:", err)
		c.JSON(http.StatusInternalServerError, resp.AsError().
			WithMessage("Something went wrong"))
		return
	}

	c.JSON(http.StatusOK, resp.WithMessage("Get list orders success").
		WithData(orders))
}
