package weights_routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"hopdf.com/dao/weights"
	"hopdf.com/helpers"
	"hopdf.com/localware"
)

func ReadWeights(c echo.Context) error {
	cc, ok := c.(*localware.LocalUserClerkDbContext)
	if !ok {
		c.Logger().Error("could not resolve cc")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "could not resolve cc",
		})
	}

	weights := weights.Weights{
		ID: 1,
	}

	updated_weights, err := weights.Read(cc.Db)
	if err != nil {
		c.Logger().Errorf("could not read the weights: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "could not read the weights",
		})
	}

	return helpers.Success(cc, updated_weights)
}
