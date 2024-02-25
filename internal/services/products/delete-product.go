package products

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wildanfaz/store-management/internal/helpers"
)

func (s *ImplementServices) DeleteProduct(c *gin.Context) {
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

	err = s.productsRepo.DeleteProduct(c.Request.Context(), idInt)
	if err != nil {
		s.log.Errorln("Unable to delete product:", err)
		c.JSON(http.StatusInternalServerError, resp.AsError().
			WithMessage("Unable to delete product"))
		return
	}

	c.JSON(http.StatusOK, resp.WithMessage("Delete product success"))
}
