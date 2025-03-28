package stations_routes

import (
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"hopdf.com/localware"
)

// TODO: START HERE
func UpdateStation(c echo.Context) error {
	cc, ok := c.(*localware.LocalUserClerkDbContext)
	if !ok {
		c.Logger().Error("could not resolve cc")
		os.Exit(1)
	}

	stationIDStr := c.Param("station_id")

	stationID, err := strconv.Atoi(stationIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid station_id",
		})
	}

	return nil
}
