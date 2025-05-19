package tiers_routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"hopdf.com/dao/tiers"
	"hopdf.com/helpers"
	"hopdf.com/localware"
)

func ReadTiers(c echo.Context) error {
	cc, ok := c.(*localware.LocalUserClerkDbContext)
	if !ok {
		c.Logger().Error("could not resolve cc")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "could not resolve cc",
		})
	}

	tiers, err := tiers.ReadTiers(cc.Db)
	if err != nil {
		cc.Logger().Error("Could not read tiers: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Could not fetch tiers",
		})
	}

	return helpers.Success(cc, tiers)
}
