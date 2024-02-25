package products

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wildanfaz/store-management/internal/helpers"
	"github.com/wildanfaz/store-management/internal/models"
)

func (s *ImplementServices) UpdateProduct(c *gin.Context) {
	var (
		resp    = helpers.NewResponse()
		id      = c.Param("id")
		product models.Product
	)

	idInt, err := strconv.Atoi(id)
	if err != nil {
		s.log.Errorln("Unable to convert string to int:", err)
		c.JSON(http.StatusBadRequest, resp.AsError().
			WithMessage("Unable to convert string to int"))
		return
	}

	if idInt < 1 {
		s.log.Errorln("Id must be greater than 0")
		c.JSON(http.StatusBadRequest, resp.AsError().
			WithMessage("Id must be greater than 0"))
		return
	}

	err = c.ShouldBindJSON(&product)
	if err != nil {
		s.log.Errorln("Unable to binding JSON:", err)
		c.JSON(http.StatusBadRequest, resp.AsError().
			WithMessage("Unable to binding JSON"))
		return
	}

	if reflect.ValueOf(product).IsZero() {
		s.log.Errorln("Product cannot be empty")
		c.JSON(http.StatusBadRequest, resp.AsError().
			WithMessage("Product cannot be empty"))
		return
	}

	if product.Price < 0 {
		s.log.Errorln("Price must not be less than 0")
		c.JSON(http.StatusBadRequest, resp.AsError().
			WithMessage("Price must not be less than 0"))
		return
	}

	if product.Quantity < 0 {
		s.log.Errorln("Quantity must not be less than 0")
		c.JSON(http.StatusBadRequest, resp.AsError().
			WithMessage("Quantity must not be less than 0"))
		return
	}

	err = s.productsRepo.UpdateProduct(c.Request.Context(), idInt, product)
	if err != nil {
		s.log.Errorln("Unable to update product:", err)
		c.JSON(http.StatusInternalServerError, resp.AsError().
			WithMessage("Unable to update product"))
		return
	}

	c.JSON(http.StatusOK, resp.WithMessage("Update product success"))
}
