package products

import (
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/wildanfaz/store-management/internal/helpers"
	"github.com/wildanfaz/store-management/internal/models"
)

func (s *ImplementServices) AddProduct(c *gin.Context) {
	var (
		resp          = helpers.NewResponse()
		product       models.Product
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

	product.UserId = profile.Id

	err = c.ShouldBindJSON(&product)
	if err != nil {
		s.log.Errorln("Unable to binding JSON:", err)
		c.JSON(http.StatusBadRequest, resp.AsError().
			WithMessage("Unable to binding JSON"))
		return
	}

	err = validation.ValidateStruct(&product,
		validation.Field(&product.Name, validation.Required),
		validation.Field(&product.Price, validation.Required),
		validation.Field(&product.Quantity, validation.Required, validation.Min(1)),
	)
	if err != nil {
		s.log.Errorln("Validate product got error:", err)
		c.JSON(http.StatusBadRequest, resp.AsError().
			WithMessage(err.Error()))
		return
	}

	err = s.productsRepo.AddProduct(c.Request.Context(), product)
	if err != nil {
		s.log.Errorln("Unable to add product:", err)
		c.JSON(http.StatusInternalServerError, resp.AsError().
			WithMessage("Unable to add product"))
		return
	}

	c.JSON(http.StatusOK, resp.WithMessage("Add product success"))
}
