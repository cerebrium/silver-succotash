package tiers_routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"hopdf.com/dao/tiers"
	"hopdf.com/helpers"
	"hopdf.com/localware"
)

func UpdateTiers(c echo.Context) error {
	cc, ok := c.(*localware.LocalUserClerkDbContext)
	if !ok {
		c.Logger().Error("could not resolve cc")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "could not resolve cc",
		})
	}

	var inc_tiers []tiers.Tiers

	err := cc.Bind(&inc_tiers)
	if err != nil {
		cc.Logger().Errorf("could not bind the request: ", err)
		return cc.JSON(http.StatusBadRequest, map[string]string{
			"error": "could not bind the request",
		})
	}

	fmt.Println("Updates coming in: ", inc_tiers)

	err = tiers.UpdateTiers(cc.Db, inc_tiers)
	if err != nil {
		cc.Logger().Errorf("Could not update tiers: ", err)
		return cc.JSON(http.StatusBadRequest, map[string]string{
			"error": "could not bind the request",
		})
	}

	// Retrieve the updated val
	updatedTiers, err := tiers.ReadTiers(cc.Db)
	if err != nil {
		c.Logger().Errorf("could not read tiers in the update: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "could not read the weights in the update",
		})
	}

	fmt.Println("updated tiers: ", updatedTiers)

	return helpers.Success(cc, updatedTiers)
}
