package weights_routes

import (
	"os"

	"github.com/labstack/echo/v4"
	"hopdf.com/dao/weights"
	"hopdf.com/helpers"
	"hopdf.com/localware"
)

func ReadWeights(c echo.Context) error {
	cc, ok := c.(*localware.LocalUserClerkDbContext)
	if !ok {
		c.Logger().Error("could not resolve cc")
		os.Exit(1)
	}
	weights := weights.Weights{
		ID: 1,
	}

	updated_weights, err := weights.Read(cc.DbContext.Db)
	if err != nil {
		c.Logger().Errorf("could not read the weights: ", err)
		os.Exit(1)

	}

	return helpers.Success(cc, updated_weights)
}
