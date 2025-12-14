package weights_routes

import (
	"fmt"
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
		c.Logger().Errorf("could not bind the request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "could not bind the request",
		})
	}

	fmt.Println("Updates coming in: ", weights)

	err = weights.Update(cc.Db)
	if err != nil {
		c.Logger().Errorf("could not update the weights: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "could not update the weights",
		})
	}

	// Retrieve the updated val
	updated_weights, err := weights.Read(cc.Db)
	if err != nil {
		c.Logger().Errorf("could not read the weights in the update: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "could not read the weights in the update",
		})
	}

	fmt.Println("Updated weights: ", updated_weights)

	return helpers.Success(cc, updated_weights)
}
