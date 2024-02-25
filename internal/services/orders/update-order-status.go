package orders

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/wildanfaz/store-management/internal/helpers"
	"github.com/wildanfaz/store-management/internal/models"
)

func (s *ImplementServices) UpdateOrderStatus(c *gin.Context) {
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

	orderId := c.Param("id")

	orderIdInt, err := strconv.Atoi(orderId)
	if err != nil {
		s.log.Errorln("Unable to convert string to int:", err)
		c.JSON(http.StatusBadRequest, resp.AsError().
			WithMessage("Something went wrong"))
		return
	}

	order.Id = orderIdInt
	order.UserId = profile.Id

	err = s.ordersRepo.UpdateOrderStatus(c.Request.Context(), order)
	if err != nil && err != pgx.ErrNoRows {
		if err.Error() == "Order already done" {
			s.log.Errorln("Unable to update order status: order already done")
			c.JSON(http.StatusBadRequest, resp.AsError().
				WithMessage("Order already done"))
			return
		}

		s.log.Errorln("Unable to update order status:", err)
		c.JSON(http.StatusInternalServerError, resp.AsError().
			WithMessage("Something went wrong"))
		return
	}

	if err == pgx.ErrNoRows {
		s.log.Errorln("Unable to update order status: order not found")
		c.JSON(http.StatusBadRequest, resp.AsError().
			WithMessage("Order not found"))
		return
	}

	c.JSON(http.StatusOK, resp.WithMessage("Update order status success"))
}
