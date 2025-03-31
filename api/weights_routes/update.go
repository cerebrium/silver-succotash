package weights_routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"hopdf.com/dao/weights"
	"hopdf.com/helpers"
	"hopdf.com/localware"
)

func UpdateWeights(c echo.Context) error {
	cc, ok := c.(*localware.LocalUserClerkDbContext)
	if !ok {
		c.Logger().Error("could not resolve cc")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "could not resolve cc",
		})
	}

	var weights weights.Weights

	err := cc.Bind(&weights)
	if err != nil {
		c.Logger().Errorf("could not bind the request: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "could not bind the request",
		})
	}

	// TODO: read in the form input, and update the
	// vals in the struct to send into the update.

	err = weights.Update(cc.Db)
	if err != nil {
		c.Logger().Errorf("could not update the weights: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "could not update the weights",
		})
	}

	// Retrieve the updated val
	updated_weights, err := weights.Read(cc.Db)
	if err != nil {
		c.Logger().Errorf("could not read the weights in the update: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "could not read the weights in the update",
		})

	}

	return helpers.Success(cc, updated_weights)
}
