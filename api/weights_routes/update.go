package weights_routes

import (
	"os"

	"github.com/labstack/echo/v4"
	"hopdf.com/dao/weights"
	"hopdf.com/helpers"
	"hopdf.com/localware"
)

func UpdateWeights(c echo.Context) error {
	cc, ok := c.(*localware.LocalUserClerkDbContext)
	if !ok {
		c.Logger().Error("could not resolve cc")
		os.Exit(1)
	}

	weights := weights.Weights{
		ID: 1,
	}

	// TODO: read in the form input, and update the
	// vals in the struct to send into the update.

	err := weights.Update(cc.DbContext.Db)
	if err != nil {
		c.Logger().Errorf("could not update the weights: ", err)
		os.Exit(1)
	}

	// Retrieve the updated val
	updated_weights, err := weights.Read(cc.DbContext.Db)
	if err != nil {
		c.Logger().Errorf("could not read the weights in the update: ", err)
		os.Exit(1)

	}

	return helpers.Success(cc, updated_weights)
}
