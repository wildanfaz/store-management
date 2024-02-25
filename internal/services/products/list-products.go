package products

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wildanfaz/store-management/internal/helpers"
	"github.com/wildanfaz/store-management/internal/models"
)

func (s *ImplementServices) ListProducts(c *gin.Context) {
	var (
		resp    = helpers.NewResponse()
		product models.Product
	)

	err := c.ShouldBindQuery(&product)
	if err != nil {
		s.log.Errorln("Unable to binding Query:", err)
		c.JSON(http.StatusBadRequest, resp.AsError().
			WithMessage("Unable to binding Query"))
		return
	}

	products, err := s.productsRepo.ListProducts(c.Request.Context(), product)
	if err != nil {
		s.log.Errorln("Unable to fetch list products:", err)
		c.JSON(http.StatusInternalServerError, resp.AsError().
			WithMessage("Unable to fetch list products"))
		return
	}

	c.JSON(http.StatusOK, resp.WithMessage("Get list products success").
		WithData(products))
}
