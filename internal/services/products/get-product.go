package products

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wildanfaz/store-management/internal/helpers"
)

func (s *ImplementServices) GetProduct(c *gin.Context) {
	var (
		resp = helpers.NewResponse()
		id   = c.Param("id")
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

	product, err := s.productsRepo.GetProduct(c.Request.Context(), idInt)
	if err != nil {
		s.log.Errorln("Unable to get product:", err)
		c.JSON(http.StatusInternalServerError, resp.AsError().
			WithMessage("Unable to get product"))
		return
	}

	if product == nil {
		s.log.Errorln("Product not found")
		c.JSON(http.StatusNotFound, resp.AsError().
			WithMessage("Product not found"))
		return
	}

	c.JSON(http.StatusOK, resp.WithMessage("Get product success").
		WithData(product))
}
